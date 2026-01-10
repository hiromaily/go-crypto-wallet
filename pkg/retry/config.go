package retry

import "time"

// Config holds retry configuration.
type Config struct {
	// MaxRetries is the maximum number of retry attempts.
	// A value of 0 means no retries (only one attempt).
	MaxRetries int

	// InitialBackoff is the initial backoff duration before the first retry.
	InitialBackoff time.Duration

	// MaxBackoff is the maximum backoff duration (cap).
	// If the calculated backoff exceeds this value, it will be capped.
	MaxBackoff time.Duration

	// Multiplier is the backoff multiplier for exponential backoff.
	// The backoff duration is multiplied by this value after each retry.
	// Default is 2.0 (exponential backoff).
	Multiplier float64

	// Jitter adds random jitter to backoff durations to prevent thundering herd.
	// When enabled, the actual backoff will be a random value between
	// 0 and the calculated backoff duration.
	Jitter bool
}

// DefaultConfig returns a default retry configuration suitable for most use cases.
// Default values:
//   - MaxRetries: 3
//   - InitialBackoff: 100ms
//   - MaxBackoff: 10s
//   - Multiplier: 2.0 (exponential)
//   - Jitter: true
func DefaultConfig() Config {
	return Config{
		MaxRetries:     3,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     10 * time.Second,
		Multiplier:     2.0,
		Jitter:         true,
	}
}

// Validate checks if the configuration is valid and returns an error if not.
func (c Config) Validate() error {
	if c.MaxRetries < 0 {
		return ErrInvalidMaxRetries
	}
	if c.InitialBackoff < 0 {
		return ErrInvalidInitialBackoff
	}
	if c.MaxBackoff < 0 {
		return ErrInvalidMaxBackoff
	}
	if c.Multiplier < 1.0 {
		return ErrInvalidMultiplier
	}
	if c.MaxBackoff > 0 && c.InitialBackoff > c.MaxBackoff {
		return ErrInitialBackoffExceedsMax
	}
	return nil
}
