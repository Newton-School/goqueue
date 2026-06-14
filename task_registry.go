package goqueue

import (
	"fmt"
	"sync"
)

// TaskRegistry stores task handlers by task name.
type TaskRegistry struct {
	mu       sync.RWMutex
	handlers map[TaskName]TaskHandler
}

// NewTaskRegistry creates an empty task registry.
func NewTaskRegistry() *TaskRegistry {
	return &TaskRegistry{handlers: make(map[TaskName]TaskHandler)}
}

// Register stores a handler for name.
func (r *TaskRegistry) Register(name TaskName, handler TaskHandler) error {
	if err := ValidateTaskName(name.String()); err != nil {
		return err
	}

	if handler == nil {
		return fmt.Errorf("%w: %s", ErrInvalidTaskHandler, name)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.handlers[name]; exists {
		return fmt.Errorf("%w: %s", ErrDuplicateTask, name)
	}

	r.handlers[name] = handler
	return nil
}
