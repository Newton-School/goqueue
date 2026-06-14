package backend

import (
	"errors"
	"testing"
)

func TestConsumerGroupRequestValidateAcceptsGroup(t *testing.T) {
	request := ConsumerGroupRequest{Queue: "default", Group: "workers"}
	if err := request.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestConsumerGroupRequestValidateRejectsMissingGroup(t *testing.T) {
	err := (ConsumerGroupRequest{Queue: "default"}).Validate()
	if !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}
