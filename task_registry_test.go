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

func TestTaskRegistryLookupReturnsRegisteredHandler(t *testing.T) {
	registry := NewTaskRegistry()
	handler := TaskHandlerFunc(func(HandlerContext, TaskPayload) (TaskResult, error) {
		return SucceededResult("ok"), nil
	})
	if err := registry.Register("email.send", handler); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	got, err := registry.Lookup("email.send")
	if err != nil {
		t.Fatalf("Lookup returned error: %v", err)
	}
	if got == nil {
		t.Fatal("Lookup returned nil handler")
	}
}

func TestTaskRegistryLookupRejectsUnknownTask(t *testing.T) {
	registry := NewTaskRegistry()

	_, err := registry.Lookup("email.send")
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("Lookup error = %v, want ErrTaskNotFound", err)
	}
}

func TestTaskRegistryNamesReturnsSortedNames(t *testing.T) {
	registry := NewTaskRegistry()
	handler := TaskHandlerFunc(func(HandlerContext, TaskPayload) (TaskResult, error) {
		return SucceededResult(nil), nil
	})

	if err := registry.Register("video.transcode", handler); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if err := registry.Register("email.send", handler); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	names := registry.Names()
	if len(names) != 2 || names[0] != "email.send" || names[1] != "video.transcode" {
		t.Fatalf("Names = %#v, want sorted task names", names)
	}
}
