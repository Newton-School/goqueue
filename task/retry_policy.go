package task

import (
	"fmt"
	"time"
)

const (
	// MaxRetryAttempts bounds retry loops and protects workers from unsafe policies.
	MaxRetryAttempts = 1000

	maxRetryDelay = time.Duration(1<<63 - 1)
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
	if p.MaxAttempts > MaxRetryAttempts {
		return fmt.Errorf("%w: max attempts cannot exceed %d", ErrInvalidRetryPolicy, MaxRetryAttempts)
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

// DelayForAttempt returns the delay before retrying after the given failed attempt.
func (p RetryPolicy) DelayForAttempt(attempt int) time.Duration {
	if attempt <= 0 || p.Backoff <= 0 {
		return 0
	}

	delay := p.Backoff
	for range attempt - 1 {
		if delay > maxRetryDelay/2 {
			if p.MaxBackoff > 0 {
				return p.MaxBackoff
			}
			return maxRetryDelay
		}

		next := delay * 2
		if p.MaxBackoff > 0 && next > p.MaxBackoff {
			return p.MaxBackoff
		}
		delay = next
	}

	if p.MaxBackoff > 0 && delay > p.MaxBackoff {
		return p.MaxBackoff
	}

	return delay
}
