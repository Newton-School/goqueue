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

func TestMoveDueScheduledScriptMovesDueTasks(t *testing.T) {
	script := moveDueScheduledScript()

	for _, fragment := range []string{"ZRANGEBYSCORE", "redis.call('GET'", "redis.call('XADD'", "redis.call('ZREM'"} {
		if !strings.Contains(script, fragment) {
			t.Fatalf("move due script missing %q in %s", fragment, script)
		}
	}
}

func TestMarkPeriodicDispatchedScriptVerifiesLeaseToken(t *testing.T) {
	script := markPeriodicDispatchedScript()

	for _, fragment := range []string{"redis.call('GET'", "token ~= ARGV[1]", "redis.call('HSET'", "redis.call('ZADD'", "redis.call('DEL'"} {
		if !strings.Contains(script, fragment) {
			t.Fatalf("periodic mark script missing %q in %s", fragment, script)
		}
	}
}
