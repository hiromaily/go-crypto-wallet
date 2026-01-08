package logger

import (
	"log/slog"
	"os"

	"github.com/phsym/console-slog"
)

//
// SlogLogger
//

type SlogLogger struct {
	log  *slog.Logger
	args []any
}

// NewSlogLogger creates a new slog logger from config parameters.
// format: "json" for JSON output, "console" for human-readable console output.
// commitID: commit ID to include in log output (empty string to skip).
func NewSlogLogger(format, levelStr, service, commitID string) Logger {
	level := getSlogLevel(levelStr)
	args := []any{
		slog.String("service", service),
	}
	if commitID != "" {
		args = append(args, slog.String("commit_id", commitID))
	}

	// logger option
	options := &slog.HandlerOptions{Level: level}

	// Choose handler based on format
	var handler slog.Handler
	switch format {
	case "console":
		handler = console.NewHandler(os.Stderr, &console.HandlerOptions{Level: level})
	default:
		// Default to JSON format for production use
		handler = slog.NewJSONHandler(os.Stdout, options)
	}

	return &SlogLogger{
		log:  slog.New(handler),
		args: args,
	}
}

// getSlogLevel converts string log level to slog.Level
func getSlogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Debug logs a debug message with the provided arguments.
func (s *SlogLogger) Debug(msg string, args ...any) {
	s.log.Debug(msg, s.appendArgs(args...)...)
}

// Info logs an info message with the provided arguments.
func (s *SlogLogger) Info(msg string, args ...any) {
	s.log.Info(msg, s.appendArgs(args...)...)
}

// Warn logs a warning message with the provided arguments.
func (s *SlogLogger) Warn(msg string, args ...any) {
	s.log.Warn(msg, s.appendArgs(args...)...)
}

// Error logs an error message with the provided arguments.
func (s *SlogLogger) Error(msg string, args ...any) {
	s.log.Error(msg, s.appendArgs(args...)...)
}

// appends the args to the default args
func (s *SlogLogger) appendArgs(args ...any) []any {
	return append(s.args, args...)
}
