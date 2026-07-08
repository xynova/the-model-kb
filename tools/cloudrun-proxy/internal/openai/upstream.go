package openai

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/xynova/cloudrun-proxy/internal/config"
)

// PrepareUpstreamChatBody adjusts client requests before forwarding to llama-server.
//
// Default (GEMMA_ENABLE_THINKING=false): disable thinking for OpenAI-compatible
// clients like RooCode that only read message.content.
//
// With thinking enabled: leave thinking on and raise low max_tokens values so
// the reasoning phase does not consume the entire generation budget before a
// final answer is emitted.
func PrepareUpstreamChatBody(body []byte, model config.ResolvedModel) ([]byte, error) {
	var req map[string]any
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("invalid JSON body: %w", err)
	}

	req["model"] = model.UpstreamModel

	kwargs, ok := req["chat_template_kwargs"].(map[string]any)
	if !ok {
		kwargs = map[string]any{}
	}

	if _, ok := kwargs["enable_thinking"]; !ok {
		kwargs["enable_thinking"] = model.EnableThinking
	}
	req["chat_template_kwargs"] = kwargs

	minTokens := model.MinMaxTokens
	if thinkingEnabled(kwargs) {
		minTokens = model.MinMaxTokensThinking
	}
	ensureMinMaxTokens(req, minTokens)

	out, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal upstream body: %w", err)
	}
	return out, nil
}

func thinkingEnabled(kwargs map[string]any) bool {
	v, ok := kwargs["enable_thinking"]
	if !ok {
		return false
	}
	switch t := v.(type) {
	case bool:
		return t
	default:
		return false
	}
}

func ensureMinMaxTokens(req map[string]any, min int) {
	if min <= 0 {
		return
	}
	current, ok := maxTokensValue(req["max_tokens"])
	if !ok || current < min {
		req["max_tokens"] = min
	}
}

func maxTokensValue(v any) (int, bool) {
	switch t := v.(type) {
	case float64:
		if t < 0 || t > math.MaxInt32 {
			return 0, false
		}
		return int(t), true
	case int:
		return t, true
	case int64:
		return int(t), true
	case json.Number:
		n, err := t.Int64()
		if err != nil || n < 0 || n > math.MaxInt32 {
			return 0, false
		}
		return int(n), true
	default:
		return 0, false
	}
}

// NormalizeChatCompletion copies reasoning_content into content when content is
// empty, so clients that only read message.content still get a response when
// thinking consumed the token budget or the client ignores reasoning_content.
func NormalizeChatCompletion(body []byte) ([]byte, error) {
	var resp map[string]any
	if err := json.Unmarshal(body, &resp); err != nil {
		return body, nil
	}

	choices, ok := resp["choices"].([]any)
	if !ok || len(choices) == 0 {
		return body, nil
	}

	choice, ok := choices[0].(map[string]any)
	if !ok {
		return body, nil
	}

	message, ok := choice["message"].(map[string]any)
	if !ok {
		return body, nil
	}

	content, _ := message["content"].(string)
	reasoning, _ := message["reasoning_content"].(string)
	if content == "" && reasoning != "" {
		message["content"] = reasoning
	}

	out, err := json.Marshal(resp)
	if err != nil {
		return body, nil
	}
	return out, nil
}
