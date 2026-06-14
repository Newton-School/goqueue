package redisbackend

import (
	"context"
	"fmt"
)

// Ping verifies Redis connectivity.
func (b *Backend) Ping(ctx context.Context) error {
	if b.client == nil {
		return fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}

	return b.client.Ping(ctx).Err()
}

// Close closes the Redis client if one is configured.
func (b *Backend) Close() error {
	if b.client == nil {
		return nil
	}

	return b.client.Close()
}
