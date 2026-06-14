package backend

import "errors"

var (
	// ErrTaskMessageNotFound is returned when a persisted task message is missing.
	ErrTaskMessageNotFound = errors.New("goqueue backend: task message not found")

	// ErrTaskStateNotFound is returned when task state has not been written.
	ErrTaskStateNotFound = errors.New("goqueue backend: task state not found")

	// ErrTaskResultNotFound is returned when task result has not been written.
	ErrTaskResultNotFound = errors.New("goqueue backend: task result not found")

	// ErrConsumerGroupNotFound is returned when a queue consumer group is missing.
	ErrConsumerGroupNotFound = errors.New("goqueue backend: consumer group not found")

	// ErrInvalidBackendRequest is returned when a backend request is incomplete.
	ErrInvalidBackendRequest = errors.New("goqueue backend: invalid request")
)
