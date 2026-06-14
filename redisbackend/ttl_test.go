package redisbackend

import (
	"testing"
	"time"
)

func TestTTLSecondsRoundsUpPositiveDuration(t *testing.T) {
	if got := ttlSeconds(1500 * time.Millisecond); got != 2 {
		t.Fatalf("ttlSeconds = %d, want 2", got)
	}
}

func TestTTLSecondsUsesOneSecondMinimum(t *testing.T) {
	if got := ttlSeconds(0); got != 1 {
		t.Fatalf("ttlSeconds = %d, want 1", got)
	}
}
