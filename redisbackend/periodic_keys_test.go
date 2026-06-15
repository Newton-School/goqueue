package redisbackend

import "testing"

func TestKeyBuilderPeriodicDefinitionsHash(t *testing.T) {
	keys := newKeyBuilder("tenant")

	if got := keys.periodicDefinitionsHash(); got != "tenant:scheduler:periodic:definitions" {
		t.Fatalf("periodicDefinitionsHash = %q", got)
	}
}

func TestKeyBuilderPeriodicDueSet(t *testing.T) {
	keys := newKeyBuilder("tenant")

	if got := keys.periodicDueSet(); got != "tenant:scheduler:periodic:due" {
		t.Fatalf("periodicDueSet = %q", got)
	}
}

func TestKeyBuilderPeriodicLease(t *testing.T) {
	keys := newKeyBuilder("tenant")

	if got := keys.periodicLease("welcome-email"); got != "tenant:scheduler:periodic:welcome-email:lease" {
		t.Fatalf("periodicLease = %q", got)
	}
}
