package scheduler

import (
	"errors"
	"testing"
	"time"
)

func TestIntervalScheduleRejectsNonPositiveDurations(t *testing.T) {
	for _, interval := range []time.Duration{0, -time.Second} {
		schedule := Every(interval)

		if err := schedule.Validate(); !errors.Is(err, ErrInvalidSchedule) {
			t.Fatalf("Validate error = %v, want ErrInvalidSchedule", err)
		}
	}
}

func TestIntervalScheduleNextAddsIntervalInUTC(t *testing.T) {
	location := time.FixedZone("IST", 5*60*60+30*60)
	now := time.Date(2026, time.June, 15, 9, 30, 0, 0, location)
	schedule := Every(15 * time.Minute)

	next := schedule.Next(now)

	want := time.Date(2026, time.June, 15, 4, 15, 0, 0, time.UTC)
	if !next.Equal(want) {
		t.Fatalf("Next = %v, want %v", next, want)
	}
	if next.Location() != time.UTC {
		t.Fatalf("Next location = %v, want UTC", next.Location())
	}
}
