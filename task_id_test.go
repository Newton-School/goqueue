package goqueue

import (
	"errors"
	"testing"
)

func TestNewTaskIDReturnsValidUniqueID(t *testing.T) {
	first, err := NewTaskID()
	if err != nil {
		t.Fatalf("NewTaskID returned error: %v", err)
	}

	second, err := NewTaskID()
	if err != nil {
		t.Fatalf("NewTaskID returned error: %v", err)
	}

	if first == second {
		t.Fatalf("NewTaskID returned duplicate IDs: %q", first)
	}
	if err := ValidateTaskID(first.String()); err != nil {
		t.Fatalf("ValidateTaskID rejected generated ID %q: %v", first, err)
	}
}

func TestValidateTaskIDRejectsInvalidID(t *testing.T) {
	err := ValidateTaskID("not a task id")
	if !errors.Is(err, ErrInvalidTaskID) {
		t.Fatalf("ValidateTaskID error = %v, want ErrInvalidTaskID", err)
	}
}
