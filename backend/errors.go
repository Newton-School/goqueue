package backend

import (
	"errors"

	"github.com/Newton-School/goqueue/task"
)

var (
	// ErrTaskMessageNotFound is returned when a persisted task message is missing.
	ErrTaskMessageNotFound = errors.New("goqueue backend: task message not found")

	// ErrDeadLetterNotFound is returned when a dead-letter record is missing.
	ErrDeadLetterNotFound = errors.New("goqueue backend: dead-letter record not found")

	// ErrTaskStateNotFound is returned when task state has not been written.
	ErrTaskStateNotFound = errors.New("goqueue backend: task state not found")

	// ErrTaskResultNotFound is returned when task result has not been written.
	ErrTaskResultNotFound = errors.New("goqueue backend: task result not found")

	// ErrConsumerGroupNotFound is returned when a queue consumer group is missing.
	ErrConsumerGroupNotFound = errors.New("goqueue backend: consumer group not found")

	// ErrPeriodicTaskLeaseLost is returned when a scheduler no longer owns a periodic lease.
	ErrPeriodicTaskLeaseLost = errors.New("goqueue backend: periodic task lease lost")

	// ErrPeriodicTaskNotFound is returned when a periodic task definition is missing.
	ErrPeriodicTaskNotFound = errors.New("goqueue backend: periodic task not found")

	// ErrInvalidBackendRequest is returned when a backend request is incomplete.
	ErrInvalidBackendRequest = errors.New("goqueue backend: invalid request")

	// ErrInvalidQueueName is returned when a backend queue request has an unsafe queue name.
	ErrInvalidQueueName = task.ErrInvalidQueueName
)
