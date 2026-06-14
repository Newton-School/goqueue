package goqueue

import (
	"errors"
	"testing"
)

func TestValidateTaskNameAcceptsNamespacedIdentifier(t *testing.T) {
	err := ValidateTaskName("email.send_welcome:v1")
	if err != nil {
		t.Fatalf("ValidateTaskName returned error: %v", err)
	}
}

func TestValidateTaskNameRejectsBlankName(t *testing.T) {
	err := ValidateTaskName(" ")
	if !errors.Is(err, ErrInvalidTaskName) {
		t.Fatalf("ValidateTaskName error = %v, want ErrInvalidTaskName", err)
	}
}

func TestValidateTaskNameRejectsUnsafeCharacters(t *testing.T) {
	err := ValidateTaskName("email/send")
	if !errors.Is(err, ErrInvalidTaskName) {
		t.Fatalf("ValidateTaskName error = %v, want ErrInvalidTaskName", err)
	}
}
