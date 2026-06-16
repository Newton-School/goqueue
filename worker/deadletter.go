package worker

import (
	"context"
	"fmt"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

const (
	malformedTaskName            = "goqueue.malformed"
	malformedOriginalIDKey       = "goqueue.malformed.original_id"
	malformedOriginalNameKey     = "goqueue.malformed.original_name"
	malformedOriginalQueueKey    = "goqueue.malformed.original_queue"
	malformedGeneratedTaskIDFlag = "goqueue.malformed.generated_task_id"
)

func (w *Worker) deadLetterMalformedMessage(ctx context.Context, message backend.ReadyMessage, cause error) error {
	malformedMessage, stateTaskID, canWriteState, err := w.malformedDeadLetterMessage(message.Message)
	if err != nil {
		return err
	}

	result := task.FailedResult(cause)
	deadLetteredAt := w.now()
	result.Metadata = task.MergeFailureMetadata(result.Metadata, task.FailureMetadata{
		Category:       task.FailureMalformedMessage,
		Attempt:        malformedMessage.Attempt,
		MaxAttempts:    malformedMessage.RetryPolicy.MaxAttempts,
		Retryable:      false,
		DeadLettered:   w.deadLetterEnabled,
		DeadLetteredAt: deadLetteredAt,
		LastError:      result.Error,
	})

	if w.deadLetterEnabled {
		if _, err := w.backend.EnqueueDeadLetter(ctx, backend.DeadLetterRequest{
			Message:        malformedMessage,
			Reason:         task.FailureMalformedMessage,
			Error:          result.Error,
			SourceStreamID: message.StreamID,
			Group:          w.group,
			Consumer:       w.consumer,
			FailedAt:       deadLetteredAt,
		}); err != nil {
			return err
		}
	}

	if canWriteState {
		state := task.TaskFailed
		if w.deadLetterEnabled {
			state = task.TaskDeadLettered
		}
		if err := w.writeState(ctx, stateTaskID, state, result.Error); err != nil {
			return err
		}
		if err := w.saveResult(ctx, stateTaskID, result); err != nil {
			return err
		}
	}

	return w.ack(ctx, message.StreamID)
}

func (w *Worker) malformedDeadLetterMessage(message task.TaskMessage) (task.TaskMessage, task.TaskID, bool, error) {
	metadata := copyStringMap(message.Metadata)

	stateTaskID := task.TaskID(message.ID)
	canWriteState := task.ValidateTaskID(message.ID) == nil
	if message.ID == "" {
		generated, err := task.NewTaskID()
		if err != nil {
			return task.TaskMessage{}, "", false, fmt.Errorf("goqueue worker: generate malformed task id: %w", err)
		}
		message.ID = generated.String()
		metadata[malformedGeneratedTaskIDFlag] = "true"
	}
	if !canWriteState && stateTaskID != "" {
		metadata[malformedOriginalIDKey] = stateTaskID.String()
	}

	if message.Name == "" {
		message.Name = malformedTaskName
		metadata[malformedOriginalNameKey] = ""
	}

	if err := task.ValidateQueueName(message.Queue); err != nil {
		metadata[malformedOriginalQueueKey] = message.Queue
		message.Queue = w.queue.String()
	}

	message.Metadata = metadata
	return message, stateTaskID, canWriteState, nil
}

func copyStringMap(values map[string]string) map[string]string {
	cloned := make(map[string]string, len(values)+4)
	for key, value := range values {
		cloned[key] = value
	}
	return cloned
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
