package backend

import (
	"context"
	"testing"
)

func TestQueueBackendInterfaceAcceptsImplementation(t *testing.T) {
	var backend QueueBackend = noopBackend{}
	if backend == nil {
		t.Fatal("QueueBackend should accept implementations")
	}
}

type noopBackend struct{}

func (noopBackend) Ping(context.Context) error { return nil }
func (noopBackend) Close() error               { return nil }
