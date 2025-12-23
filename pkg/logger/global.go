package logger

import (
	"sync"
)

var (
	// globalLogger is the global logger instance used by package-level logging functions.
	globalLogger Logger
	// mu protects the globalLogger from concurrent access.
	mu sync.RWMutex
)

// SetGlobal sets the global logger instance.
// This function is thread-safe and should be called during application initialization.
func SetGlobal(logger Logger) {
	mu.Lock()
	defer mu.Unlock()
	globalLogger = logger
}

// getGlobalLogger returns the global logger instance.
// If no logger has been set, it initializes and returns a no-op logger.
// This function is thread-safe for concurrent reads.
func getGlobalLogger() Logger {
	// Acquire a read lock to allow concurrent reads while blocking writes
	mu.RLock()
	defer mu.RUnlock()
	if globalLogger == nil {
		// Set a no-op logger as the default
		globalLogger = NewNoopLogger()
	}
	return globalLogger
}

// Debug logs a debug message using the global logger.
func Debug(msg string, args ...any) {
	getGlobalLogger().Debug(msg, args...)
}

// Info logs an info message using the global logger.
func Info(msg string, args ...any) {
	getGlobalLogger().Info(msg, args...)
}

// Warn logs a warning message using the global logger.
func Warn(msg string, args ...any) {
	getGlobalLogger().Warn(msg, args...)
}

// Error logs an error message using the global logger.
func Error(msg string, args ...any) {
	getGlobalLogger().Error(msg, args...)
}
