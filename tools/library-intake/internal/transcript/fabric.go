package transcript

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	derrors "github.com/xynova/library-intake/internal/errors"
	"github.com/xynova/library-intake/internal/execx"
)

// Config holds Fabric transcript construction inputs.
type Config struct {
	FabricBin string
	Runner    *execx.Runner
}

// Service grabs YouTube transcripts via fabric.
type Service struct {
	fabricBin string
	runner    *execx.Runner
}

// CreateService builds a transcript Service.
func (c Config) CreateService() *Service {
	runner := c.Runner
	if runner == nil {
		runner = execx.NewRunner()
	}
	bin := strings.TrimSpace(c.FabricBin)
	if bin == "" {
		bin = "fabric"
	}
	return &Service{fabricBin: bin, runner: runner}
}

// Grab writes a transcript-only fabric capture to outPath and returns the text.
func (s *Service) Grab(ctx context.Context, youtubeURL, outPath string) (string, error) {
	if s == nil {
		return "", derrors.New(derrors.CodeFailed, "transcript.Grab", "service is nil")
	}
	youtubeURL = strings.TrimSpace(youtubeURL)
	outPath = strings.TrimSpace(outPath)
	if youtubeURL == "" {
		return "", derrors.New(derrors.CodeInvalidArgument, "transcript.Grab", "url is empty")
	}
	if outPath == "" {
		return "", derrors.New(derrors.CodeInvalidArgument, "transcript.Grab", "out path is empty")
	}
	if err := ctx.Err(); err != nil {
		return "", derrors.Wrap(err, derrors.CodeFailed, "transcript.Grab", "context done")
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return "", derrors.Wrap(err, derrors.CodeFailed, "transcript.Grab", "mkdir out dir").
			With("path", outPath)
	}

	// Transcript-only: no -p pattern so fabric exits without a chat LLM call.
	res, err := s.runner.Run(ctx, s.fabricBin, "-y", youtubeURL, "--transcript", "-o", outPath)
	if err != nil {
		return "", err
	}
	data, readErr := os.ReadFile(outPath)
	if readErr != nil {
		// Fabric may have printed only; fall back to stdout.
		text := strings.TrimSpace(res.Stdout)
		if text == "" {
			return "", derrors.Wrap(readErr, derrors.CodeFailed, "transcript.Grab", "read transcript file").
				With("path", outPath)
		}
		if writeErr := os.WriteFile(outPath, []byte(text+"\n"), 0o644); writeErr != nil {
			return "", derrors.Wrap(writeErr, derrors.CodeFailed, "transcript.Grab", "write transcript fallback").
				With("path", outPath)
		}
		return text, nil
	}
	text := strings.TrimSpace(string(data))
	if text == "" {
		return "", derrors.New(derrors.CodeFailed, "transcript.Grab", "transcript is empty").
			With("path", outPath)
	}
	return text, nil
}
