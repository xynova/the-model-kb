package polypus_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/xynova/library-intake/internal/clients/polypus"
)

func TestExtractJSONObjectStripsFence(t *testing.T) {
	t.Parallel()
	raw, err := polypus.ExtractJSONObject("```json\n{\"a\":1}\n```")
	if err != nil {
		t.Fatal(err)
	}
	var obj map[string]int
	if err := json.Unmarshal(raw, &obj); err != nil {
		t.Fatal(err)
	}
	if obj["a"] != 1 {
		t.Fatalf("got %#v", obj)
	}
}

func TestChatCompletions(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]string{"content": `{"ok":true}`}},
			},
		})
	}))
	t.Cleanup(srv.Close)

	client := polypus.Config{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
		Timeout:    5 * time.Second,
	}.CreateClient()
	content, err := client.ChatCompletions(context.Background(), "test-model", 0.1, []polypus.Message{
		{Role: "user", Content: "hi"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(content, `"ok":true`) {
		t.Fatalf("content=%q", content)
	}
}

func TestChatCompletionsHTTPError(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":{"message":"down"}}`))
	}))
	t.Cleanup(srv.Close)

	client := polypus.Config{BaseURL: srv.URL, HTTPClient: srv.Client(), Timeout: time.Second}.CreateClient()
	_, err := client.ChatCompletions(context.Background(), "m", 0, []polypus.Message{{Role: "user", Content: "x"}})
	if err == nil {
		t.Fatal("expected error")
	}
}
