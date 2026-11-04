package observability

// Logger is a minimal structured logger for services.
type Logger interface {
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
}

// NopLogger discards log lines.
type NopLogger struct{}

func (NopLogger) Info(string, ...any) {}
func (NopLogger) Warn(string, ...any) {}
