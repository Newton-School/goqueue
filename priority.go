package goqueue

import "fmt"

const (
	// MinPriority is the lowest supported task priority.
	MinPriority Priority = 0

	// DefaultPriority is used when no task priority is provided.
	DefaultPriority Priority = 5

	// MaxPriority is the highest supported task priority.
	MaxPriority Priority = 9
)

// Priority describes task importance within a queue.
type Priority int

// ValidatePriority verifies that priority is in the supported range.
func ValidatePriority(priority Priority) error {
	if priority < MinPriority || priority > MaxPriority {
		return fmt.Errorf("%w: priority must be between %d and %d", ErrInvalidPriority, MinPriority, MaxPriority)
	}

	return nil
}
