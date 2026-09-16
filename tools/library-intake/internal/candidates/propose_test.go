package candidates_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/xynova/library-intake/internal/candidates"
	"github.com/xynova/library-intake/internal/clients/polypus"
)

func TestProposeFromFile(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]string{"content": `{
  "candidates": [
    {
      "title": "VoiceMem",
      "name_hint": "voicemem",
      "category": "agents-memory",
      "why": "Local voice memory",
      "cold_read": "Useful for companion agents."
    }
  ]
}`}},
			},
		})
	}))
	t.Cleanup(srv.Close)

	dir := t.TempDir()
	tx := filepath.Join(dir, "x.transcript.txt")
	out := filepath.Join(dir, "x.candidates.parsed.json")
	if err := os.WriteFile(tx, []byte("VoiceMem is a local voice memory tool."), 0o644); err != nil {
		t.Fatal(err)
	}

	client := polypus.Config{BaseURL: srv.URL, HTTPClient: srv.Client(), Timeout: time.Second}.CreateClient()
	proposer, err := candidates.ProposeConfig{Client: client, Model: "glm"}.CreateProposer()
	if err != nil {
		t.Fatal(err)
	}
	list, err := proposer.ProposeFromFile(context.Background(), tx, out)
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Candidates) != 1 || list.Candidates[0].Title != "VoiceMem" {
		t.Fatalf("%#v", list)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatal(err)
	}
}
