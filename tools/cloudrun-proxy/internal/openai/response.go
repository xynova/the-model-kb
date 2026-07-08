package openai

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type modelList struct {
	Object string  `json:"object"`
	Data   []model `json:"data"`
}

type model struct {
	ID     string `json:"id"`
	Object string `json:"object"`
}

type ErrorBody struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Message string  `json:"message"`
	Type    string  `json:"type"`
	Code    *string `json:"code,omitempty"`
}

func WriteModels(w http.ResponseWriter, modelIDs []string) {
	w.Header().Set("Content-Type", "application/json")
	data := make([]model, len(modelIDs))
	for i, id := range modelIDs {
		data[i] = model{ID: id, Object: "model"}
	}
	_ = json.NewEncoder(w).Encode(modelList{
		Object: "list",
		Data:   data,
	})
}

func WriteError(w http.ResponseWriter, status int, errType, message string) {
	WriteErrorWithRetryAfter(w, status, errType, message, 0)
}

func WriteErrorWithRetryAfter(w http.ResponseWriter, status int, errType, message string, retryAfterSec int) {
	w.Header().Set("Content-Type", "application/json")
	if retryAfterSec > 0 {
		w.Header().Set("Retry-After", strconv.Itoa(retryAfterSec))
	}
	w.WriteHeader(status)
	code := http.StatusText(status)
	_ = json.NewEncoder(w).Encode(ErrorBody{
		Error: ErrorDetail{
			Message: message,
			Type:    errType,
			Code:    &code,
		},
	})
}
