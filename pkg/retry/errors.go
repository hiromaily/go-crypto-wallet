package retry

import "errors"

var (
	// ErrInvalidMaxRetries is returned when MaxRetries is negative.
	ErrInvalidMaxRetries = errors.New("max retries must be non-negative")

	// ErrInvalidInitialBackoff is returned when InitialBackoff is negative.
	ErrInvalidInitialBackoff = errors.New("initial backoff must be non-negative")

	// ErrInvalidMaxBackoff is returned when MaxBackoff is negative.
	ErrInvalidMaxBackoff = errors.New("max backoff must be non-negative")

	// ErrInvalidMultiplier is returned when Multiplier is less than 1.0.
	ErrInvalidMultiplier = errors.New("multiplier must be at least 1.0")

	// ErrInitialBackoffExceedsMax is returned when InitialBackoff exceeds MaxBackoff.
	ErrInitialBackoffExceedsMax = errors.New("initial backoff must not exceed max backoff")
)
