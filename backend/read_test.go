package backend

import (
	"errors"
	"testing"
	"time"
)

func TestReadReadyRequestValidateAcceptsConsumerRead(t *testing.T) {
	request := ReadReadyRequest{
		Queue:    "default",
		Group:    "workers",
		Consumer: "pod-1",
		Count:    10,
		Block:    time.Second,
	}

	if err := request.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestReadReadyRequestValidateRejectsMissingGroup(t *testing.T) {
	err := (ReadReadyRequest{Queue: "default", Consumer: "pod-1"}).Validate()
	if !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func TestAckRequestValidateRequiresStreamID(t *testing.T) {
	err := (AckRequest{Queue: "default", Group: "workers"}).Validate()
	if !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}
