package goqueue

import (
	"context"
	"time"

	"github.com/Newton-School/goqueue/task"
)

const (
	MinPriority     = task.MinPriority
	DefaultPriority = task.DefaultPriority
	MaxPriority     = task.MaxPriority
)

const (
	TaskPending      = task.TaskPending
	TaskScheduled    = task.TaskScheduled
	TaskReceived     = task.TaskReceived
	TaskStarted      = task.TaskStarted
	TaskRetrying     = task.TaskRetrying
	TaskSucceeded    = task.TaskSucceeded
	TaskFailed       = task.TaskFailed
	TaskRevoked      = task.TaskRevoked
	TaskExpired      = task.TaskExpired
	TaskDeadLettered = task.TaskDeadLettered
)

type (
	TaskName          = task.TaskName
	QueueName         = task.QueueName
	TaskID            = task.TaskID
	Priority          = task.Priority
	TaskState         = task.TaskState
	RetryPolicy       = task.RetryPolicy
	TaskTiming        = task.TaskTiming
	TaskPayload       = task.TaskPayload
	TaskMetadata      = task.TaskMetadata
	PayloadCodec      = task.PayloadCodec
	JSONPayloadCodec  = task.JSONPayloadCodec
	TaskEnvelope      = task.TaskEnvelope
	TaskEnvelopeInput = task.TaskEnvelopeInput
	TaskMessage       = task.TaskMessage
	HandlerContext    = task.HandlerContext
	TaskResult        = task.TaskResult
	TaskHandler       = task.TaskHandler
	TaskHandlerFunc   = task.TaskHandlerFunc
	TaskRegistry      = task.TaskRegistry
)

var (
	ErrInvalidTaskName    = task.ErrInvalidTaskName
	ErrInvalidTaskID      = task.ErrInvalidTaskID
	ErrInvalidPriority    = task.ErrInvalidPriority
	ErrInvalidTaskState   = task.ErrInvalidTaskState
	ErrInvalidRetryPolicy = task.ErrInvalidRetryPolicy
	ErrInvalidTaskTiming  = task.ErrInvalidTaskTiming
	ErrInvalidPayload     = task.ErrInvalidPayload
	ErrDuplicateTask      = task.ErrDuplicateTask
	ErrInvalidTaskHandler = task.ErrInvalidTaskHandler
	ErrTaskNotFound       = task.ErrTaskNotFound
)

func ValidateTaskName(name string) error {
	return task.ValidateTaskName(name)
}

func ValidateQueueName(name string) error {
	return task.ValidateQueueName(name)
}

func NewTaskID() (TaskID, error) {
	return task.NewTaskID()
}

func ValidateTaskID(id string) error {
	return task.ValidateTaskID(id)
}

func ValidatePriority(priority Priority) error {
	return task.ValidatePriority(priority)
}

func ValidateTaskState(state TaskState) error {
	return task.ValidateTaskState(state)
}

func DefaultRetryPolicy() RetryPolicy {
	return task.DefaultRetryPolicy()
}

func TaskTimingFromCountdown(now time.Time, countdown time.Duration) (TaskTiming, error) {
	return task.TaskTimingFromCountdown(now, countdown)
}

func NewTaskPayload(args []any, kwargs map[string]any) TaskPayload {
	return task.NewTaskPayload(args, kwargs)
}

func NewTaskMetadata(values map[string]string) TaskMetadata {
	return task.NewTaskMetadata(values)
}

func NewTaskEnvelope(input TaskEnvelopeInput) (TaskEnvelope, error) {
	return task.NewTaskEnvelope(input)
}

func TaskEnvelopeToMessage(envelope TaskEnvelope, codec PayloadCodec) (TaskMessage, error) {
	return task.TaskEnvelopeToMessage(envelope, codec)
}

func TaskMessageToEnvelope(message TaskMessage, codec PayloadCodec) (TaskEnvelope, error) {
	return task.TaskMessageToEnvelope(message, codec)
}

func NewHandlerContext(ctx context.Context, envelope TaskEnvelope) HandlerContext {
	return task.NewHandlerContext(ctx, envelope)
}

func SucceededResult(value any) TaskResult {
	return task.SucceededResult(value)
}

func FailedResult(err error) TaskResult {
	return task.FailedResult(err)
}

func NewTaskRegistry() *TaskRegistry {
	return task.NewTaskRegistry()
}
