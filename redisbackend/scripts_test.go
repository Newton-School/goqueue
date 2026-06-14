package redisbackend

import (
	"strings"
	"testing"
)

func TestReadyEnqueueScriptStoresMessageAndAddsStreamEntry(t *testing.T) {
	script := readyEnqueueScript()

	for _, fragment := range []string{"redis.call('SET'", "redis.call('XADD'", "ARGV[1]", "ARGV[2]"} {
		if !strings.Contains(script, fragment) {
			t.Fatalf("ready script missing %q in %s", fragment, script)
		}
	}
}

func TestScheduledEnqueueScriptStoresMessageAndAddsSortedSetMember(t *testing.T) {
	script := scheduledEnqueueScript()

	for _, fragment := range []string{"redis.call('SET'", "redis.call('ZADD'", "ARGV[3]", "ARGV[4]"} {
		if !strings.Contains(script, fragment) {
			t.Fatalf("scheduled script missing %q in %s", fragment, script)
		}
	}
}
