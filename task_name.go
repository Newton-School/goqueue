package goqueue

import "fmt"

// TaskName identifies a registered task handler.
type TaskName string

// String returns the task name as a string.
func (n TaskName) String() string {
	return string(n)
}

// ValidateTaskName verifies that name can be safely used in task messages.
func ValidateTaskName(name string) error {
	if !validName(name) {
		return fmt.Errorf("%w: task names must use 1-%d characters from [A-Za-z0-9._:-]", ErrInvalidTaskName, maxNameLength)
	}

	return nil
}
