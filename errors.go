package goqueue

import "errors"

var (
	// ErrMissingRedisURL is returned when an app is created without Redis.
	ErrMissingRedisURL = errors.New("goqueue: missing redis url")

	// ErrInvalidRedisURL is returned when RedisURL is malformed or unsupported.
	ErrInvalidRedisURL = errors.New("goqueue: invalid redis url")

	// ErrInvalidQueueName is returned when DefaultQueue is not Redis-key safe.
	ErrInvalidQueueName = errors.New("goqueue: invalid queue name")

	// ErrInvalidNamespace is returned when Namespace is not Redis-key safe.
	ErrInvalidNamespace = errors.New("goqueue: invalid namespace")

	// ErrInvalidTaskName is returned when a task name is empty or unsafe.
	ErrInvalidTaskName = errors.New("goqueue: invalid task name")

	// ErrInvalidTaskID is returned when a task ID is not a UUID string.
	ErrInvalidTaskID = errors.New("goqueue: invalid task id")

	// ErrInvalidPriority is returned when a task priority is outside the supported range.
	ErrInvalidPriority = errors.New("goqueue: invalid priority")

	// ErrInvalidTaskState is returned when a task state is not recognized.
	ErrInvalidTaskState = errors.New("goqueue: invalid task state")
)
