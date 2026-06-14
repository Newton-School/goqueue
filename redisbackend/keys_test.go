package redisbackend

import "testing"

func TestKeyBuilderReadyStream(t *testing.T) {
	keys := newKeyBuilder("payments")

	got := keys.readyStream("emails")
	want := "payments:queue:emails:ready"
	if got != want {
		t.Fatalf("readyStream = %q, want %q", got, want)
	}
}

func TestKeyBuilderScheduledSet(t *testing.T) {
	keys := newKeyBuilder("payments")

	got := keys.scheduledSet("emails")
	want := "payments:queue:emails:scheduled"
	if got != want {
		t.Fatalf("scheduledSet = %q, want %q", got, want)
	}
}

func TestKeyBuilderDeadLetterStream(t *testing.T) {
	keys := newKeyBuilder("payments")

	got := keys.deadLetterStream("emails")
	want := "payments:queue:emails:dead"
	if got != want {
		t.Fatalf("deadLetterStream = %q, want %q", got, want)
	}
}

func TestKeyBuilderTaskKeys(t *testing.T) {
	keys := newKeyBuilder("payments")
	id := "4ac0a01f-1b16-4330-b3e7-e99826eacb1a"

	if got := keys.message(id); got != "payments:task:4ac0a01f-1b16-4330-b3e7-e99826eacb1a:message" {
		t.Fatalf("message key = %q", got)
	}
	if got := keys.state(id); got != "payments:task:4ac0a01f-1b16-4330-b3e7-e99826eacb1a:state" {
		t.Fatalf("state key = %q", got)
	}
	if got := keys.result(id); got != "payments:task:4ac0a01f-1b16-4330-b3e7-e99826eacb1a:result" {
		t.Fatalf("result key = %q", got)
	}
}

func TestKeyBuilderTaskPrefix(t *testing.T) {
	keys := newKeyBuilder("payments")

	if got := keys.taskPrefix(); got != "payments:task:" {
		t.Fatalf("taskPrefix = %q, want payments:task:", got)
	}
}
