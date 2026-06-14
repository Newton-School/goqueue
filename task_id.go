package goqueue

import (
	"crypto/rand"
	"fmt"
)

// TaskID uniquely identifies a task invocation.
type TaskID string

// NewTaskID returns a random RFC 4122 version 4 task ID.
func NewTaskID() (TaskID, error) {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", fmt.Errorf("goqueue: generate task id: %w", err)
	}

	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80

	return TaskID(fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		bytes[0:4],
		bytes[4:6],
		bytes[6:8],
		bytes[8:10],
		bytes[10:16],
	)), nil
}

// String returns the task ID as a string.
func (id TaskID) String() string {
	return string(id)
}

// ValidateTaskID verifies that id is an RFC 4122 UUID string.
func ValidateTaskID(id string) error {
	if len(id) != 36 {
		return fmt.Errorf("%w: expected RFC 4122 UUID string", ErrInvalidTaskID)
	}

	for index, char := range id {
		switch index {
		case 8, 13, 18, 23:
			if char != '-' {
				return fmt.Errorf("%w: expected hyphen at index %d", ErrInvalidTaskID, index)
			}
			continue
		}

		if !isHex(char) {
			return fmt.Errorf("%w: expected hexadecimal character at index %d", ErrInvalidTaskID, index)
		}
	}

	return nil
}

func isHex(char rune) bool {
	return (char >= '0' && char <= '9') ||
		(char >= 'a' && char <= 'f') ||
		(char >= 'A' && char <= 'F')
}
