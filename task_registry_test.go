package goqueue

import (
	"errors"
	"testing"
)

func TestTaskRegistryRegisterStoresHandler(t *testing.T) {
	registry := NewTaskRegistry()
	handler := TaskHandlerFunc(func(HandlerContext, TaskPayload) (TaskResult, error) {
		return SucceededResult(nil), nil
	})

	if err := registry.Register("email.send", handler); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
}

func TestTaskRegistryRegisterRejectsDuplicateName(t *testing.T) {
	registry := NewTaskRegistry()
	handler := TaskHandlerFunc(func(HandlerContext, TaskPayload) (TaskResult, error) {
		return SucceededResult(nil), nil
	})

	if err := registry.Register("email.send", handler); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	err := registry.Register("email.send", handler)
	if !errors.Is(err, ErrDuplicateTask) {
		t.Fatalf("Register error = %v, want ErrDuplicateTask", err)
	}
}

func TestTaskRegistryRegisterRejectsNilHandler(t *testing.T) {
	registry := NewTaskRegistry()

	err := registry.Register("email.send", nil)
	if !errors.Is(err, ErrInvalidTaskHandler) {
		t.Fatalf("Register error = %v, want ErrInvalidTaskHandler", err)
	}
}
