package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/xynova/library-intake/internal/candidates"
	"github.com/xynova/library-intake/internal/catalog"
	"github.com/xynova/library-intake/internal/clients/polypus"
	"github.com/xynova/library-intake/internal/config"
	derrors "github.com/xynova/library-intake/internal/errors"
	"github.com/xynova/library-intake/internal/execx"
	"github.com/xynova/library-intake/internal/frames"
	"github.com/xynova/library-intake/internal/transcript"
	"github.com/xynova/library-intake/internal/workflow"
	"github.com/xynova/library-intake/internal/youtubeid"
)

const usageTmpl = `library-intake — YouTube intake prep for library/

Usage:
  library-intake <command> [flags]

Commands:
  transcript   Grab a YouTube transcript via fabric
  candidates   Propose candidates from a transcript via Polypus GLM
  frames       Extract sparse JPEG frames via yt-dlp + ffmpeg
  curate       Map frames to candidates via Polypus Gemma vision
  prepare      transcript + candidates, then soft-fail frames + curate
  index        Refresh library INDEX.md files and weekly additions

Environment:
  POLYPUS_BASE_URL                 default {{.BaseURL}}
  LIBRARY_INTAKE_POLYPUS_MODEL     default {{.ChatModel}}
  LIBRARY_INTAKE_VISION_MODEL      default {{.VisionModel}}
  LIBRARY_INTAKE_WORK_ROOT         default {{.WorkRoot}}
`

// Run dispatches a subcommand.
func Run(ctx context.Context, args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		return printUsage()
	}
	cfg := config.Load()
	switch args[0] {
	case "transcript":
		return runTranscript(ctx, cfg, args[1:])
	case "candidates":
		return runCandidates(ctx, cfg, args[1:])
	case "frames":
		return runFrames(ctx, cfg, args[1:])
	case "curate":
		return runCurate(ctx, cfg, args[1:])
	case "prepare":
		return runPrepare(ctx, cfg, args[1:])
	case "index":
		return runIndex(ctx, cfg, args[1:])
	default:
		return derrors.New(derrors.CodeInvalidArgument, "cli.Run", "unknown command: "+args[0])
	}
}

func printUsage() error {
	t, err := template.New("usage").Parse(usageTmpl)
	if err != nil {
		return derrors.Wrap(err, derrors.CodeFailed, "cli.printUsage", "parse usage template")
	}
	cfg := config.Load()
	return t.Execute(os.Stdout, cfg)
}

func runTranscript(ctx context.Context, cfg config.Config, args []string) error {
	fs := flag.NewFlagSet("transcript", flag.ContinueOnError)
	url := fs.String("url", "", "YouTube URL")
	out := fs.String("out", "", "transcript output path")
	if err := fs.Parse(args); err != nil {
		return derrors.Wrap(err, derrors.CodeInvalidArgument, "cli.transcript", "parse flags")
	}
	if strings.TrimSpace(*url) == "" {
		return derrors.New(derrors.CodeInvalidArgument, "cli.transcript", "--url is required")
	}
	outPath := strings.TrimSpace(*out)
	if outPath == "" {
		id, err := youtubeid.FromURL(*url)
		if err != nil {
			return err
		}
		outPath = filepath.Join(cfg.WorkRoot, id+".transcript.txt")
	}
	svc := transcript.Config{FabricBin: cfg.FabricBin, Runner: execx.NewRunner()}.CreateService()
	text, err := svc.Grab(ctx, *url, outPath)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "wrote %s (%d bytes)\n", outPath, len(text))
	return nil
}

func runCandidates(ctx context.Context, cfg config.Config, args []string) error {
	fs := flag.NewFlagSet("candidates", flag.ContinueOnError)
	transcriptPath := fs.String("transcript", "", "transcript file path")
	out := fs.String("out", "", "candidates JSON output path")
	if err := fs.Parse(args); err != nil {
		return derrors.Wrap(err, derrors.CodeInvalidArgument, "cli.candidates", "parse flags")
	}
	if strings.TrimSpace(*transcriptPath) == "" {
		return derrors.New(derrors.CodeInvalidArgument, "cli.candidates", "--transcript is required")
	}
	outPath := strings.TrimSpace(*out)
	if outPath == "" {
		base := strings.TrimSuffix(filepath.Base(*transcriptPath), ".transcript.txt")
		base = strings.TrimSuffix(base, filepath.Ext(base))
		outPath = filepath.Join(filepath.Dir(*transcriptPath), base+".candidates.parsed.json")
	}
	client := polypus.Config{BaseURL: cfg.BaseURL, Timeout: cfg.ChatTimeout}.CreateClient()
	proposer, err := candidates.ProposeConfig{Client: client, Model: cfg.ChatModel}.CreateProposer()
	if err != nil {
		return err
	}
	list, err := proposer.ProposeFromFile(ctx, *transcriptPath, outPath)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "wrote %s (%d candidates)\n", outPath, len(list.Candidates))
	return nil
}

func runFrames(ctx context.Context, cfg config.Config, args []string) error {
	fs := flag.NewFlagSet("frames", flag.ContinueOnError)
	url := fs.String("url", "", "YouTube URL")
	outDir := fs.String("outdir", "", "frames output directory")
	interval := fs.Int("interval", cfg.FrameIntervalS, "seconds between frames")
	if err := fs.Parse(args); err != nil {
		return derrors.Wrap(err, derrors.CodeInvalidArgument, "cli.frames", "parse flags")
	}
	if strings.TrimSpace(*url) == "" {
		return derrors.New(derrors.CodeInvalidArgument, "cli.frames", "--url is required")
	}
	dir := strings.TrimSpace(*outDir)
	if dir == "" {
		id, err := youtubeid.FromURL(*url)
		if err != nil {
			return err
		}
		dir = filepath.Join(cfg.WorkRoot, id, "frames")
	}
	extractor := frames.ExtractConfig{
		YtDlpBin:           cfg.YtDlpBin,
		FFmpegBin:          cfg.FFmpegBin,
		YtDlpExtractorArgs: cfg.YtDlpExtractorArgs,
		Runner:             execx.NewRunner(),
	}.CreateExtractor()
	n, err := extractor.Extract(ctx, *url, dir, *interval)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "wrote %d frames to %s\n", n, dir)
	return nil
}

func runCurate(ctx context.Context, cfg config.Config, args []string) error {
	fs := flag.NewFlagSet("curate", flag.ContinueOnError)
	framesDir := fs.String("frames-dir", "", "directory of frame_*.jpg")
	candidatesPath := fs.String("candidates", "", "candidates.parsed.json path")
	out := fs.String("out", "", "frame-map JSON output path")
	batch := fs.Int("batch-size", cfg.VisionBatch, "frames per vision request")
	if err := fs.Parse(args); err != nil {
		return derrors.Wrap(err, derrors.CodeInvalidArgument, "cli.curate", "parse flags")
	}
	if strings.TrimSpace(*framesDir) == "" {
		return derrors.New(derrors.CodeInvalidArgument, "cli.curate", "--frames-dir is required")
	}
	if strings.TrimSpace(*candidatesPath) == "" {
		return derrors.New(derrors.CodeInvalidArgument, "cli.curate", "--candidates is required")
	}
	outPath := strings.TrimSpace(*out)
	if outPath == "" {
		outPath = filepath.Join(filepath.Dir(*candidatesPath),
			strings.TrimSuffix(filepath.Base(*candidatesPath), ".candidates.parsed.json")+".frame-map.json")
	}
	raw, err := os.ReadFile(*candidatesPath)
	if err != nil {
		return derrors.Wrap(err, derrors.CodeNotFound, "cli.curate", "read candidates").
			With("path", *candidatesPath)
	}
	var list candidates.List
	if err := json.Unmarshal(raw, &list); err != nil {
		return derrors.Wrap(err, derrors.CodeParse, "cli.curate", "parse candidates JSON")
	}
	client := polypus.Config{BaseURL: cfg.BaseURL, Timeout: cfg.VisionTimeout}.CreateClient()
	curator, err := frames.CurateConfig{
		Client:    client,
		Model:     cfg.VisionModel,
		BatchSize: *batch,
	}.CreateCurator()
	if err != nil {
		return err
	}
	m, err := curator.Curate(ctx, *framesDir, &list, outPath)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "wrote %s (%d assignments)\n", outPath, len(m.Assignments))
	return nil
}

func runPrepare(ctx context.Context, cfg config.Config, args []string) error {
	fs := flag.NewFlagSet("prepare", flag.ContinueOnError)
	url := fs.String("url", "", "YouTube URL")
	workRoot := fs.String("workdir", cfg.WorkRoot, "working directory for artifacts")
	if err := fs.Parse(args); err != nil {
		return derrors.Wrap(err, derrors.CodeInvalidArgument, "cli.prepare", "parse flags")
	}
	if strings.TrimSpace(*url) == "" {
		return derrors.New(derrors.CodeInvalidArgument, "cli.prepare", "--url is required")
	}
	preparer := workflow.Config{App: cfg, Runner: execx.NewRunner()}.CreatePreparer()
	result, err := preparer.Prepare(ctx, *url, *workRoot)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		return derrors.Wrap(err, derrors.CodeFailed, "cli.prepare", "encode result")
	}
	return nil
}

func runIndex(ctx context.Context, cfg config.Config, args []string) error {
	_ = ctx
	_ = cfg
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		return derrors.New(derrors.CodeInvalidArgument, "cli.index", "usage: library-intake index refresh --library <path>")
	}
	switch args[0] {
	case "refresh":
		fs := flag.NewFlagSet("index refresh", flag.ContinueOnError)
		library := fs.String("library", "", "path to library/ root")
		if err := fs.Parse(args[1:]); err != nil {
			return derrors.Wrap(err, derrors.CodeInvalidArgument, "cli.index.refresh", "parse flags")
		}
		root, err := resolveLibraryRoot(*library)
		if err != nil {
			return err
		}
		if err := catalog.RefreshIndexes(root); err != nil {
			return err
		}
		if err := catalog.RebuildAdditions(root, time.Now()); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "refreshed indexes under %s\n", root)
		return nil
	default:
		return derrors.New(derrors.CodeInvalidArgument, "cli.index", "unknown subcommand: "+args[0])
	}
}

func resolveLibraryRoot(flagValue string) (string, error) {
	op := "cli.resolveLibraryRoot"
	root := strings.TrimSpace(flagValue)
	if root == "" {
		root = strings.TrimSpace(os.Getenv("LIBRARY_INTAKE_LIBRARY_ROOT"))
	}
	if root == "" {
		// Default: ../../library relative to tools/library-intake when run from that dir.
		cwd, err := os.Getwd()
		if err != nil {
			return "", derrors.Wrap(err, derrors.CodeFailed, op, "getwd")
		}
		candidate := filepath.Clean(filepath.Join(cwd, "..", "..", "library"))
		if st, err := os.Stat(candidate); err == nil && st.IsDir() {
			root = candidate
		}
	}
	if root == "" {
		return "", derrors.New(derrors.CodeInvalidArgument, op, "--library is required (or set LIBRARY_INTAKE_LIBRARY_ROOT)")
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", derrors.Wrap(err, derrors.CodeInvalidArgument, op, "abs library path").With("library", root)
	}
	st, err := os.Stat(abs)
	if err != nil || !st.IsDir() {
		return "", derrors.Wrap(err, derrors.CodeNotFound, op, "library root is not a directory").
			With("library", abs)
	}
	return abs, nil
}
