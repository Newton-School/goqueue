package goqueue

import (
	"errors"
	"testing"
)

func TestSucceededResultBuildsSuccessState(t *testing.T) {
	result := SucceededResult("ok")

	if result.State != TaskSucceeded {
		t.Fatalf("State = %s, want %s", result.State, TaskSucceeded)
	}
	if result.Value != "ok" {
		t.Fatalf("Value = %v, want ok", result.Value)
	}
}

func TestFailedResultBuildsFailureState(t *testing.T) {
	result := FailedResult(errors.New("boom"))

	if result.State != TaskFailed {
		t.Fatalf("State = %s, want %s", result.State, TaskFailed)
	}
	if result.Error == "" {
		t.Fatal("Error should contain failure message")
	}
}

func TestTaskResultValidateRejectsUnknownState(t *testing.T) {
	result := TaskResult{State: TaskState("UNKNOWN")}

	err := result.Validate()
	if !errors.Is(err, ErrInvalidTaskState) {
		t.Fatalf("Validate error = %v, want ErrInvalidTaskState", err)
	}
}
