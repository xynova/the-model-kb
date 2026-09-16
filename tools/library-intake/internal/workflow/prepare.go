package workflow

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xynova/library-intake/internal/candidates"
	"github.com/xynova/library-intake/internal/clients/polypus"
	"github.com/xynova/library-intake/internal/config"
	derrors "github.com/xynova/library-intake/internal/errors"
	"github.com/xynova/library-intake/internal/execx"
	"github.com/xynova/library-intake/internal/frames"
	"github.com/xynova/library-intake/internal/transcript"
	"github.com/xynova/library-intake/internal/youtubeid"
)

// PrepareResult summarizes a prepare run for operators.
type PrepareResult struct {
	VideoID         string
	TranscriptPath  string
	CandidatesPath  string
	FramesDir       string
	FrameMapPath    string
	CandidateCount  int
	FrameCount      int
	MediaSkipped    bool
	MediaSkipReason string
}

// Preparer runs the YouTube intake preparation pipeline.
type Preparer struct {
	cfg    config.Config
	runner *execx.Runner
}

// Config holds preparer construction inputs.
type Config struct {
	App    config.Config
	Runner *execx.Runner
}

// CreatePreparer builds a Preparer.
func (c Config) CreatePreparer() *Preparer {
	runner := c.Runner
	if runner == nil {
		runner = execx.NewRunner()
	}
	return &Preparer{cfg: c.App, runner: runner}
}

// Prepare runs transcript + candidates, then soft-fails frames + curation.
func (p *Preparer) Prepare(ctx context.Context, youtubeURL, workRoot string) (*PrepareResult, error) {
	if p == nil {
		return nil, derrors.New(derrors.CodeFailed, "workflow.Prepare", "preparer is nil")
	}
	if workRoot == "" {
		workRoot = p.cfg.WorkRoot
	}
	videoID, err := youtubeid.FromURL(youtubeURL)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(workRoot, 0o755); err != nil {
		return nil, derrors.Wrap(err, derrors.CodeFailed, "workflow.Prepare", "mkdir work root").
			With("path", workRoot)
	}

	transcriptPath := filepath.Join(workRoot, videoID+".transcript.txt")
	candidatesPath := filepath.Join(workRoot, videoID+".candidates.parsed.json")
	framesDir := filepath.Join(workRoot, videoID, "frames")
	frameMapPath := filepath.Join(workRoot, videoID+".frame-map.json")

	tx := transcript.Config{FabricBin: p.cfg.FabricBin, Runner: p.runner}.CreateService()
	if _, err := tx.Grab(ctx, youtubeURL, transcriptPath); err != nil {
		return nil, err
	}

	chatClient := polypus.Config{
		BaseURL: p.cfg.BaseURL,
		Timeout: p.cfg.ChatTimeout,
	}.CreateClient()
	proposer, err := candidates.ProposeConfig{
		Client: chatClient,
		Model:  p.cfg.ChatModel,
	}.CreateProposer()
	if err != nil {
		return nil, err
	}
	list, err := proposer.ProposeFromFile(ctx, transcriptPath, candidatesPath)
	if err != nil {
		return nil, err
	}

	result := &PrepareResult{
		VideoID:        videoID,
		TranscriptPath: transcriptPath,
		CandidatesPath: candidatesPath,
		FramesDir:      framesDir,
		FrameMapPath:   frameMapPath,
		CandidateCount: len(list.Candidates),
	}

	extractor := frames.ExtractConfig{
		YtDlpBin:  p.cfg.YtDlpBin,
		FFmpegBin: p.cfg.FFmpegBin,
		Runner:    p.runner,
	}.CreateExtractor()
	frameCount, frameErr := extractor.Extract(ctx, youtubeURL, framesDir, p.cfg.FrameIntervalS)
	if frameErr != nil {
		result.MediaSkipped = true
		result.MediaSkipReason = frameErr.Error()
		fmt.Fprintf(os.Stderr, "library-intake: frames soft-fail: %v\n", frameErr)
		return result, nil
	}
	result.FrameCount = frameCount

	visionClient := polypus.Config{
		BaseURL: p.cfg.BaseURL,
		Timeout: p.cfg.VisionTimeout,
	}.CreateClient()
	curator, err := frames.CurateConfig{
		Client:    visionClient,
		Model:     p.cfg.VisionModel,
		BatchSize: p.cfg.VisionBatch,
	}.CreateCurator()
	if err != nil {
		result.MediaSkipped = true
		result.MediaSkipReason = err.Error()
		fmt.Fprintf(os.Stderr, "library-intake: curate soft-fail: %v\n", err)
		return result, nil
	}
	if _, err := curator.Curate(ctx, framesDir, list, frameMapPath); err != nil {
		result.MediaSkipped = true
		result.MediaSkipReason = err.Error()
		fmt.Fprintf(os.Stderr, "library-intake: curate soft-fail: %v\n", err)
		return result, nil
	}
	return result, nil
}
