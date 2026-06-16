package inspect

import "errors"

var (
	ErrNilInspector       = errors.New("inspect: inspector is nil")
	ErrInspectorBackend   = errors.New("inspect: inspector backend is nil")
	ErrEmptyQueueName     = errors.New("inspect: queue name is required")
	ErrInvalidDeadLetters = errors.New("inspect: invalid dead-letter request")
)
