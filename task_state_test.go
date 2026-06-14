package goqueue

import (
	"errors"
	"testing"
)

func TestTaskStateTerminalReportsTerminalStates(t *testing.T) {
	for _, state := range []TaskState{TaskSucceeded, TaskFailed, TaskRevoked, TaskExpired, TaskDeadLettered} {
		if !state.Terminal() {
			t.Fatalf("%s should be terminal", state)
		}
	}
}

func TestTaskStateTerminalReportsActiveStates(t *testing.T) {
	for _, state := range []TaskState{TaskPending, TaskScheduled, TaskReceived, TaskStarted, TaskRetrying} {
		if state.Terminal() {
			t.Fatalf("%s should not be terminal", state)
		}
	}
}

func TestValidateTaskStateRejectsUnknownState(t *testing.T) {
	err := ValidateTaskState(TaskState("lost"))
	if !errors.Is(err, ErrInvalidTaskState) {
		t.Fatalf("ValidateTaskState error = %v, want ErrInvalidTaskState", err)
	}
}
