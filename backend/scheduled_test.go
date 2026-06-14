package backend

import (
	"errors"
	"testing"
	"time"
)

func TestMoveDueScheduledRequestValidateAcceptsDueRequest(t *testing.T) {
	request := MoveDueScheduledRequest{
		Queue: "default",
		Now:   time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC),
		Limit: 10,
	}

	if err := request.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestMoveDueScheduledRequestValidateRejectsNegativeLimit(t *testing.T) {
	err := (MoveDueScheduledRequest{Queue: "default", Limit: -1}).Validate()
	if !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}
