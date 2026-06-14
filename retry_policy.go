package goqueue

import (
	"fmt"
	"time"
)

// RetryPolicy controls how many times a task may be attempted.
type RetryPolicy struct {
	MaxAttempts int
	Backoff     time.Duration
	MaxBackoff  time.Duration
}

// DefaultRetryPolicy returns the safe default retry policy.
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{MaxAttempts: 1}
}

// Validate verifies retry settings are bounded and internally consistent.
func (p RetryPolicy) Validate() error {
	if p.MaxAttempts < 1 {
		return fmt.Errorf("%w: max attempts must be at least 1", ErrInvalidRetryPolicy)
	}

	if p.Backoff < 0 {
		return fmt.Errorf("%w: backoff cannot be negative", ErrInvalidRetryPolicy)
	}

	if p.MaxBackoff < 0 {
		return fmt.Errorf("%w: max backoff cannot be negative", ErrInvalidRetryPolicy)
	}

	if p.MaxBackoff > 0 && p.Backoff > p.MaxBackoff {
		return fmt.Errorf("%w: backoff cannot exceed max backoff", ErrInvalidRetryPolicy)
	}

	return nil
}
