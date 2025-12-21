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

func NewLogger(
	level slog.Level,
	commitID string,
) *SlogLogger {
	args := []any{
		slog.String("commit_id", commitID),
		// slog.Group(
		// 	"dd",
		// 	slog.String("service", ddService),
		// 	slog.String("version", ddVersion),
		// ),
	}
	// logger option
	options := &slog.HandlerOptions{Level: level}

	return &SlogLogger{
		log:  slog.New(slog.NewJSONHandler(os.Stdout, options)),
		args: args,
	}
}

// Debug
func (s *SlogLogger) Debug(msg string, args ...any) {
	s.log.Debug(msg, s.appendArgs(args...)...)
}

// Info
func (s *SlogLogger) Info(msg string, args ...any) {
	s.log.Info(msg, s.appendArgs(args...)...)
}

// Warn
func (s *SlogLogger) Warn(msg string, args ...any) {
	s.log.Warn(msg, s.appendArgs(args...)...)
}

// Error
func (s *SlogLogger) Error(msg string, args ...any) {
	s.log.Error(msg, s.appendArgs(args...)...)
}

// appends the args to the default args
func (s *SlogLogger) appendArgs(args ...any) []any {
	return append(s.args, args...)
}

// NewSlogLoggerWithLevel builder with fixed args
// this is test or local batch use
func NewSlogLoggerWithLevel(level slog.Level) *SlogLogger {
	return NewLogger(
		level,
		"commitid",
	)
}

//
// SlogConsoleLogger
// use https://github.com/phsym/console-slog
//

type SlogConsoleLogger struct {
	log  *slog.Logger
	args []any
}

// Localでのみ利用するため、重要ではない情報は保持しない

func NewConsoleLogger(
	level slog.Level,
) *SlogConsoleLogger {
	options := &console.HandlerOptions{Level: level}
	return &SlogConsoleLogger{
		log:  slog.New(console.NewHandler(os.Stderr, options)),
		args: []any{},
	}
}

// Debug
func (s *SlogConsoleLogger) Debug(msg string, args ...any) {
	s.log.Debug(msg, args...)
}

// Info
func (s *SlogConsoleLogger) Info(msg string, args ...any) {
	s.log.Info(msg, args...)
}

// Warn
func (s *SlogConsoleLogger) Warn(msg string, args ...any) {
	s.log.Warn(msg, args...)
}

// Error
func (s *SlogConsoleLogger) Error(msg string, args ...any) {
	s.log.Error(msg, args...)
}
