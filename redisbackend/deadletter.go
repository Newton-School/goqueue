package redisbackend

import (
	"context"
	"fmt"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/redis/go-redis/v9"
)

// EnqueueDeadLetter stores an unrecoverable task message for inspection.
func (b *Backend) EnqueueDeadLetter(ctx context.Context, request backend.DeadLetterRequest) (backend.DeadLetterRecord, error) {
	if b.client == nil {
		return backend.DeadLetterRecord{}, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return backend.DeadLetterRecord{}, err
	}

	failedAt := request.FailedAt
	if failedAt.IsZero() {
		failedAt = time.Now().UTC()
	}
	record := backend.DeadLetterRecord{
		Message:        request.Message,
		Reason:         request.Reason,
		Error:          request.Error,
		SourceStreamID: request.SourceStreamID,
		Group:          request.Group,
		Consumer:       request.Consumer,
		FailedAt:       failedAt,
	}

	values, err := (deadLetterCodec{}).encode(record)
	if err != nil {
		return backend.DeadLetterRecord{}, err
	}

	streamID, err := b.client.XAdd(ctx, &redis.XAddArgs{
		Stream: b.keys.deadLetterStream(request.Message.Queue),
		Values: values,
	}).Result()
	if err != nil {
		return backend.DeadLetterRecord{}, err
	}

	record.StreamID = streamID
	return record, nil
}

// ReadDeadLetters reads recent dead-lettered task messages.
func (b *Backend) ReadDeadLetters(ctx context.Context, request backend.ReadDeadLettersRequest) ([]backend.DeadLetterRecord, error) {
	if b.client == nil {
		return nil, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	count := request.Count
	if count == 0 {
		count = 100
	}

	messages, err := b.client.XRevRangeN(ctx, b.keys.deadLetterStream(request.Queue.String()), "+", "-", count).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	records := make([]backend.DeadLetterRecord, 0, len(messages))
	codec := deadLetterCodec{}
	for _, message := range messages {
		record, err := codec.decode(message.ID, message.Values)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}
