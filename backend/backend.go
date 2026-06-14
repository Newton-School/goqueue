package backend

import "context"

// QueueBackend is the storage boundary used by producers, workers, and schedulers.
type QueueBackend interface {
	EnqueueReady(context.Context, EnqueueRequest) (EnqueueResponse, error)
	EnqueueScheduled(context.Context, EnqueueRequest) (EnqueueResponse, error)
	Ping(context.Context) error
	Close() error
}
