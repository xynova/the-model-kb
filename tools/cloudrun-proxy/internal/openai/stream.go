package openai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
)

// RewriteSSEStream passes SSE through, copying reasoning_content into content
// deltas when content is empty so OpenAI clients like RooCode see output.
func RewriteSSEStream(r io.Reader, w io.Writer, flush func()) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := append([]byte(nil), scanner.Bytes()...)
		line = append(line, '\n')

		if out, ok := normalizeSSELine(scanner.Bytes()); ok {
			if _, err := w.Write(out); err != nil {
				return err
			}
			if _, err := w.Write([]byte("\n")); err != nil {
				return err
			}
		} else {
			if _, err := w.Write(line); err != nil {
				return err
			}
		}
		if flush != nil {
			flush()
		}
	}
	return scanner.Err()
}

func normalizeSSELine(line []byte) ([]byte, bool) {
	const prefix = "data:"
	if !bytes.HasPrefix(line, []byte(prefix)) {
		return nil, false
	}
	payload := bytes.TrimSpace(line[len(prefix):])
	if len(payload) == 0 || bytes.Equal(payload, []byte("[DONE]")) {
		return nil, false
	}

	var obj map[string]any
	if err := json.Unmarshal(payload, &obj); err != nil {
		return nil, false
	}

	choices, ok := obj["choices"].([]any)
	if !ok || len(choices) == 0 {
		return nil, false
	}
	choice, ok := choices[0].(map[string]any)
	if !ok {
		return nil, false
	}

	delta, ok := choice["delta"].(map[string]any)
	if !ok {
		return nil, false
	}

	content, _ := delta["content"].(string)
	reasoning, _ := delta["reasoning_content"].(string)
	if content == "" && reasoning != "" {
		delta["content"] = reasoning
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return nil, false
	}
	return append(append([]byte(prefix+" "), out...), '\n'), true
}
