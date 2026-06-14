package redisbackend

import (
	"context"
	"fmt"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
	"github.com/redis/go-redis/v9"
)

// EnqueueReady atomically stores a task message and appends it to the ready stream.
func (b *Backend) EnqueueReady(ctx context.Context, request backend.EnqueueRequest) (backend.EnqueueResponse, error) {
	if b.client == nil {
		return backend.EnqueueResponse{}, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return backend.EnqueueResponse{}, err
	}

	encoded, err := (messageCodec{}).encode(request.Message)
	if err != nil {
		return backend.EnqueueResponse{}, err
	}

	streamID, err := redis.NewScript(readyEnqueueScript()).Run(
		ctx,
		b.client,
		[]string{
			b.keys.message(request.Message.ID),
			b.keys.readyStream(request.Message.Queue),
		},
		string(encoded),
		ttlSeconds(b.options.MessageTTL),
		request.Message.ID,
	).Text()
	if err != nil {
		return backend.EnqueueResponse{}, err
	}

	return backend.EnqueueResponse{
		TaskID:   task.TaskID(request.Message.ID),
		StreamID: streamID,
	}, nil
}

// EnqueueScheduled atomically stores a task message and schedules its task ID.
func (b *Backend) EnqueueScheduled(ctx context.Context, request backend.EnqueueRequest) (backend.EnqueueResponse, error) {
	if request.Message.Timing.ETA.IsZero() {
		return backend.EnqueueResponse{}, fmt.Errorf("%w: scheduled task requires eta", ErrInvalidRedisMessage)
	}
	if b.client == nil {
		return backend.EnqueueResponse{}, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return backend.EnqueueResponse{}, err
	}

	encoded, err := (messageCodec{}).encode(request.Message)
	if err != nil {
		return backend.EnqueueResponse{}, err
	}

	_, err = redis.NewScript(scheduledEnqueueScript()).Run(
		ctx,
		b.client,
		[]string{
			b.keys.message(request.Message.ID),
			b.keys.scheduledSet(request.Message.Queue),
		},
		string(encoded),
		ttlSeconds(b.options.MessageTTL),
		unixMillis(request.Message.Timing.ETA),
		request.Message.ID,
	).Text()
	if err != nil {
		return backend.EnqueueResponse{}, err
	}

	return backend.EnqueueResponse{
		TaskID:    task.TaskID(request.Message.ID),
		Scheduled: true,
	}, nil
}
