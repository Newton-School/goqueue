package redisbackend

import (
	"context"
	"fmt"

	"github.com/Newton-School/goqueue/backend"
)

// Ack acknowledges a ready stream message for a consumer group.
func (b *Backend) Ack(ctx context.Context, request backend.AckRequest) error {
	if b.client == nil {
		return fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return err
	}

	return b.client.XAck(ctx, b.keys.readyStream(request.Queue.String()), request.Group, request.StreamID).Err()
}
