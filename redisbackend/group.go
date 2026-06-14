package redisbackend

import (
	"context"
	"fmt"
	"strings"

	"github.com/Newton-School/goqueue/backend"
)

// EnsureConsumerGroup creates a consumer group for a queue if it is missing.
func (b *Backend) EnsureConsumerGroup(ctx context.Context, request backend.ConsumerGroupRequest) error {
	if b.client == nil {
		return fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return err
	}

	err := b.client.XGroupCreateMkStream(ctx, b.keys.readyStream(request.Queue.String()), request.Group, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return err
	}

	return nil
}
