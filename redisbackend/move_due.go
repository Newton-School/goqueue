package redisbackend

import (
	"context"
	"fmt"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/redis/go-redis/v9"
)

// MoveDueScheduled moves due scheduled tasks into the ready stream.
func (b *Backend) MoveDueScheduled(ctx context.Context, request backend.MoveDueScheduledRequest) ([]backend.MovedScheduledMessage, error) {
	if b.client == nil {
		return nil, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	now := request.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	limit := request.Limit
	if limit == 0 {
		limit = 100
	}

	values, err := redis.NewScript(moveDueScheduledScript()).Run(
		ctx,
		b.client,
		[]string{
			b.keys.scheduledSet(request.Queue.String()),
			b.keys.readyStream(request.Queue.String()),
		},
		unixMillis(now),
		limit,
		b.keys.taskPrefix(),
	).Slice()
	if err != nil {
		return nil, err
	}

	return parseMovedScheduledMessages(values)
}
