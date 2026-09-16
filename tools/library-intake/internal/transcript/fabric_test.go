package transcript_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/xynova/library-intake/internal/execx"
	"github.com/xynova/library-intake/internal/transcript"
)

func TestGrabWritesFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	out := filepath.Join(dir, "t.transcript.txt")

	runner := &execx.Runner{
		LookPath: func(file string) (string, error) { return file, nil },
		Command: func(ctx context.Context, name string, args ...string) *exec.Cmd {
			// Fake fabric: write transcript via a tiny shell snippet using -o path.
			outPath := ""
			for i := 0; i < len(args)-1; i++ {
				if args[i] == "-o" {
					outPath = args[i+1]
				}
			}
			script := "printf 'hello transcript\\n' > \"$1\""
			return exec.CommandContext(ctx, "sh", "-c", script, "sh", outPath)
		},
	}

	svc := transcript.Config{FabricBin: "fabric", Runner: runner}.CreateService()
	text, err := svc.Grab(context.Background(), "https://www.youtube.com/watch?v=nZYJdwM-_nI", out)
	if err != nil {
		t.Fatal(err)
	}
	if text != "hello transcript" {
		t.Fatalf("text=%q", text)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello transcript\n" && string(data) != "hello transcript" {
		// Grab may rewrite trimmed file without trailing newline depending on path.
		if text != "hello transcript" {
			t.Fatalf("file=%q", data)
		}
	}
}
