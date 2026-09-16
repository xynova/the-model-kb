package frames

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	derrors "github.com/xynova/library-intake/internal/errors"
	"github.com/xynova/library-intake/internal/execx"
)

// ExtractConfig holds frame extraction construction inputs.
type ExtractConfig struct {
	YtDlpBin  string
	FFmpegBin string
	Runner    *execx.Runner
}

// Extractor downloads a YouTube video and extracts sparse JPEG frames.
type Extractor struct {
	ytDlp  string
	ffmpeg string
	runner *execx.Runner
}

// CreateExtractor builds an Extractor.
func (c ExtractConfig) CreateExtractor() *Extractor {
	runner := c.Runner
	if runner == nil {
		runner = execx.NewRunner()
	}
	yt := strings.TrimSpace(c.YtDlpBin)
	if yt == "" {
		yt = "yt-dlp"
	}
	ff := strings.TrimSpace(c.FFmpegBin)
	if ff == "" {
		ff = "ffmpeg"
	}
	return &Extractor{ytDlp: yt, ffmpeg: ff, runner: runner}
}

// Extract downloads the video into a scratch dir under outDir and writes frame_*.jpg.
func (e *Extractor) Extract(ctx context.Context, youtubeURL, outDir string, intervalSeconds int) (int, error) {
	if e == nil {
		return 0, derrors.New(derrors.CodeFailed, "frames.Extract", "extractor is nil")
	}
	youtubeURL = strings.TrimSpace(youtubeURL)
	outDir = strings.TrimSpace(outDir)
	if youtubeURL == "" {
		return 0, derrors.New(derrors.CodeInvalidArgument, "frames.Extract", "url is empty")
	}
	if outDir == "" {
		return 0, derrors.New(derrors.CodeInvalidArgument, "frames.Extract", "out dir is empty")
	}
	if intervalSeconds <= 0 {
		intervalSeconds = 8
	}
	if err := ctx.Err(); err != nil {
		return 0, derrors.Wrap(err, derrors.CodeFailed, "frames.Extract", "context done")
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return 0, derrors.Wrap(err, derrors.CodeFailed, "frames.Extract", "mkdir frames dir").
			With("path", outDir)
	}

	scratch, err := os.MkdirTemp(filepath.Dir(outDir), "library-intake-video-*")
	if err != nil {
		return 0, derrors.Wrap(err, derrors.CodeFailed, "frames.Extract", "mkdir scratch")
	}
	defer os.RemoveAll(scratch)

	outTemplate := filepath.Join(scratch, "video.%(ext)s")
	_, err = e.runner.Run(ctx, e.ytDlp,
		"--no-update",
		"-f", "bv*[height<=720][ext=mp4]+ba[ext=m4a]/b[height<=720]/best",
		"-o", outTemplate,
		"--merge-output-format", "mp4",
		youtubeURL,
	)
	if err != nil {
		return 0, err
	}

	videoPath, err := findVideo(scratch)
	if err != nil {
		return 0, err
	}

	// Clear previous frames for deterministic reruns of the same outDir.
	prev, _ := filepath.Glob(filepath.Join(outDir, "frame_*.jpg"))
	for _, p := range prev {
		_ = os.Remove(p)
	}

	framePattern := filepath.Join(outDir, "frame_%04d.jpg")
	_, err = e.runner.Run(ctx, e.ffmpeg,
		"-hide_banner", "-loglevel", "error", "-y",
		"-i", videoPath,
		"-vf", fmt.Sprintf("fps=1/%d", intervalSeconds),
		"-q:v", "3",
		framePattern,
	)
	if err != nil {
		return 0, err
	}

	files, err := filepath.Glob(filepath.Join(outDir, "frame_*.jpg"))
	if err != nil {
		return 0, derrors.Wrap(err, derrors.CodeFailed, "frames.Extract", "list frames")
	}
	if len(files) == 0 {
		return 0, derrors.New(derrors.CodeFailed, "frames.Extract", "no frames written").
			With("outdir", outDir)
	}
	return len(files), nil
}

func findVideo(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", derrors.Wrap(err, derrors.CodeFailed, "frames.findVideo", "read scratch")
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, "video.") {
			return filepath.Join(dir, name), nil
		}
	}
	return "", derrors.New(derrors.CodeFailed, "frames.findVideo", "yt-dlp produced no video file").
		With("dir", dir)
}
