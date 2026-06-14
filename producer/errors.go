package producer

import "errors"

var (
	// ErrNilBackend is returned when a producer is created without storage.
	ErrNilBackend = errors.New("goqueue producer: nil backend")

	// ErrMissingTaskName is returned when a task name is required but missing.
	ErrMissingTaskName = errors.New("goqueue producer: task name is required")

	// ErrMissingApplyOption is returned when apply options are invalid.
	ErrMissingApplyOption = errors.New("goqueue producer: invalid apply option")
)
