package task

import (
	"errors"
	"testing"
)

func TestValidateQueueNameAcceptsRedisSafeName(t *testing.T) {
	err := ValidateQueueName("critical.emails:v1")
	if err != nil {
		t.Fatalf("ValidateQueueName returned error: %v", err)
	}
}

func TestValidateQueueNameRejectsWhitespace(t *testing.T) {
	err := ValidateQueueName("critical emails")
	if !errors.Is(err, ErrInvalidQueueName) {
		t.Fatalf("ValidateQueueName error = %v, want ErrInvalidQueueName", err)
	}
}
