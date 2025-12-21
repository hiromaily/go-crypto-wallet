package logger

//
// Noop Logger
//

type NoopLogger struct{}

func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

// Debug logs a debug message (no-op implementation).
func (*NoopLogger) Debug(msg string, args ...any) {
	// DummyLogger disables logging
}

// Info logs an info message (no-op implementation).
func (*NoopLogger) Info(msg string, args ...any) {
	// DummyLogger disables logging
}

// Warn logs a warning message (no-op implementation).
func (*NoopLogger) Warn(msg string, args ...any) {
	// DummyLogger disables logging
}

// Error logs an error message (no-op implementation).
func (*NoopLogger) Error(msg string, args ...any) {
	// DummyLogger disables logging
}
