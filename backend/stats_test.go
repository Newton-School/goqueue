package backend

import (
	"errors"
	"testing"
)

func TestQueueStatsRequestValidateAcceptsQueue(t *testing.T) {
	request := QueueStatsRequest{Queue: "default"}
	if err := request.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestQueueStatsRequestValidateRejectsBlankQueue(t *testing.T) {
	err := (QueueStatsRequest{}).Validate()
	if !errors.Is(err, ErrInvalidQueueName) {
		t.Fatalf("Validate error = %v, want ErrInvalidQueueName", err)
	}
}
