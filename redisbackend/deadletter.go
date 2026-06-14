package redisbackend

import (
	"context"
	"fmt"

	"github.com/Newton-School/goqueue/backend"
)

// EnqueueDeadLetter stores an unrecoverable task message for inspection.
func (b *Backend) EnqueueDeadLetter(ctx context.Context, request backend.DeadLetterRequest) (backend.DeadLetterRecord, error) {
	if b.client == nil {
		return backend.DeadLetterRecord{}, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return backend.DeadLetterRecord{}, err
	}

	return backend.DeadLetterRecord{}, fmt.Errorf("%w: dead letter enqueue is not implemented", ErrInvalidRedisMessage)
}

// ReadDeadLetters reads recent dead-lettered task messages.
func (b *Backend) ReadDeadLetters(ctx context.Context, request backend.ReadDeadLettersRequest) ([]backend.DeadLetterRecord, error) {
	if b.client == nil {
		return nil, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("%w: dead letter read is not implemented", ErrInvalidRedisMessage)
}
