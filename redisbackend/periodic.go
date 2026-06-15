package redisbackend

import (
	"context"
	"fmt"

	"github.com/Newton-School/goqueue/backend"
)

// UpsertPeriodicTask stores or replaces a periodic task definition.
func (b *Backend) UpsertPeriodicTask(ctx context.Context, request backend.UpsertPeriodicTaskRequest) error {
	if b.client == nil {
		return fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return err
	}

	return fmt.Errorf("%w: periodic upsert not implemented", ErrInvalidRedisMessage)
}

// DeletePeriodicTask removes a periodic task definition.
func (b *Backend) DeletePeriodicTask(ctx context.Context, request backend.DeletePeriodicTaskRequest) error {
	if b.client == nil {
		return fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return err
	}

	return fmt.Errorf("%w: periodic delete not implemented", ErrInvalidRedisMessage)
}

// ListDuePeriodicTasks leases due periodic task definitions for a scheduler.
func (b *Backend) ListDuePeriodicTasks(ctx context.Context, request backend.ListDuePeriodicTasksRequest) ([]backend.DuePeriodicTask, error) {
	if b.client == nil {
		return nil, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("%w: periodic due scan not implemented", ErrInvalidRedisMessage)
}

// MarkPeriodicTaskDispatched advances a periodic definition after dispatch.
func (b *Backend) MarkPeriodicTaskDispatched(ctx context.Context, request backend.MarkPeriodicTaskDispatchedRequest) error {
	if b.client == nil {
		return fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return err
	}

	return fmt.Errorf("%w: periodic dispatch mark not implemented", ErrInvalidRedisMessage)
}
