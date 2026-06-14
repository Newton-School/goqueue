package backend

import (
	"testing"
	"time"
)

func TestClaimStaleReadyRequestValidateAcceptsValidRequest(t *testing.T) {
	request := ClaimStaleReadyRequest{
		Queue:    "default",
		Group:    "workers",
		Consumer: "pod-2",
		MinIdle:  5 * time.Minute,
		Count:    10,
		StartID:  "0-0",
	}

	if err := request.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestClaimStaleReadyRequestRejectsMissingGroup(t *testing.T) {
	request := ClaimStaleReadyRequest{
		Queue:    "default",
		Consumer: "pod-2",
		MinIdle:  time.Minute,
	}

	if err := request.Validate(); err == nil {
		t.Fatal("Validate expected error for missing group")
	}
}

func TestClaimStaleReadyRequestRejectsMissingConsumer(t *testing.T) {
	request := ClaimStaleReadyRequest{
		Queue:   "default",
		Group:   "workers",
		MinIdle: time.Minute,
	}

	if err := request.Validate(); err == nil {
		t.Fatal("Validate expected error for missing consumer")
	}
}

func TestClaimStaleReadyRequestRejectsNegativeMinIdle(t *testing.T) {
	request := ClaimStaleReadyRequest{
		Queue:    "default",
		Group:    "workers",
		Consumer: "pod-2",
		MinIdle:  -time.Second,
	}

	if err := request.Validate(); err == nil {
		t.Fatal("Validate expected error for negative min idle")
	}
}

func TestClaimStaleReadyRequestRejectsNegativeCount(t *testing.T) {
	request := ClaimStaleReadyRequest{
		Queue:    "default",
		Group:    "workers",
		Consumer: "pod-2",
		MinIdle:  time.Minute,
		Count:    -1,
	}

	if err := request.Validate(); err == nil {
		t.Fatal("Validate expected error for negative count")
	}
}
