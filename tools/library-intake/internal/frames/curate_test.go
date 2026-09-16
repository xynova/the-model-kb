package frames_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/xynova/library-intake/internal/candidates"
	"github.com/xynova/library-intake/internal/clients/polypus"
	"github.com/xynova/library-intake/internal/frames"
)

func TestCurateWritesMap(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]string{"content": `{
  "assignments": [
    {
      "file": "frame_0001.jpg",
      "candidate_title": "VoiceMem",
      "role": "cover",
      "reason": "demo still"
    }
  ]
}`}},
			},
		})
	}))
	t.Cleanup(srv.Close)

	dir := t.TempDir()
	framesDir := filepath.Join(dir, "frames")
	if err := os.MkdirAll(framesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Minimal JPEG-ish bytes are fine; curator only base64-encodes.
	jpegish := []byte{0xff, 0xd8, 0xff, 0xd9}
	if err := os.WriteFile(filepath.Join(framesDir, "frame_0001.jpg"), jpegish, 0o644); err != nil {
		t.Fatal(err)
	}
	_ = base64.StdEncoding.EncodeToString(jpegish)

	out := filepath.Join(dir, "frame-map.json")
	client := polypus.Config{BaseURL: srv.URL, HTTPClient: srv.Client(), Timeout: time.Second}.CreateClient()
	curator, err := frames.CurateConfig{Client: client, Model: "gemma", BatchSize: 4}.CreateCurator()
	if err != nil {
		t.Fatal(err)
	}
	list := &candidates.List{Candidates: []candidates.Candidate{{
		Title: "VoiceMem", Category: "agents-memory", Why: "x", ColdRead: "y",
	}}}
	m, err := curator.Curate(context.Background(), framesDir, list, out)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Assignments) != 1 || m.Assignments[0].Role != "cover" {
		t.Fatalf("%#v", m)
	}
}
