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

// GetTaskResult returns a stored task result.
func (b *Backend) GetTaskResult(ctx context.Context, taskID task.TaskID) (backend.TaskResultRecord, error) {
	if b.client == nil {
		return backend.TaskResultRecord{}, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := task.ValidateTaskID(taskID.String()); err != nil {
		return backend.TaskResultRecord{}, err
	}

	data, err := b.client.Get(ctx, b.keys.result(taskID.String())).Bytes()
	if errors.Is(err, redis.Nil) {
		return backend.TaskResultRecord{}, backend.ErrTaskResultNotFound
	}
	if err != nil {
		return backend.TaskResultRecord{}, err
	}

	return (resultCodec{}).decode(data)
}
