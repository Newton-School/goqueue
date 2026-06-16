package inspect

import (
	"context"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

// QueueStats returns queue storage and dead-letter metrics.
func (i *Inspector) QueueStats(ctx context.Context, queue task.QueueName) (backend.QueueStats, error) {
	if i == nil {
		return backend.QueueStats{}, ErrNilInspector
	}
	if i.backend == nil {
		return backend.QueueStats{}, ErrInspectorBackend
	}
	if err := task.ValidateQueueName(queue.String()); err != nil {
		return backend.QueueStats{}, err
	}

	return i.backend.QueueStats(ctx, backend.QueueStatsRequest{Queue: queue})
}
