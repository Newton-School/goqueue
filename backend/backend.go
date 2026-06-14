package backend

import "context"

// QueueBackend is the storage boundary used by producers, workers, and schedulers.
type QueueBackend interface {
	Ping(context.Context) error
	Close() error
}
