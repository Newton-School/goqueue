package task

import (
	"fmt"
	"time"
)

// TaskTiming describes when a task may be executed and when it expires.
type TaskTiming struct {
	ETA       time.Time
	ExpiresAt time.Time
}

// TaskTimingFromCountdown returns timing with ETA set relative to now.
func TaskTimingFromCountdown(now time.Time, countdown time.Duration) (TaskTiming, error) {
	if countdown < 0 {
		return TaskTiming{}, fmt.Errorf("%w: countdown cannot be negative", ErrInvalidTaskTiming)
	}

	return TaskTiming{ETA: now.Add(countdown)}, nil
}

// Scheduled reports whether the task has a future execution timestamp.
func (t TaskTiming) Scheduled() bool {
	return !t.ETA.IsZero()
}

// Validate verifies that timing fields are internally consistent.
func (t TaskTiming) Validate() error {
	if !t.ETA.IsZero() && !t.ExpiresAt.IsZero() && t.ExpiresAt.Before(t.ETA) {
		return fmt.Errorf("%w: expiration cannot be before ETA", ErrInvalidTaskTiming)
	}

	return nil
}
