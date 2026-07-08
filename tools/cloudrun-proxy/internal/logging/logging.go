package logging

import (
	"log/slog"
	"os"
	"strings"
)

func New(debug bool) *slog.Logger {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
}

func RedactAuthHeader(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if len(value) <= 12 {
		return "[REDACTED]"
	}
	return value[:7] + "…[REDACTED]"
}
