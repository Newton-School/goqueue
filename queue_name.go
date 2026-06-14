package goqueue

import "fmt"

// QueueName identifies a task queue.
type QueueName string

// String returns the queue name as a string.
func (n QueueName) String() string {
	return string(n)
}

// ValidateQueueName verifies that name can be safely used as a queue key part.
func ValidateQueueName(name string) error {
	if !validName(name) {
		return fmt.Errorf("%w: queue names must use 1-%d characters from [A-Za-z0-9._:-]", ErrInvalidQueueName, maxNameLength)
	}

	return nil
}
