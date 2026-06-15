package scheduler

import "errors"

var (
	// ErrNilBackend is returned when a scheduler is created without storage.
	ErrNilBackend = errors.New("goqueue scheduler: backend is nil")

	// ErrInvalidSchedule is returned when a schedule cannot calculate safe run times.
	ErrInvalidSchedule = errors.New("goqueue scheduler: invalid schedule")

	// ErrInvalidPeriodicTask is returned when a periodic task definition is unsafe.
	ErrInvalidPeriodicTask = errors.New("goqueue scheduler: invalid periodic task")

	// ErrInvalidSchedulerOption is returned when scheduler configuration is unsafe.
	ErrInvalidSchedulerOption = errors.New("goqueue scheduler: invalid option")
)
