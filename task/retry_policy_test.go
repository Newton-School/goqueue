package task

import (
	"errors"
	"testing"
	"time"
)

func TestDefaultRetryPolicyDoesNotRetryImplicitly(t *testing.T) {
	policy := DefaultRetryPolicy()

	if policy.MaxAttempts != 1 {
		t.Fatalf("MaxAttempts = %d, want 1", policy.MaxAttempts)
	}
	if err := policy.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestRetryPolicyValidateRejectsInvalidAttempts(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 0}

	err := policy.Validate()
	if !errors.Is(err, ErrInvalidRetryPolicy) {
		t.Fatalf("Validate error = %v, want ErrInvalidRetryPolicy", err)
	}
}

func TestRetryPolicyValidateRejectsExcessiveAttempts(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: MaxRetryAttempts + 1}

	err := policy.Validate()
	if !errors.Is(err, ErrInvalidRetryPolicy) {
		t.Fatalf("Validate error = %v, want ErrInvalidRetryPolicy", err)
	}
}

func TestRetryPolicyValidateRejectsInvalidBackoff(t *testing.T) {
	policy := RetryPolicy{
		MaxAttempts: 2,
		Backoff:     10 * time.Second,
		MaxBackoff:  5 * time.Second,
	}

	err := policy.Validate()
	if !errors.Is(err, ErrInvalidRetryPolicy) {
		t.Fatalf("Validate error = %v, want ErrInvalidRetryPolicy", err)
	}
}

func TestRetryPolicyDelayForAttemptCapsOverflow(t *testing.T) {
	policy := RetryPolicy{
		MaxAttempts: 3,
		Backoff:     maxRetryDelay/2 + time.Nanosecond,
	}

	if got := policy.DelayForAttempt(2); got != maxRetryDelay {
		t.Fatalf("DelayForAttempt overflow cap = %s, want %s", got, maxRetryDelay)
	}
}

func TestRetryPolicyDelayForAttemptUsesExponentialBackoff(t *testing.T) {
	policy := RetryPolicy{
		MaxAttempts: 4,
		Backoff:     2 * time.Second,
		MaxBackoff:  5 * time.Second,
	}

	tests := map[int]time.Duration{
		1: 2 * time.Second,
		2: 4 * time.Second,
		3: 5 * time.Second,
	}

	for attempt, want := range tests {
		if got := policy.DelayForAttempt(attempt); got != want {
			t.Fatalf("DelayForAttempt(%d) = %s, want %s", attempt, got, want)
		}
	}
}
