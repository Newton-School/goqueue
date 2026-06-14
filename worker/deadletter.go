package worker

import (
	"context"
	"fmt"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

func (w *Worker) deadLetterMalformedMessage(ctx context.Context, message backend.ReadyMessage, cause error) error {
	envelope, err := task.NewTaskEnvelope(task.TaskEnvelopeInput{
		ID:          task.TaskID(message.Message.ID),
		Name:        task.TaskName(message.Message.Name),
		Queue:       task.QueueName(message.Message.Queue),
		Metadata:    message.Message.Metadata,
		Timing:      message.Message.Timing,
		Priority:    message.Message.Priority,
		RetryPolicy: message.Message.RetryPolicy,
		CreatedAt:   message.Message.CreatedAt,
		Attempt:     message.Message.Attempt,
	})
	if err != nil {
		return fmt.Errorf("goqueue worker: build malformed message envelope: %w", err)
	}

	result := task.FailedResult(cause)
	if err := w.deadLetterTask(ctx, message.StreamID, envelope, message, task.FailureMalformedMessage, result); err != nil {
		return err
	}
	return w.ack(ctx, message.StreamID)
}

func (w *Worker) deadLetterTask(
	ctx context.Context,
	streamID string,
	envelope task.TaskEnvelope,
	message backend.ReadyMessage,
	reason task.FailureCategory,
	result task.TaskResult,
) error {
	if !w.deadLetterEnabled {
		if err := w.writeState(ctx, envelope.ID, task.TaskFailed, result.Error); err != nil {
			return err
		}
		return w.saveResult(ctx, envelope.ID, result)
	}

	result, err := w.recordDeadLetter(ctx, streamID, envelope, message, reason, result)
	if err != nil {
		return err
	}
	if err := w.writeState(ctx, envelope.ID, task.TaskDeadLettered, result.Error); err != nil {
		return err
	}
	return w.saveResult(ctx, envelope.ID, result)
}

func (w *Worker) recordDeadLetter(
	ctx context.Context,
	streamID string,
	envelope task.TaskEnvelope,
	message backend.ReadyMessage,
	reason task.FailureCategory,
	result task.TaskResult,
) (task.TaskResult, error) {
	if !w.deadLetterEnabled {
		return result, nil
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
		return task.TaskResult{}, err
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
	return result, nil
}
