package goqueue

import (
	"errors"
	"testing"
	"time"
)

func TestTaskTimingValidateAcceptsExpirationAfterETA(t *testing.T) {
	eta := time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC)
	timing := TaskTiming{
		ETA:       eta,
		ExpiresAt: eta.Add(time.Minute),
	}

	if err := timing.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestTaskTimingValidateRejectsExpirationBeforeETA(t *testing.T) {
	eta := time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC)
	timing := TaskTiming{
		ETA:       eta,
		ExpiresAt: eta.Add(-time.Second),
	}

	err := timing.Validate()
	if !errors.Is(err, ErrInvalidTaskTiming) {
		t.Fatalf("Validate error = %v, want ErrInvalidTaskTiming", err)
	}
}
