package scheduler

import (
	"fmt"
	"time"
)

const (
	// ScheduleKindInterval identifies fixed interval schedules.
	ScheduleKindInterval = "interval"
)

// IntervalSchedule runs a periodic task after each fixed interval.
type IntervalSchedule struct {
	Interval time.Duration
}

// Every creates an interval schedule.
func Every(interval time.Duration) IntervalSchedule {
	return IntervalSchedule{Interval: interval}
}

// Kind returns the stable storage kind for interval schedules.
func (s IntervalSchedule) Kind() string {
	return ScheduleKindInterval
}

// Validate verifies that the interval can produce future run times.
func (s IntervalSchedule) Validate() error {
	if s.Interval <= 0 {
		return fmt.Errorf("%w: interval must be positive", ErrInvalidSchedule)
	}

	return nil
}

// Next returns the next run time after now in UTC.
func (s IntervalSchedule) Next(now time.Time) time.Time {
	return now.UTC().Add(s.Interval)
}
