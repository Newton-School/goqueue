package redisbackend

import (
	"context"
	"fmt"

	"github.com/Newton-School/goqueue/backend"
	"github.com/redis/go-redis/v9"
)

// UpsertPeriodicTask stores or replaces a periodic task definition.
func (b *Backend) UpsertPeriodicTask(ctx context.Context, request backend.UpsertPeriodicTaskRequest) error {
	if b.client == nil {
		return fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return err
	}

	encoded, err := (periodicTaskCodec{}).encode(request.Record)
	if err != nil {
		return err
	}

	pipe := b.client.TxPipeline()
	pipe.HSet(ctx, b.keys.periodicDefinitionsHash(), request.Record.Name, string(encoded))
	pipe.ZAdd(ctx, b.keys.periodicDueSet(), redis.Z{
		Score:  float64(unixMillis(request.Record.NextDueAt)),
		Member: request.Record.Name,
	})

	_, err = pipe.Exec(ctx)
	return err
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
