package worker

import (
	"context"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

func (w *Worker) deadLetterTask(
	ctx context.Context,
	streamID string,
	envelope task.TaskEnvelope,
	message backend.ReadyMessage,
	reason task.FailureCategory,
	result task.TaskResult,
) error {
	if !w.deadLetterEnabled {
		return nil
	}

	_, err := w.backend.EnqueueDeadLetter(ctx, backend.DeadLetterRequest{
		Message:        message.Message,
		Reason:         reason,
		Error:          result.Error,
		SourceStreamID: streamID,
		Group:          w.group,
		Consumer:       w.consumer,
		FailedAt:       w.now(),
	})
	if err != nil {
		return err
	}

	result.Metadata = task.MergeFailureMetadata(result.Metadata, task.FailureMetadata{
		Category:       reason,
		Attempt:        envelope.Attempt,
		MaxAttempts:    envelope.RetryPolicy.MaxAttempts,
		Retryable:      false,
		DeadLettered:   true,
		DeadLetteredAt: w.now(),
		LastError:      result.Error,
	})
	if err := w.writeState(ctx, envelope.ID, task.TaskDeadLettered, result.Error); err != nil {
		return err
	}
	return w.saveResult(ctx, envelope.ID, result)
}
