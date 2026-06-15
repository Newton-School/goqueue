package scheduler

import "errors"

var (
	// ErrInvalidSchedule is returned when a schedule cannot calculate safe run times.
	ErrInvalidSchedule = errors.New("goqueue scheduler: invalid schedule")

	// ErrInvalidPeriodicTask is returned when a periodic task definition is unsafe.
	ErrInvalidPeriodicTask = errors.New("goqueue scheduler: invalid periodic task")
)
