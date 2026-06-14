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

func (noopBackend) EnqueueReady(context.Context, EnqueueRequest) (EnqueueResponse, error) {
	return EnqueueResponse{}, nil
}
func (noopBackend) EnqueueScheduled(context.Context, EnqueueRequest) (EnqueueResponse, error) {
	return EnqueueResponse{}, nil
}
func (noopBackend) Ping(context.Context) error { return nil }
func (noopBackend) Close() error               { return nil }
