package redisbackend

import (
	"context"
	"fmt"
	"time"

	"github.com/Newton-School/goqueue/backend"
)

// SetTaskState stores the latest task state.
func (b *Backend) SetTaskState(ctx context.Context, record backend.TaskStateRecord) error {
	if b.client == nil {
		return fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := record.Validate(); err != nil {
		return err
	}

	if record.UpdatedAt.IsZero() {
		record.UpdatedAt = time.Now().UTC()
	}

	encoded, err := (stateCodec{}).encode(record)
	if err != nil {
		return err
	}

	ttl := record.TTL
	if ttl == 0 {
		ttl = b.options.StateTTL
	}

	return b.client.Set(ctx, b.keys.state(record.TaskID.String()), encoded, ttl).Err()
}
