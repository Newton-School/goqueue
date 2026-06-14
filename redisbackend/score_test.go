package redisbackend

import (
	"testing"
	"time"
)

func TestUnixMillisConvertsTime(t *testing.T) {
	timestamp := time.Unix(10, 250*int64(time.Millisecond)).UTC()

	if got := unixMillis(timestamp); got != 10250 {
		t.Fatalf("unixMillis = %d, want 10250", got)
	}
}
