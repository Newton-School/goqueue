package redisbackend

import (
	"context"
	"fmt"

	"github.com/Newton-School/goqueue/backend"
)

// QueueStats returns ready and scheduled counts for a queue.
func (b *Backend) QueueStats(ctx context.Context, request backend.QueueStatsRequest) (backend.QueueStats, error) {
	if b.client == nil {
		return backend.QueueStats{}, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return backend.QueueStats{}, err
	}

	readyCount, err := b.client.XLen(ctx, b.keys.readyStream(request.Queue.String())).Result()
	if err != nil {
		return backend.QueueStats{}, err
	}

	scheduledCount, err := b.client.ZCard(ctx, b.keys.scheduledSet(request.Queue.String())).Result()
	if err != nil {
		return backend.QueueStats{}, err
	}

	deadLetterCount, err := b.client.XLen(ctx, b.keys.deadLetterStream(request.Queue.String())).Result()
	if err != nil {
		return backend.QueueStats{}, err
	}

	return backend.QueueStats{
		Queue:           request.Queue,
		ReadyCount:      readyCount,
		ScheduledCount:  scheduledCount,
		DeadLetterCount: deadLetterCount,
	}, nil
}
