package task

import "fmt"

const (
	TaskPending      TaskState = "PENDING"
	TaskScheduled    TaskState = "SCHEDULED"
	TaskReceived     TaskState = "RECEIVED"
	TaskStarted      TaskState = "STARTED"
	TaskRetrying     TaskState = "RETRYING"
	TaskSucceeded    TaskState = "SUCCEEDED"
	TaskFailed       TaskState = "FAILED"
	TaskRevoked      TaskState = "REVOKED"
	TaskExpired      TaskState = "EXPIRED"
	TaskDeadLettered TaskState = "DEAD_LETTERED"
)

// TaskState describes the lifecycle state of a task invocation.
type TaskState string

// String returns the task state as a string.
func (s TaskState) String() string {
	return string(s)
}

// Terminal reports whether no further worker execution should happen.
func (s TaskState) Terminal() bool {
	switch s {
	case TaskSucceeded, TaskFailed, TaskRevoked, TaskExpired, TaskDeadLettered:
		return true
	default:
		return false
	}
}

// ValidateTaskState verifies that state is one of the supported lifecycle states.
func ValidateTaskState(state TaskState) error {
	switch state {
	case TaskPending, TaskScheduled, TaskReceived, TaskStarted, TaskRetrying,
		TaskSucceeded, TaskFailed, TaskRevoked, TaskExpired, TaskDeadLettered:
		return nil
	default:
		return fmt.Errorf("%w: %q", ErrInvalidTaskState, state)
	}
}
