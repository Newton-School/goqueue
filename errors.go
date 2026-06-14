package goqueue

import (
	"errors"

	"github.com/Newton-School/goqueue/task"
)

var (
	// ErrMissingRedisURL is returned when an app is created without Redis.
	ErrMissingRedisURL = errors.New("goqueue: missing redis url")

	// ErrInvalidRedisURL is returned when RedisURL is malformed or unsupported.
	ErrInvalidRedisURL = errors.New("goqueue: invalid redis url")

	// ErrInvalidQueueName is returned when a queue name is not Redis-key safe.
	ErrInvalidQueueName = task.ErrInvalidQueueName

	// ErrInvalidNamespace is returned when Namespace is not Redis-key safe.
	ErrInvalidNamespace = errors.New("goqueue: invalid namespace")
)
