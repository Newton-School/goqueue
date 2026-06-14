package redisbackend

import (
	"context"
	"testing"
)

func TestBackendPingRejectsNilClient(t *testing.T) {
	backend := &Backend{}

	err := backend.Ping(context.Background())
	if err == nil {
		t.Fatal("Ping should fail without a Redis client")
	}
}

func TestBackendCloseAllowsNilClient(t *testing.T) {
	backend := &Backend{}

	if err := backend.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
}
