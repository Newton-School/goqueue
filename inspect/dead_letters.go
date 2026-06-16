package inspect

import (
	"context"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

// ReadDeadLetters reads dead-lettered messages for queue inspection.
func (i *Inspector) ReadDeadLetters(ctx context.Context, queue task.QueueName, count int64) ([]backend.DeadLetterRecord, error) {
	if i == nil {
		return nil, ErrNilInspector
	}
	if i.backend == nil {
		return nil, ErrInspectorBackend
	}
	if err := queue.Validate(); err != nil {
		return nil, err
	}
	if count < 0 {
		return nil, ErrInvalidDeadLetters
	}

	return i.backend.ReadDeadLetters(ctx, backend.ReadDeadLettersRequest{
		Queue: queue,
		Count: count,
	})
}
