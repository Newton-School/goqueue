package redisbackend

import (
	"context"
	"fmt"

	"github.com/Newton-School/goqueue/backend"
	"github.com/redis/go-redis/v9"
)

// ReadReady reads ready messages for a Redis consumer group.
func (b *Backend) ReadReady(ctx context.Context, request backend.ReadReadyRequest) ([]backend.ReadyMessage, error) {
	if b.client == nil {
		return nil, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	count := request.Count
	if count == 0 {
		count = 1
	}

	streams, err := b.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    request.Group,
		Consumer: request.Consumer,
		Streams:  []string{b.keys.readyStream(request.Queue.String()), ">"},
		Count:    count,
		Block:    request.Block,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	return parseReadyStreamMessages(streams)
}
