package openai

import "encoding/json"

// ChatSummary is a redacted view of a chat completion request for logging.
type ChatSummary struct {
	Model           string `json:"model"`
	UpstreamModel   string `json:"upstream_model,omitempty"`
	Stream          bool   `json:"stream"`
	MaxTokens       any    `json:"max_tokens,omitempty"`
	MessageCount    int    `json:"message_count"`
	EnableThinking  any    `json:"enable_thinking,omitempty"`
	BodyBytes       int    `json:"body_bytes"`
	UpstreamBytes   int    `json:"upstream_bytes,omitempty"`
}

func SummarizeChatRequest(body []byte) ChatSummary {
	var req map[string]any
	if err := json.Unmarshal(body, &req); err != nil {
		return ChatSummary{BodyBytes: len(body)}
	}

	summary := ChatSummary{
		Model:      stringField(req["model"]),
		Stream:     streamField(req["stream"]),
		MaxTokens:  req["max_tokens"],
		BodyBytes:  len(body),
		MessageCount: messageCount(req["messages"]),
	}
	if kwargs, ok := req["chat_template_kwargs"].(map[string]any); ok {
		summary.EnableThinking = kwargs["enable_thinking"]
	}
	return summary
}

func SummarizeUpstreamRequest(body []byte, upstreamModel string) ChatSummary {
	s := SummarizeChatRequest(body)
	s.UpstreamModel = upstreamModel
	s.UpstreamBytes = len(body)
	return s
}

func stringField(v any) string {
	s, _ := v.(string)
	return s
}

func streamField(v any) bool {
	switch t := v.(type) {
	case bool:
		return t
	default:
		return false
	}
}

func messageCount(v any) int {
	msgs, ok := v.([]any)
	if !ok {
		return 0
	}
	return len(msgs)
}
