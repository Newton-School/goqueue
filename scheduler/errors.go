package scheduler

import "errors"

var (
	// ErrInvalidSchedule is returned when a schedule cannot calculate safe run times.
	ErrInvalidSchedule = errors.New("goqueue scheduler: invalid schedule")
)
