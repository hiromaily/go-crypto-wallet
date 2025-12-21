package logger

//
// Noop Logger
//

type NoopLogger struct{}

func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

// Debug
func (*NoopLogger) Debug(msg string, args ...any) {
	// DummyLogger disables logging
}

// Info
func (*NoopLogger) Info(msg string, args ...any) {
	// DummyLogger disables logging
}

// Warn
func (*NoopLogger) Warn(msg string, args ...any) {
	// DummyLogger disables logging
}

// Error
func (*NoopLogger) Error(msg string, args ...any) {
	// DummyLogger disables logging
}
