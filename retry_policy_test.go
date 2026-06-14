package goqueue

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
