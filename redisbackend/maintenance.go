package redisbackend

import (
	"context"
	"errors"
	"fmt"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
	"github.com/redis/go-redis/v9"
)

// GetTaskMessage loads a persisted task message by ID.
func (b *Backend) GetTaskMessage(ctx context.Context, taskID task.TaskID) (task.TaskMessage, error) {
	if b.client == nil {
		return task.TaskMessage{}, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := task.ValidateTaskID(taskID.String()); err != nil {
		return task.TaskMessage{}, err
	}

	raw, err := b.client.Get(ctx, b.keys.message(taskID.String())).Bytes()
	if errors.Is(err, redis.Nil) {
		return task.TaskMessage{}, backend.ErrTaskMessageNotFound
	}
	if err != nil {
		return task.TaskMessage{}, err
	}

	return (messageCodec{}).decode(raw)
}

// DeleteTaskMessage removes a persisted task message by ID.
func (b *Backend) DeleteTaskMessage(ctx context.Context, taskID task.TaskID) error {
	if b.client == nil {
		return fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := task.ValidateTaskID(taskID.String()); err != nil {
		return err
	}

	return b.client.Del(ctx, b.keys.message(taskID.String())).Err()
}

// ReadDeadLetter returns a single dead-letter stream entry.
func (b *Backend) ReadDeadLetter(ctx context.Context, queue task.QueueName, streamID string) (backend.DeadLetterRecord, error) {
	if b.client == nil {
		return backend.DeadLetterRecord{}, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := task.ValidateQueueName(queue.String()); err != nil {
		return backend.DeadLetterRecord{}, err
	}
	if streamID == "" {
		return backend.DeadLetterRecord{}, fmt.Errorf("%w: dead-letter stream id is required", backend.ErrInvalidBackendRequest)
	}

	entries, err := b.client.XRangeN(ctx, b.keys.deadLetterStream(queue.String()), streamID, streamID, 1).Result()
	if err != nil {
		if err == redis.Nil {
			return backend.DeadLetterRecord{}, backend.ErrDeadLetterNotFound
		}
		return backend.DeadLetterRecord{}, err
	}
	if len(entries) == 0 {
		return backend.DeadLetterRecord{}, backend.ErrDeadLetterNotFound
	}

	decoded, err := (deadLetterCodec{}).decode(entries[0].ID, entries[0].Values)
	if err != nil {
		return backend.DeadLetterRecord{}, err
	}

	if decoded.StreamID == "" {
		decoded.StreamID = streamID
	}

	return decoded, nil
}

// DeleteDeadLetters removes dead-letter stream entries by stream ID.
func (b *Backend) DeleteDeadLetters(ctx context.Context, queue task.QueueName, streamIDs ...string) (int64, error) {
	if b.client == nil {
		return 0, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := task.ValidateQueueName(queue.String()); err != nil {
		return 0, err
	}
	if len(streamIDs) == 0 {
		return 0, nil
	}

	validIDs := make([]string, 0, len(streamIDs))
	for _, id := range streamIDs {
		if id == "" {
			return 0, fmt.Errorf("%w: dead-letter stream id is required", backend.ErrInvalidBackendRequest)
		}
		validIDs = append(validIDs, id)
	}

	return b.client.XDel(ctx, b.keys.deadLetterStream(queue.String()), validIDs...).Result()
}

// PurgeQueue deletes queue streams, pending payload metadata, and optional task state/result rows.
func (b *Backend) PurgeQueue(ctx context.Context, request backend.PurgeQueueRequest) (backend.PurgeQueueResult, error) {
	if b.client == nil {
		return backend.PurgeQueueResult{}, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return backend.PurgeQueueResult{}, err
	}

	result := backend.PurgeQueueResult{Queue: request.Queue}

	readyStream := b.keys.readyStream(request.Queue.String())
	scheduledSet := b.keys.scheduledSet(request.Queue.String())
	deadLetterStream := b.keys.deadLetterStream(request.Queue.String())

	taskIDs := make(map[task.TaskID]struct{})

	readyMessages, err := b.client.XRange(ctx, readyStream, "-", "+").Result()
	if err != nil && err != redis.Nil {
		return backend.PurgeQueueResult{}, err
	}
	readyParsed, err := parseReadyMessages(readyMessages)
	if err != nil {
		return backend.PurgeQueueResult{}, err
	}
	for _, readyMessage := range readyParsed {
		taskIDs[task.TaskID(readyMessage.Message.ID)] = struct{}{}
	}

	scheduledTaskIDs, err := b.client.ZRange(ctx, scheduledSet, 0, -1).Result()
	if err != nil && err != redis.Nil {
		return backend.PurgeQueueResult{}, err
	}
	for _, taskID := range scheduledTaskIDs {
		taskIDs[task.TaskID(taskID)] = struct{}{}
	}

	deadLetterEntries, err := b.client.XRange(ctx, deadLetterStream, "-", "+").Result()
	if err != nil && err != redis.Nil {
		return backend.PurgeQueueResult{}, err
	}
	for _, entry := range deadLetterEntries {
		decoded, err := (deadLetterCodec{}).decode(entry.ID, entry.Values)
		if err != nil {
			return backend.PurgeQueueResult{}, err
		}
		taskIDs[task.TaskID(decoded.Message.ID)] = struct{}{}
	}

	if request.DeleteMessages {
		messageKeys := make([]string, 0, len(taskIDs))
		for id := range taskIDs {
			messageKeys = append(messageKeys, b.keys.message(id.String()))
		}

		deleted, err := b.deleteMessageKeys(ctx, messageKeys)
		if err != nil {
			return backend.PurgeQueueResult{}, err
		}
		result.TaskMessages = deleted
	}

	if request.DeleteStates {
		stateKeys := make([]string, 0, len(taskIDs))
		for id := range taskIDs {
			stateKeys = append(stateKeys, b.keys.state(id.String()))
		}

		deleted, err := b.deleteMessageKeys(ctx, stateKeys)
		if err != nil {
			return backend.PurgeQueueResult{}, err
		}
		result.TaskStates = deleted
	}

	if request.DeleteResults {
		resultKeys := make([]string, 0, len(taskIDs))
		for id := range taskIDs {
			resultKeys = append(resultKeys, b.keys.result(id.String()))
		}

		deleted, err := b.deleteMessageKeys(ctx, resultKeys)
		if err != nil {
			return backend.PurgeQueueResult{}, err
		}
		result.TaskResults = deleted
	}

	readyDeleted, err := b.client.Del(ctx, readyStream).Result()
	if err != nil {
		return backend.PurgeQueueResult{}, err
	}
	if readyDeleted > 0 {
		result.ReadyStream = 1
	}

	scheduledDeleted, err := b.client.Del(ctx, scheduledSet).Result()
	if err != nil {
		return backend.PurgeQueueResult{}, err
	}
	if scheduledDeleted > 0 {
		result.ScheduledSet = 1
	}

	deadLetterDeleted, err := b.client.Del(ctx, deadLetterStream).Result()
	if err != nil {
		return backend.PurgeQueueResult{}, err
	}
	if deadLetterDeleted > 0 {
		result.DeadLetterStream = 1
	}

	return result, nil
}

func (b *Backend) deleteMessageKeys(ctx context.Context, keys []string) (int64, error) {
	if len(keys) == 0 {
		return 0, nil
	}

	deleted, err := b.client.Del(ctx, keys...).Result()
	if err != nil {
		return 0, err
	}

	return deleted, nil
}
