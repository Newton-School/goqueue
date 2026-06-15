package redisbackend

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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

	pipe := b.client.TxPipeline()
	pipe.HDel(ctx, b.keys.periodicDefinitionsHash(), request.Name)
	pipe.ZRem(ctx, b.keys.periodicDueSet(), request.Name)
	pipe.Del(ctx, b.keys.periodicLease(request.Name))

	_, err := pipe.Exec(ctx)
	return err
}

// ListDuePeriodicTasks leases due periodic task definitions for a scheduler.
func (b *Backend) ListDuePeriodicTasks(ctx context.Context, request backend.ListDuePeriodicTasksRequest) ([]backend.DuePeriodicTask, error) {
	if b.client == nil {
		return nil, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	names, err := b.client.ZRangeByScore(ctx, b.keys.periodicDueSet(), &redis.ZRangeBy{
		Min:    "-inf",
		Max:    fmt.Sprintf("%d", unixMillis(request.Now)),
		Offset: 0,
		Count:  request.Limit,
	}).Result()
	if err != nil {
		return nil, err
	}

	due := make([]backend.DuePeriodicTask, 0, len(names))
	codec := periodicTaskCodec{}
	for _, name := range names {
		token, err := newPeriodicLockToken()
		if err != nil {
			return nil, err
		}

		lockKey := b.keys.periodicLease(name)
		acquired, err := b.client.SetNX(ctx, lockKey, token, request.LockTTL).Result()
		if err != nil {
			return nil, err
		}
		if !acquired {
			continue
		}

		encoded, err := b.client.HGet(ctx, b.keys.periodicDefinitionsHash(), name).Bytes()
		if err == redis.Nil {
			pipe := b.client.TxPipeline()
			pipe.Del(ctx, lockKey)
			pipe.ZRem(ctx, b.keys.periodicDueSet(), name)
			if _, execErr := pipe.Exec(ctx); execErr != nil {
				return nil, execErr
			}
			continue
		}
		if err != nil {
			return nil, err
		}

		record, err := codec.decode(encoded)
		if err != nil {
			return nil, err
		}

		due = append(due, backend.DuePeriodicTask{
			Record:      record,
			LockToken:   token,
			LockedUntil: request.Now.UTC().Add(request.LockTTL),
		})
	}

	return due, nil
}

// MarkPeriodicTaskDispatched advances a periodic definition after dispatch.
func (b *Backend) MarkPeriodicTaskDispatched(ctx context.Context, request backend.MarkPeriodicTaskDispatchedRequest) error {
	if b.client == nil {
		return fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return err
	}

	encoded, err := b.client.HGet(ctx, b.keys.periodicDefinitionsHash(), request.Name).Bytes()
	if err == redis.Nil {
		return backend.ErrPeriodicTaskNotFound
	}
	if err != nil {
		return err
	}

	record, err := (periodicTaskCodec{}).decode(encoded)
	if err != nil {
		return err
	}
	record.NextDueAt = request.NextDueAt.UTC()
	record.UpdatedAt = request.DispatchedAt.UTC()

	updated, err := (periodicTaskCodec{}).encode(record)
	if err != nil {
		return err
	}

	result, err := redis.NewScript(markPeriodicDispatchedScript()).Run(
		ctx,
		b.client,
		[]string{
			b.keys.periodicDefinitionsHash(),
			b.keys.periodicDueSet(),
			b.keys.periodicLease(request.Name),
		},
		request.LockToken,
		request.Name,
		string(updated),
		unixMillis(request.NextDueAt),
	).Int64()
	if err != nil {
		return err
	}
	if result == 0 {
		return backend.ErrPeriodicTaskLeaseLost
	}

	return nil
}

func newPeriodicLockToken() (string, error) {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", fmt.Errorf("%w: generate periodic lock token: %v", ErrInvalidRedisMessage, err)
	}

	return hex.EncodeToString(bytes[:]), nil
}
