package openai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type chatRequest struct {
	Model  string `json:"model"`
	Stream *bool  `json:"stream"`
}

func RequestModel(body []byte) (string, error) {
	var req chatRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return "", fmt.Errorf("invalid JSON body: %w", err)
	}
	model := strings.TrimSpace(req.Model)
	if model == "" {
		return "", fmt.Errorf("model is required")
	}
	return model, nil
}

func ValidateChatRequest(body []byte) error {
	if len(body) == 0 {
		return fmt.Errorf("request body is empty")
	}

	var req chatRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}
	return nil
}

func IsStreaming(body []byte) bool {
	var req chatRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return false
	}
	return req.Stream != nil && *req.Stream
}

func CheckLocalSecret(r *http.Request, secret string) error {
	if secret == "" {
		return nil
	}

	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if auth == "" {
		return fmt.Errorf("missing Authorization header")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return fmt.Errorf("Authorization must use Bearer scheme")
	}
	if strings.TrimPrefix(auth, prefix) != secret {
		return fmt.Errorf("invalid local secret")
	}
	return nil
}

func ReadBody(r *http.Request, limit int64) ([]byte, error) {
	defer r.Body.Close()
	reader := io.LimitReader(r.Body, limit)
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	if int64(len(body)) == limit {
		return nil, fmt.Errorf("request body exceeds %d bytes", limit)
	}
	return body, nil
}
