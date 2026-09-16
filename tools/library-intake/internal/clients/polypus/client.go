package polypus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	derrors "github.com/xynova/library-intake/internal/errors"
)

var fenceRe = regexp.MustCompile("(?s)```(?:json)?\\s*(.*?)```")

// Config holds Polypus client construction inputs.
type Config struct {
	BaseURL    string
	HTTPClient *http.Client
	Timeout    time.Duration
}

// Client talks to an OpenAI-compatible Polypus gateway.
type Client struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

// CreateClient builds a Client from Config.
func (c Config) CreateClient() *Client {
	timeout := c.Timeout
	if timeout <= 0 {
		timeout = 180 * time.Second
	}
	hc := c.HTTPClient
	if hc == nil {
		hc = &http.Client{Timeout: timeout}
	}
	base := strings.TrimRight(strings.TrimSpace(c.BaseURL), "/")
	return &Client{baseURL: base, httpClient: hc, timeout: timeout}
}

// CreateClient is a package-level alias for Config.CreateClient.
func CreateClient(cfg Config) *Client {
	return cfg.CreateClient()
}

type chatRequest struct {
	Model       string    `json:"model"`
	Temperature float64   `json:"temperature"`
	Messages    []Message `json:"messages"`
}

// Message is an OpenAI-style chat message.
type Message struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

// TextPart is a multimodal text content part.
type TextPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ImagePart is a multimodal image content part.
type ImagePart struct {
	Type     string   `json:"type"`
	ImageURL ImageURL `json:"image_url"`
}

// ImageURL wraps a data URL or remote URL.
type ImageURL struct {
	URL string `json:"url"`
}

type chatResponse struct {
	Error   *apiError `json:"error"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type apiError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// ChatCompletions posts a chat completion and returns assistant text content.
func (c *Client) ChatCompletions(ctx context.Context, model string, temperature float64, messages []Message) (string, error) {
	if c == nil {
		return "", derrors.New(derrors.CodeFailed, "polypus.ChatCompletions", "client is nil")
	}
	if strings.TrimSpace(model) == "" {
		return "", derrors.New(derrors.CodeInvalidArgument, "polypus.ChatCompletions", "model is empty")
	}
	if len(messages) == 0 {
		return "", derrors.New(derrors.CodeInvalidArgument, "polypus.ChatCompletions", "messages are empty")
	}
	if err := ctx.Err(); err != nil {
		return "", derrors.Wrap(err, derrors.CodeFailed, "polypus.ChatCompletions", "context done")
	}

	body, err := json.Marshal(chatRequest{
		Model:       model,
		Temperature: temperature,
		Messages:    messages,
	})
	if err != nil {
		return "", derrors.Wrap(err, derrors.CodeFailed, "polypus.ChatCompletions", "marshal request")
	}

	reqCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	url := c.baseURL + "/v1/chat/completions"
	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", derrors.Wrap(err, derrors.CodeFailed, "polypus.ChatCompletions", "build request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", derrors.Wrap(err, derrors.CodeUnavailable, "polypus.ChatCompletions", "post chat").
			With("url", url).
			With("model", model)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", derrors.Wrap(err, derrors.CodeFailed, "polypus.ChatCompletions", "read body")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", derrors.New(derrors.CodeUnavailable, "polypus.ChatCompletions",
			fmt.Sprintf("unexpected status %d", resp.StatusCode)).
			With("body", trim(string(raw))).
			With("model", model)
	}

	var parsed chatResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", derrors.Wrap(err, derrors.CodeParse, "polypus.ChatCompletions", "decode response").
			With("body", trim(string(raw)))
	}
	if parsed.Error != nil {
		return "", derrors.New(derrors.CodeUnavailable, "polypus.ChatCompletions", parsed.Error.Message).
			With("type", parsed.Error.Type).
			With("model", model)
	}
	if len(parsed.Choices) == 0 {
		return "", derrors.New(derrors.CodeParse, "polypus.ChatCompletions", "empty choices")
	}
	content := strings.TrimSpace(parsed.Choices[0].Message.Content)
	if content == "" {
		return "", derrors.New(derrors.CodeParse, "polypus.ChatCompletions", "empty message content")
	}
	return content, nil
}

// ExtractJSONObject strips optional markdown fences and returns the first JSON object.
func ExtractJSONObject(content string) (json.RawMessage, error) {
	content = strings.TrimSpace(content)
	if m := fenceRe.FindStringSubmatch(content); len(m) == 2 {
		content = strings.TrimSpace(m[1])
	}
	start := strings.IndexByte(content, '{')
	end := strings.LastIndexByte(content, '}')
	if start < 0 || end < start {
		return nil, derrors.New(derrors.CodeParse, "polypus.ExtractJSONObject", "no JSON object found").
			With("content", trim(content))
	}
	raw := json.RawMessage(content[start : end+1])
	if !json.Valid(raw) {
		return nil, derrors.New(derrors.CodeParse, "polypus.ExtractJSONObject", "invalid JSON object").
			With("content", trim(string(raw)))
	}
	return raw, nil
}

func trim(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 2000 {
		return s[:2000] + "…"
	}
	return s
}
