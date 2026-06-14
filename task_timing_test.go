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

func TestTaskTimingFromCountdownSetsETA(t *testing.T) {
	now := time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC)

	timing, err := TaskTimingFromCountdown(now, 30*time.Second)
	if err != nil {
		t.Fatalf("TaskTimingFromCountdown returned error: %v", err)
	}

	if !timing.ETA.Equal(now.Add(30 * time.Second)) {
		t.Fatalf("ETA = %s, want %s", timing.ETA, now.Add(30*time.Second))
	}
}

func TestTaskTimingFromCountdownRejectsNegativeCountdown(t *testing.T) {
	_, err := TaskTimingFromCountdown(time.Now(), -time.Second)
	if !errors.Is(err, ErrInvalidTaskTiming) {
		t.Fatalf("TaskTimingFromCountdown error = %v, want ErrInvalidTaskTiming", err)
	}
}
