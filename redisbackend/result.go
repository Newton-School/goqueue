package redisbackend

import (
	"context"
	"fmt"
	"time"

	"github.com/Newton-School/goqueue/backend"
)

// SaveTaskResult stores a task result.
func (b *Backend) SaveTaskResult(ctx context.Context, record backend.TaskResultRecord) error {
	if b.client == nil {
		return fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := record.Validate(); err != nil {
		return err
	}

	if record.UpdatedAt.IsZero() {
		record.UpdatedAt = time.Now().UTC()
	}

	encoded, err := (resultCodec{}).encode(record)
	if err != nil {
		return err
	}

	ttl := record.TTL
	if ttl == 0 {
		ttl = b.options.ResultTTL
	}

	return b.client.Set(ctx, b.keys.result(record.TaskID.String()), encoded, ttl).Err()
}
