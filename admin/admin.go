package admin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

type controlBackend interface {
	backend.QueueBackend

	GetTaskMessage(context.Context, task.TaskID) (task.TaskMessage, error)
	ReadDeadLetter(context.Context, task.QueueName, string) (backend.DeadLetterRecord, error)
	DeleteDeadLetters(context.Context, task.QueueName, ...string) (int64, error)
	PurgeQueue(context.Context, backend.PurgeQueueRequest) (backend.PurgeQueueResult, error)
}

// Admin provides operational queue/task control APIs.
type Admin struct {
	backend controlBackend
}

// NewAdmin creates an admin client bound to a control-capable backend.
func NewAdmin(queueBackend backend.QueueBackend) (*Admin, error) {
	if queueBackend == nil {
		return nil, ErrNilAdmin
	}

	concrete, ok := queueBackend.(controlBackend)
	if !ok {
		return nil, ErrAdminBackend
	}

	return &Admin{backend: concrete}, nil
}

// RetryTask reloads a persisted message and queues it for re-execution.
func (a *Admin) RetryTask(ctx context.Context, taskID task.TaskID, options RetryTaskOptions) (RetryTaskResult, error) {
	if a == nil {
		return RetryTaskResult{}, ErrNilAdmin
	}
	if err := task.ValidateTaskID(taskID.String()); err != nil {
		return RetryTaskResult{}, ErrNilTaskID
	}
	if err := validateRetryTaskOptions(options); err != nil {
		return RetryTaskResult{}, err
	}

	stored, err := a.backend.GetTaskMessage(ctx, taskID)
	if err != nil {
		return RetryTaskResult{}, err
	}

	parsed, err := task.TaskMessageToEnvelope(stored, task.JSONPayloadCodec{})
	if err != nil {
		return RetryTaskResult{}, err
	}

	originalQueue := parsed.Queue
	retryQueue := parsed.Queue
	if options.Queue != "" {
		if err := task.ValidateQueueName(string(options.Queue)); err != nil {
			return RetryTaskResult{}, err
		}
		retryQueue = options.Queue
	}

	parsed.Queue = retryQueue
	parsed.Timing.ETA = retryTaskETA(parsed.Timing, options)
	if !options.PreserveAttempt {
		parsed.Attempt = 0
	}

	retryEnvelope := task.TaskEnvelopeInput{
		ID:          parsed.ID,
		Name:        parsed.Name,
		Queue:       retryQueue,
		Args:        parsed.Payload.Args(),
		Kwargs:      parsed.Payload.Kwargs(),
		Metadata:    parsed.Metadata.Values(),
		Timing:      parsed.Timing,
		Priority:    parsed.Priority,
		RetryPolicy: parsed.RetryPolicy,
		Attempt:     parsed.Attempt,
		CreatedAt:   parsed.CreatedAt,
	}

	envelopeMessage, err := task.NewTaskEnvelope(retryEnvelope)
	if err != nil {
		return RetryTaskResult{}, err
	}

	retryMessage, err := task.TaskEnvelopeToMessage(envelopeMessage, task.JSONPayloadCodec{})
	if err != nil {
		return RetryTaskResult{}, err
	}

	var response backend.EnqueueResponse
	if retryEnvelope.Timing.ETA.IsZero() {
		response, err = a.backend.EnqueueReady(ctx, backend.EnqueueRequest{Message: retryMessage})
	} else {
		response, err = a.backend.EnqueueScheduled(ctx, backend.EnqueueRequest{Message: retryMessage})
	}
	if err != nil {
		return RetryTaskResult{}, err
	}

	if options.ClearState {
		if err := a.backend.SetTaskState(ctx, backend.TaskStateRecord{
			TaskID: taskID,
			State:  task.TaskPending,
		}); err != nil {
			return RetryTaskResult{}, err
		}
	}

	if options.ClearResult {
		if err := a.backend.ForgetTaskResult(ctx, taskID); err != nil && !errors.Is(err, backend.ErrTaskResultNotFound) {
			return RetryTaskResult{}, err
		}
	}

	return RetryTaskResult{
		TaskID:        taskID,
		Queue:         retryQueue,
		OriginalQueue: originalQueue,
		ScheduledAt:   retryEnvelope.Timing.ETA,
		EnqueueResult: response,
		Attempt:       retryEnvelope.Attempt,
	}, nil
}

// RevokeTask sets task state to revoked for intervention.
func (a *Admin) RevokeTask(ctx context.Context, taskID task.TaskID, reason string) (RevokeTaskResult, error) {
	if a == nil {
		return RevokeTaskResult{}, ErrNilAdmin
	}
	if err := task.ValidateTaskID(taskID.String()); err != nil {
		return RevokeTaskResult{}, ErrNilTaskID
	}

	if err := a.backend.SetTaskState(ctx, backend.TaskStateRecord{
		TaskID: taskID,
		State:  task.TaskRevoked,
		Error:  reason,
	}); err != nil {
		return RevokeTaskResult{}, err
	}

	return RevokeTaskResult{TaskID: taskID, State: task.TaskRevoked}, nil
}

// ReplayDeadLetter requeues a dead-letter stream entry.
func (a *Admin) ReplayDeadLetter(ctx context.Context, queue task.QueueName, streamID string, options ReplayDeadLetterOptions) (ReplayDeadLetterResult, error) {
	if a == nil {
		return ReplayDeadLetterResult{}, ErrNilAdmin
	}
	if err := task.ValidateQueueName(string(queue)); err != nil {
		return ReplayDeadLetterResult{}, ErrInvalidQueue
	}
	if streamID == "" {
		return ReplayDeadLetterResult{}, fmt.Errorf("%w: dead-letter stream id is required", ErrInvalidControlOption)
	}

	record, err := a.backend.ReadDeadLetter(ctx, queue, streamID)
	if err != nil {
		return ReplayDeadLetterResult{}, err
	}

	retryQueue := task.QueueName(record.Message.Queue)
	if options.DestinationQueue != "" {
		if err := task.ValidateQueueName(string(options.DestinationQueue)); err != nil {
			return ReplayDeadLetterResult{}, err
		}
		retryQueue = options.DestinationQueue
	}

	message := record.Message
	message.Queue = string(retryQueue)
	message.Attempt = 0

	response, err := a.backend.EnqueueReady(ctx, backend.EnqueueRequest{Message: message})
	if err != nil {
		return ReplayDeadLetterResult{}, err
	}

	deleted := false
	if options.DeleteSource {
		if deletedCount, err := a.backend.DeleteDeadLetters(ctx, queue, streamID); err != nil {
			return ReplayDeadLetterResult{}, err
		} else if deletedCount == 1 {
			deleted = true
		}
	}

	return ReplayDeadLetterResult{
		StreamID:       streamID,
		Queue:          queue,
		Destination:    retryQueue,
		OriginalTaskID: task.TaskID(message.ID),
		EnqueueResult:  response,
		SourceDeleted:  deleted,
	}, nil
}

// DeleteDeadLetters removes dead-letter stream IDs.
func (a *Admin) DeleteDeadLetters(ctx context.Context, queue task.QueueName, streamIDs ...string) (DeleteDeadLettersResult, error) {
	if a == nil {
		return DeleteDeadLettersResult{}, ErrNilAdmin
	}
	if err := task.ValidateQueueName(string(queue)); err != nil {
		return DeleteDeadLettersResult{}, ErrInvalidQueue
	}
	if len(streamIDs) == 0 {
		return DeleteDeadLettersResult{}, fmt.Errorf("%w: at least one stream id is required", ErrInvalidControlOption)
	}

	for _, id := range streamIDs {
		if id == "" {
			return DeleteDeadLettersResult{}, fmt.Errorf("%w: stream id is required", ErrInvalidControlOption)
		}
	}

	deleted, err := a.backend.DeleteDeadLetters(ctx, queue, streamIDs...)
	if err != nil {
		return DeleteDeadLettersResult{}, err
	}

	return DeleteDeadLettersResult{Queue: queue, Deleted: deleted}, nil
}

// PurgeQueue clears storage by queue with optional state/result/message retention options.
func (a *Admin) PurgeQueue(ctx context.Context, options PurgeQueueOptions) (PurgeQueueResult, error) {
	if a == nil {
		return PurgeQueueResult{}, ErrNilAdmin
	}
	if err := task.ValidateQueueName(string(options.Queue)); err != nil {
		return PurgeQueueResult{}, ErrInvalidQueue
	}

	result, err := a.backend.PurgeQueue(ctx, backend.PurgeQueueRequest{
		Queue:          options.Queue,
		DeleteMessages: options.DeleteMessage,
		DeleteStates:   options.DeleteState,
		DeleteResults:  options.DeleteResult,
	})
	if err != nil {
		return PurgeQueueResult{}, err
	}

	return PurgeQueueResult{
		Queue:            result.Queue,
		ReadyStream:      result.ReadyStream,
		ScheduledSet:     result.ScheduledSet,
		DeadLetterStream: result.DeadLetterStream,
		TaskMessages:     result.TaskMessages,
		TaskStates:       result.TaskStates,
		TaskResults:      result.TaskResults,
	}, nil
}

func validateRetryTaskOptions(options RetryTaskOptions) error {
	if options.CountDown < 0 {
		return ErrInvalidControlOption
	}
	if !options.ScheduledAt.IsZero() && options.CountDown != 0 {
		return fmt.Errorf("%w: cannot set both scheduled_at and countdown", ErrInvalidControlOption)
	}

	return nil
}

func retryTaskETA(timing task.TaskTiming, options RetryTaskOptions) time.Time {
	if !options.ScheduledAt.IsZero() {
		return options.ScheduledAt
	}
	if options.CountDown > 0 {
		now := time.Now().UTC()
		if options.Now != nil {
			now = options.Now()
		}

		return now.Add(options.CountDown)
	}

	return timing.ETA
}
