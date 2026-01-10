// Package retry provides reusable retry functionality with exponential backoff.
//
// This package supports configurable retry strategies with exponential backoff,
// optional jitter, context cancellation, and generic return types.
//
// Example usage:
//
//	// Simple retry with default configuration
//	err := retry.Retry(ctx, retry.DefaultConfig(), func() error {
//	    return someOperation()
//	})
//
//	// Retry with custom configuration
//	cfg := retry.Config{
//	    MaxRetries:     5,
//	    InitialBackoff: 200 * time.Millisecond,
//	    MaxBackoff:     5 * time.Second,
//	    Multiplier:     2.0,
//	    Jitter:         true,
//	}
//	err := retry.Retry(ctx, cfg, func() error {
//	    return someOperation()
//	})
//
//	// Retry with result
//	result, err := retry.RetryWithResult(ctx, retry.DefaultConfig(), func() (*MyResult, error) {
//	    return fetchData()
//	})
//
//	// Conditional retry (only retry specific errors)
//	err := retry.RetryIf(ctx, cfg, func() error {
//	    return someOperation()
//	}, func(err error) bool {
//	    return errors.Is(err, ErrTemporary)
//	})
package retry

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"
)

// IsRetryable is a function type to determine if an error should be retried.
type IsRetryable func(error) bool

// Retry executes the operation with retry logic using exponential backoff.
// It returns the last error encountered if all retries are exhausted.
// The operation is retried based on the provided configuration.
//
// Parameters:
//   - ctx: Context for cancellation
//   - cfg: Retry configuration
//   - operation: Function to execute (returns error)
//
// Returns:
//   - error: nil if operation succeeds, last error if all retries are exhausted
func Retry(ctx context.Context, cfg Config, operation func() error) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid retry config: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		// Check context cancellation before attempt
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context cancelled before attempt %d: %w", attempt, err)
		}

		// Execute operation
		lastErr = operation()
		if lastErr == nil {
			return nil // Success
		}

		// Don't sleep after the last attempt
		if attempt >= cfg.MaxRetries {
			break
		}

		// Calculate backoff duration
		backoff := calculateBackoff(attempt, cfg)

		// Sleep with context cancellation support
		if err := sleepWithContext(ctx, backoff); err != nil {
			return fmt.Errorf("context cancelled during backoff after attempt %d: %w", attempt, err)
		}
	}

	return fmt.Errorf("failed after %d attempts: %w", cfg.MaxRetries+1, lastErr)
}

// RetryWithResult executes the operation and returns the result with retry logic.
// It supports generic return types using Go generics.
//
// Parameters:
//   - ctx: Context for cancellation
//   - cfg: Retry configuration
//   - operation: Function to execute (returns T and error)
//
// Returns:
//   - T: Result of the operation (zero value if failed)
//   - error: nil if operation succeeds, last error if all retries are exhausted
func RetryWithResult[T any](ctx context.Context, cfg Config, operation func() (T, error)) (T, error) {
	if err := cfg.Validate(); err != nil {
		var zero T
		return zero, fmt.Errorf("invalid retry config: %w", err)
	}

	var lastErr error
	var result T

	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		// Check context cancellation before attempt
		if err := ctx.Err(); err != nil {
			var zero T
			return zero, fmt.Errorf("context cancelled before attempt %d: %w", attempt, err)
		}

		// Execute operation
		result, lastErr = operation()
		if lastErr == nil {
			return result, nil // Success
		}

		// Don't sleep after the last attempt
		if attempt >= cfg.MaxRetries {
			break
		}

		// Calculate backoff duration
		backoff := calculateBackoff(attempt, cfg)

		// Sleep with context cancellation support
		if err := sleepWithContext(ctx, backoff); err != nil {
			var zero T
			return zero, fmt.Errorf("context cancelled during backoff after attempt %d: %w", attempt, err)
		}
	}

	var zero T
	return zero, fmt.Errorf("failed after %d attempts: %w", cfg.MaxRetries+1, lastErr)
}

// RetryIf executes the operation with conditional retry based on error type.
// Only errors for which isRetryable returns true will be retried.
// Other errors will be returned immediately without retry.
//
// Parameters:
//   - ctx: Context for cancellation
//   - cfg: Retry configuration
//   - operation: Function to execute (returns error)
//   - isRetryable: Function to determine if an error should be retried
//
// Returns:
//   - error: nil if operation succeeds, error if operation fails or all retries are exhausted
func RetryIf(ctx context.Context, cfg Config, operation func() error, isRetryable IsRetryable) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid retry config: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		// Check context cancellation before attempt
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context cancelled before attempt %d: %w", attempt, err)
		}

		// Execute operation
		lastErr = operation()
		if lastErr == nil {
			return nil // Success
		}

		// Check if error is retryable
		if !isRetryable(lastErr) {
			return fmt.Errorf("non-retryable error on attempt %d: %w", attempt, lastErr)
		}

		// Don't sleep after the last attempt
		if attempt >= cfg.MaxRetries {
			break
		}

		// Calculate backoff duration
		backoff := calculateBackoff(attempt, cfg)

		// Sleep with context cancellation support
		if err := sleepWithContext(ctx, backoff); err != nil {
			return fmt.Errorf("context cancelled during backoff after attempt %d: %w", attempt, err)
		}
	}

	return fmt.Errorf("failed after %d attempts: %w", cfg.MaxRetries+1, lastErr)
}

// calculateBackoff calculates the backoff duration for the given attempt.
// It uses exponential backoff with optional jitter and max backoff cap.
func calculateBackoff(attempt int, cfg Config) time.Duration {
	// Calculate exponential backoff: InitialBackoff * (Multiplier ^ attempt)
	backoff := float64(cfg.InitialBackoff)
	for range attempt {
		backoff *= cfg.Multiplier
	}

	// Convert to duration
	duration := time.Duration(backoff)

	// Apply max backoff cap
	if cfg.MaxBackoff > 0 && duration > cfg.MaxBackoff {
		duration = cfg.MaxBackoff
	}

	// Apply jitter if enabled
	if cfg.Jitter && duration > 0 {
		// Random value between 0 and duration
		// Using math/rand/v2 is acceptable here as this is for backoff jitter,
		// not cryptographic purposes (nosec G404)
		duration = time.Duration(rand.Int64N(int64(duration))) //nolint:gosec
	}

	return duration
}

// sleepWithContext sleeps for the specified duration or until context is cancelled.
func sleepWithContext(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
