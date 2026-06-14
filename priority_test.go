package goqueue

import (
	"errors"
	"testing"
)

func TestValidatePriorityAcceptsSupportedRange(t *testing.T) {
	for _, priority := range []Priority{MinPriority, DefaultPriority, MaxPriority} {
		if err := ValidatePriority(priority); err != nil {
			t.Fatalf("ValidatePriority(%d) returned error: %v", priority, err)
		}
	}
}

func TestValidatePriorityRejectsUnsupportedRange(t *testing.T) {
	err := ValidatePriority(MaxPriority + 1)
	if !errors.Is(err, ErrInvalidPriority) {
		t.Fatalf("ValidatePriority error = %v, want ErrInvalidPriority", err)
	}
}
