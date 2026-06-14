package worker

import "errors"

var (
	// ErrNilWorker is returned when worker instance is nil.
	ErrNilWorker = errors.New("goqueue worker: nil worker")

	// ErrNilBackend is returned when a worker is created without a queue backend.
	ErrNilBackend = errors.New("goqueue worker: nil backend")

	// ErrNilTaskRegistry is returned when no task registry is available.
	ErrNilTaskRegistry = errors.New("goqueue worker: nil task registry")

	// ErrInvalidWorkerOption is returned when required worker configuration is invalid.
	ErrInvalidWorkerOption = errors.New("goqueue worker: invalid option")
)
