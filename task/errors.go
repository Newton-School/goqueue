package task

import "errors"

var (
	// ErrInvalidTaskName is returned when a task name is empty or unsafe.
	ErrInvalidTaskName = errors.New("goqueue: invalid task name")

	// ErrInvalidTaskID is returned when a task ID is not a UUID string.
	ErrInvalidTaskID = errors.New("goqueue: invalid task id")

	// ErrInvalidQueueName is returned when a queue name is not Redis-key safe.
	ErrInvalidQueueName = errors.New("goqueue: invalid queue name")

	// ErrInvalidPriority is returned when a task priority is outside the supported range.
	ErrInvalidPriority = errors.New("goqueue: invalid priority")

	// ErrInvalidTaskState is returned when a task state is not recognized.
	ErrInvalidTaskState = errors.New("goqueue: invalid task state")

	// ErrInvalidRetryPolicy is returned when retry settings are unsafe.
	ErrInvalidRetryPolicy = errors.New("goqueue: invalid retry policy")

	// ErrInvalidTaskTiming is returned when task scheduling fields conflict.
	ErrInvalidTaskTiming = errors.New("goqueue: invalid task timing")

	// ErrInvalidPayload is returned when a task payload cannot be encoded or decoded.
	ErrInvalidPayload = errors.New("goqueue: invalid payload")

	// ErrDuplicateTask is returned when a task name is registered more than once.
	ErrDuplicateTask = errors.New("goqueue: duplicate task")

	// ErrInvalidTaskHandler is returned when a task handler is nil.
	ErrInvalidTaskHandler = errors.New("goqueue: invalid task handler")

	// ErrTaskNotFound is returned when a task name has no registered handler.
	ErrTaskNotFound = errors.New("goqueue: task not found")
)
