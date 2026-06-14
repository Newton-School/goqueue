package redisbackend

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
	"github.com/redis/go-redis/v9"
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

// GetTaskState returns the latest stored task state.
func (b *Backend) GetTaskState(ctx context.Context, taskID task.TaskID) (backend.TaskStateRecord, error) {
	if b.client == nil {
		return backend.TaskStateRecord{}, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := task.ValidateTaskID(taskID.String()); err != nil {
		return backend.TaskStateRecord{}, err
	}

	data, err := b.client.Get(ctx, b.keys.state(taskID.String())).Bytes()
	if errors.Is(err, redis.Nil) {
		return backend.TaskStateRecord{}, backend.ErrTaskStateNotFound
	}
	if err != nil {
		return backend.TaskStateRecord{}, err
	}

	return (stateCodec{}).decode(data)
}
