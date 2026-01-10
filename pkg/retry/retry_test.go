package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

var (
	errTemporary = errors.New("temporary error")
	errPermanent = errors.New("permanent error")
)

func TestRetry_Success(t *testing.T) {
	t.Parallel()

	cfg := Config{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
		Jitter:         false,
	}

	attemptCount := 0
	operation := func(attempt uint) error {
		attemptCount++
		if attempt != uint(attemptCount) {
			t.Errorf("expected attempt %d, got %d", attemptCount, attempt)
		}
		return nil // Success on first attempt
	}

	err := Retry(context.Background(), cfg, operation)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if attemptCount != 1 {
		t.Errorf("expected 1 attempt, got %d", attemptCount)
	}
}

func TestRetry_SuccessAfterRetries(t *testing.T) {
	t.Parallel()

	cfg := Config{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
		Jitter:         false,
	}

	attemptCount := 0
	operation := func(attempt uint) error {
		attemptCount++
		if attempt != uint(attemptCount) {
			t.Errorf("expected attempt %d, got %d", attemptCount, attempt)
		}
		if attemptCount < 3 {
			return errTemporary // Fail first 2 attempts
		}
		return nil // Success on 3rd attempt
	}

	err := Retry(context.Background(), cfg, operation)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if attemptCount != 3 {
		t.Errorf("expected 3 attempts, got %d", attemptCount)
	}
}

func TestRetry_MaxRetriesExceeded(t *testing.T) {
	t.Parallel()

	cfg := Config{
		MaxRetries:     2,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
		Jitter:         false,
	}

	attemptCount := 0
	operation := func(attempt uint) error {
		attemptCount++
		if attempt != uint(attemptCount) {
			t.Errorf("expected attempt %d, got %d", attemptCount, attempt)
		}
		return errTemporary // Always fail
	}

	err := Retry(context.Background(), cfg, operation)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// MaxRetries=2 means 3 total attempts (initial + 2 retries)
	expectedAttempts := cfg.MaxRetries + 1
	if attemptCount != expectedAttempts {
		t.Errorf("expected %d attempts, got %d", expectedAttempts, attemptCount)
	}

	if !errors.Is(err, errTemporary) {
		t.Errorf("expected error to wrap errTemporary, got: %v", err)
	}
}

func TestRetry_ContextCancellation(t *testing.T) {
	t.Parallel()

	cfg := Config{
		MaxRetries:     5,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     1 * time.Second,
		Multiplier:     2.0,
		Jitter:         false,
	}

	ctx, cancel := context.WithCancel(context.Background())

	attemptCount := 0
	operation := func(attempt uint) error {
		attemptCount++
		if attempt != uint(attemptCount) {
			t.Errorf("expected attempt %d, got %d", attemptCount, attempt)
		}
		if attemptCount == 2 {
			cancel() // Cancel context on 2nd attempt
		}
		return errTemporary
	}

	err := Retry(ctx, cfg, operation)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got: %v", err)
	}

	// Should stop retrying after context cancellation
	if attemptCount > 3 {
		t.Errorf("expected at most 3 attempts before cancellation, got %d", attemptCount)
	}
}

func TestRetry_ContextTimeout(t *testing.T) {
	t.Parallel()

	cfg := Config{
		MaxRetries:     10,
		InitialBackoff: 50 * time.Millisecond,
		MaxBackoff:     200 * time.Millisecond,
		Multiplier:     2.0,
		Jitter:         false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	attemptCount := 0
	operation := func(attempt uint) error {
		attemptCount++
		if attempt != uint(attemptCount) {
			t.Errorf("expected attempt %d, got %d", attemptCount, attempt)
		}
		return errTemporary
	}

	err := Retry(ctx, cfg, operation)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context.DeadlineExceeded error, got: %v", err)
	}
}

func TestRetry_ExponentialBackoff(t *testing.T) {
	t.Parallel()

	cfg := Config{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     0, // No cap
		Multiplier:     2.0,
		Jitter:         false,
	}

	attemptTimes := []time.Time{}
	operation := func(attempt uint) error {
		attemptTimes = append(attemptTimes, time.Now())
		return errTemporary
	}

	_ = Retry(context.Background(), cfg, operation)

	// Verify exponential backoff timing (within tolerance)
	if len(attemptTimes) != cfg.MaxRetries+1 {
		t.Fatalf("expected %d attempts, got %d", cfg.MaxRetries+1, len(attemptTimes))
	}

	// Check backoff durations between attempts
	expectedBackoffs := []time.Duration{
		10 * time.Millisecond, // After 1st attempt
		20 * time.Millisecond, // After 2nd attempt
		40 * time.Millisecond, // After 3rd attempt
	}

	tolerance := 5 * time.Millisecond
	for i := range expectedBackoffs {
		actualBackoff := attemptTimes[i+1].Sub(attemptTimes[i])
		expected := expectedBackoffs[i]

		if actualBackoff < expected-tolerance || actualBackoff > expected+tolerance {
			t.Errorf("attempt %d: expected backoff ~%v, got %v", i+1, expected, actualBackoff)
		}
	}
}

func TestRetry_MaxBackoffCap(t *testing.T) {
	t.Parallel()

	cfg := Config{
		MaxRetries:     5,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     50 * time.Millisecond, // Cap at 50ms
		Multiplier:     2.0,
		Jitter:         false,
	}

	attemptTimes := []time.Time{}
	operation := func(attempt uint) error {
		attemptTimes = append(attemptTimes, time.Now())
		return errTemporary
	}

	_ = Retry(context.Background(), cfg, operation)

	// After attempt 3, backoff should be capped at MaxBackoff
	// Attempt 1: 10ms, Attempt 2: 20ms, Attempt 3: 40ms, Attempt 4+: 50ms (capped)
	for i := 3; i < len(attemptTimes)-1; i++ {
		actualBackoff := attemptTimes[i+1].Sub(attemptTimes[i])
		maxWithTolerance := cfg.MaxBackoff + 10*time.Millisecond

		if actualBackoff > maxWithTolerance {
			t.Errorf("attempt %d: backoff %v exceeds max backoff %v", i+1, actualBackoff, cfg.MaxBackoff)
		}
	}
}

func TestRetry_InvalidConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "negative max retries",
			config: Config{
				MaxRetries:     -1,
				InitialBackoff: 10 * time.Millisecond,
				MaxBackoff:     100 * time.Millisecond,
				Multiplier:     2.0,
			},
		},
		{
			name: "negative initial backoff",
			config: Config{
				MaxRetries:     3,
				InitialBackoff: -10 * time.Millisecond,
				MaxBackoff:     100 * time.Millisecond,
				Multiplier:     2.0,
			},
		},
		{
			name: "invalid multiplier",
			config: Config{
				MaxRetries:     3,
				InitialBackoff: 10 * time.Millisecond,
				MaxBackoff:     100 * time.Millisecond,
				Multiplier:     0.5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Retry(context.Background(), tt.config, func(attempt uint) error { return nil })
			if err == nil {
				t.Error("expected validation error, got nil")
			}
		})
	}
}

func TestRetryWithResult_Success(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig()

	attemptCount := 0
	operation := func(attempt uint) (string, error) {
		attemptCount++
		if attempt != uint(attemptCount) {
			t.Errorf("expected attempt %d, got %d", attemptCount, attempt)
		}
		return "success", nil
	}

	result, err := RetryWithResult(context.Background(), cfg, operation)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result != "success" {
		t.Errorf("expected result 'success', got '%s'", result)
	}

	if attemptCount != 1 {
		t.Errorf("expected 1 attempt, got %d", attemptCount)
	}
}

func TestRetryWithResult_SuccessAfterRetries(t *testing.T) {
	t.Parallel()

	cfg := Config{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
		Jitter:         false,
	}

	attemptCount := 0
	operation := func(attempt uint) (int, error) {
		attemptCount++
		if attempt != uint(attemptCount) {
			t.Errorf("expected attempt %d, got %d", attemptCount, attempt)
		}
		if attemptCount < 3 {
			return 0, errTemporary
		}
		return 42, nil
	}

	result, err := RetryWithResult(context.Background(), cfg, operation)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result != 42 {
		t.Errorf("expected result 42, got %d", result)
	}

	if attemptCount != 3 {
		t.Errorf("expected 3 attempts, got %d", attemptCount)
	}
}

func TestRetryWithResult_MaxRetriesExceeded(t *testing.T) {
	t.Parallel()

	cfg := Config{
		MaxRetries:     2,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
		Jitter:         false,
	}

	attemptCount := 0
	operation := func(attempt uint) (*int, error) {
		attemptCount++
		if attempt != uint(attemptCount) {
			t.Errorf("expected attempt %d, got %d", attemptCount, attempt)
		}
		return nil, errTemporary
	}

	result, err := RetryWithResult(context.Background(), cfg, operation)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	expectedAttempts := cfg.MaxRetries + 1
	if attemptCount != expectedAttempts {
		t.Errorf("expected %d attempts, got %d", expectedAttempts, attemptCount)
	}
}

func TestRetryWithResult_ContextCancellation(t *testing.T) {
	t.Parallel()

	cfg := Config{
		MaxRetries:     5,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     1 * time.Second,
		Multiplier:     2.0,
		Jitter:         false,
	}

	ctx, cancel := context.WithCancel(context.Background())

	attemptCount := 0
	operation := func(attempt uint) (string, error) {
		attemptCount++
		if attempt != uint(attemptCount) {
			t.Errorf("expected attempt %d, got %d", attemptCount, attempt)
		}
		if attemptCount == 2 {
			cancel()
		}
		return "", errTemporary
	}

	result, err := RetryWithResult(ctx, cfg, operation)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if result != "" {
		t.Errorf("expected empty result, got '%s'", result)
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got: %v", err)
	}
}

func TestRetryIf_NonRetryableError(t *testing.T) {
	t.Parallel()

	cfg := Config{
		MaxRetries:     5,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
		Jitter:         false,
	}

	isRetryable := func(err error) bool {
		return errors.Is(err, errTemporary)
	}

	attemptCount := 0
	operation := func(attempt uint) error {
		attemptCount++
		if attempt != uint(attemptCount) {
			t.Errorf("expected attempt %d, got %d", attemptCount, attempt)
		}
		return errPermanent // Non-retryable error
	}

	err := RetryIf(context.Background(), cfg, operation, isRetryable)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Should fail immediately without retry
	if attemptCount != 1 {
		t.Errorf("expected 1 attempt, got %d", attemptCount)
	}

	if !errors.Is(err, errPermanent) {
		t.Errorf("expected error to wrap errPermanent, got: %v", err)
	}
}

func TestRetryIf_RetryableError(t *testing.T) {
	t.Parallel()

	cfg := Config{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
		Jitter:         false,
	}

	isRetryable := func(err error) bool {
		return errors.Is(err, errTemporary)
	}

	attemptCount := 0
	operation := func(attempt uint) error {
		attemptCount++
		if attempt != uint(attemptCount) {
			t.Errorf("expected attempt %d, got %d", attemptCount, attempt)
		}
		if attemptCount < 3 {
			return errTemporary // Retryable error
		}
		return nil // Success on 3rd attempt
	}

	err := RetryIf(context.Background(), cfg, operation, isRetryable)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if attemptCount != 3 {
		t.Errorf("expected 3 attempts, got %d", attemptCount)
	}
}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig()

	if cfg.MaxRetries != 3 {
		t.Errorf("expected MaxRetries=3, got %d", cfg.MaxRetries)
	}
	if cfg.InitialBackoff != 100*time.Millisecond {
		t.Errorf("expected InitialBackoff=100ms, got %v", cfg.InitialBackoff)
	}
	if cfg.MaxBackoff != 10*time.Second {
		t.Errorf("expected MaxBackoff=10s, got %v", cfg.MaxBackoff)
	}
	if cfg.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", cfg.Multiplier)
	}
	if !cfg.Jitter {
		t.Error("expected Jitter=true, got false")
	}

	if err := cfg.Validate(); err != nil {
		t.Errorf("default config should be valid, got error: %v", err)
	}
}

func TestConfigValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		config      Config
		expectError error
	}{
		{
			name: "valid config",
			config: Config{
				MaxRetries:     3,
				InitialBackoff: 100 * time.Millisecond,
				MaxBackoff:     10 * time.Second,
				Multiplier:     2.0,
				Jitter:         true,
			},
			expectError: nil,
		},
		{
			name: "negative max retries",
			config: Config{
				MaxRetries:     -1,
				InitialBackoff: 100 * time.Millisecond,
				MaxBackoff:     10 * time.Second,
				Multiplier:     2.0,
			},
			expectError: ErrInvalidMaxRetries,
		},
		{
			name: "negative initial backoff",
			config: Config{
				MaxRetries:     3,
				InitialBackoff: -100 * time.Millisecond,
				MaxBackoff:     10 * time.Second,
				Multiplier:     2.0,
			},
			expectError: ErrInvalidInitialBackoff,
		},
		{
			name: "negative max backoff",
			config: Config{
				MaxRetries:     3,
				InitialBackoff: 100 * time.Millisecond,
				MaxBackoff:     -10 * time.Second,
				Multiplier:     2.0,
			},
			expectError: ErrInvalidMaxBackoff,
		},
		{
			name: "invalid multiplier",
			config: Config{
				MaxRetries:     3,
				InitialBackoff: 100 * time.Millisecond,
				MaxBackoff:     10 * time.Second,
				Multiplier:     0.5,
			},
			expectError: ErrInvalidMultiplier,
		},
		{
			name: "initial backoff exceeds max",
			config: Config{
				MaxRetries:     3,
				InitialBackoff: 20 * time.Second,
				MaxBackoff:     10 * time.Second,
				Multiplier:     2.0,
			},
			expectError: ErrInitialBackoffExceedsMax,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.config.Validate()
			if tt.expectError == nil {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectError)
				} else if !errors.Is(err, tt.expectError) {
					t.Errorf("expected error %v, got: %v", tt.expectError, err)
				}
			}
		})
	}
}

func TestCalculateBackoff(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		attempt  int
		config   Config
		minValue time.Duration // For jitter, we check if value is at least 0
		maxValue time.Duration // For jitter, we check if value is less than this
	}{
		{
			name:    "first retry, no jitter",
			attempt: 0,
			config: Config{
				InitialBackoff: 100 * time.Millisecond,
				MaxBackoff:     10 * time.Second,
				Multiplier:     2.0,
				Jitter:         false,
			},
			minValue: 100 * time.Millisecond,
			maxValue: 100 * time.Millisecond,
		},
		{
			name:    "second retry, no jitter",
			attempt: 1,
			config: Config{
				InitialBackoff: 100 * time.Millisecond,
				MaxBackoff:     10 * time.Second,
				Multiplier:     2.0,
				Jitter:         false,
			},
			minValue: 200 * time.Millisecond,
			maxValue: 200 * time.Millisecond,
		},
		{
			name:    "capped by max backoff",
			attempt: 10,
			config: Config{
				InitialBackoff: 100 * time.Millisecond,
				MaxBackoff:     1 * time.Second,
				Multiplier:     2.0,
				Jitter:         false,
			},
			minValue: 1 * time.Second,
			maxValue: 1 * time.Second,
		},
		{
			name:    "with jitter",
			attempt: 1,
			config: Config{
				InitialBackoff: 100 * time.Millisecond,
				MaxBackoff:     10 * time.Second,
				Multiplier:     2.0,
				Jitter:         true,
			},
			minValue: 0,
			maxValue: 200 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			backoff := calculateBackoff(tt.attempt, tt.config)

			if !tt.config.Jitter {
				// Without jitter, backoff should be exact
				if backoff != tt.minValue {
					t.Errorf("expected backoff %v, got %v", tt.minValue, backoff)
				}
			} else {
				// With jitter, backoff should be in range [0, maxValue)
				if backoff < tt.minValue || backoff >= tt.maxValue {
					t.Errorf("expected backoff in range [%v, %v), got %v", tt.minValue, tt.maxValue, backoff)
				}
			}
		})
	}
}

func TestSleepWithContext(t *testing.T) {
	t.Parallel()

	t.Run("sleep completes", func(t *testing.T) {
		t.Parallel()

		start := time.Now()
		err := sleepWithContext(context.Background(), 50*time.Millisecond)
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		if elapsed < 50*time.Millisecond {
			t.Errorf("expected sleep duration >= 50ms, got %v", elapsed)
		}
	})

	t.Run("context cancelled", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		err := sleepWithContext(ctx, 1*time.Second)
		if err == nil {
			t.Error("expected error, got nil")
		}

		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled error, got: %v", err)
		}
	})
}
