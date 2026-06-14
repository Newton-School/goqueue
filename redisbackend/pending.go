package redisbackend

import (
	"context"
	"fmt"

	"github.com/Newton-School/goqueue/backend"
	"github.com/redis/go-redis/v9"
)

// ClaimStaleReady claims pending ready messages that have been idle too long.
func (b *Backend) ClaimStaleReady(ctx context.Context, request backend.ClaimStaleReadyRequest) ([]backend.ReadyMessage, error) {
	if b.client == nil {
		return nil, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	count := request.Count
	if count == 0 {
		count = 100
	}
	startID := request.StartID
	if startID == "" {
		startID = "0-0"
	}

	messages, _, err := b.client.XAutoClaim(ctx, &redis.XAutoClaimArgs{
		Stream:   b.keys.readyStream(request.Queue.String()),
		Group:    request.Group,
		Consumer: request.Consumer,
		MinIdle:  request.MinIdle,
		Start:    startID,
		Count:    count,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	return parseReadyMessages(messages)
}
