package goqueue

import (
	"context"
	"time"

	"github.com/Newton-School/goqueue/producer"
	"github.com/Newton-School/goqueue/scheduler"
	"github.com/Newton-School/goqueue/task"
	"github.com/Newton-School/goqueue/worker"
	"github.com/Newton-School/goqueue/workflow"
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
	FailureCategory   = task.FailureCategory
	FailureMetadata   = task.FailureMetadata
	Producer          = producer.Producer
	ApplyOption       = producer.ApplyOption
	ProducerOption    = producer.ProducerOption
	AsyncResult       = producer.AsyncResult
	Scheduler         = scheduler.Scheduler
	SchedulerOption   = scheduler.SchedulerOption
	PeriodicTask      = scheduler.PeriodicTask
	PeriodicTaskName  = scheduler.PeriodicTaskName
	IntervalSchedule  = scheduler.IntervalSchedule
	Signature         = workflow.Signature
	Chain             = workflow.Chain
	Group             = workflow.Group
	Chord             = workflow.Chord
	Canvas            = workflow.Canvas
	CanvasOption      = workflow.CanvasOption
	ChainResult       = workflow.ChainResult
	GroupResult       = workflow.GroupResult
	ChordResult       = workflow.ChordResult
	Worker            = worker.Worker
	WorkerOption      = worker.WorkerOption
)

var (
	ErrInvalidTaskName     = task.ErrInvalidTaskName
	ErrInvalidTaskID       = task.ErrInvalidTaskID
	ErrInvalidPriority     = task.ErrInvalidPriority
	ErrInvalidTaskState    = task.ErrInvalidTaskState
	ErrInvalidRetryPolicy  = task.ErrInvalidRetryPolicy
	ErrInvalidTaskTiming   = task.ErrInvalidTaskTiming
	ErrInvalidPayload      = task.ErrInvalidPayload
	ErrDuplicateTask       = task.ErrDuplicateTask
	ErrInvalidTaskHandler  = task.ErrInvalidTaskHandler
	ErrTaskNotFound        = task.ErrTaskNotFound
	ErrNilBackend          = producer.ErrNilBackend
	ErrNilWorker           = worker.ErrNilWorker
	ErrNilTaskRegistry     = worker.ErrNilTaskRegistry
	ErrMissingTaskName     = producer.ErrMissingTaskName
	ErrMissingApplyOption  = producer.ErrMissingApplyOption
	ErrInvalidSchedule     = scheduler.ErrInvalidSchedule
	ErrInvalidPeriodicTask = scheduler.ErrInvalidPeriodicTask
	ErrNilSchedulerBackend = scheduler.ErrNilBackend
	ErrInvalidWorkflow     = workflow.ErrInvalidWorkflow
	ErrInvalidSignature    = workflow.ErrInvalidSignature
	ErrNilCanvasBackend    = workflow.ErrNilBackend
	ErrInvalidWorkerOption = worker.ErrInvalidWorkerOption
)

const (
	ScheduleKindInterval             = scheduler.ScheduleKindInterval
	PeriodicMetadataNameKey          = scheduler.PeriodicMetadataNameKey
	PeriodicMetadataDueAtKey         = scheduler.PeriodicMetadataDueAtKey
	WorkflowKindChain                = workflow.WorkflowKindChain
	WorkflowKindGroup                = workflow.WorkflowKindGroup
	WorkflowKindChord                = workflow.WorkflowKindChord
	WorkflowMetadataKindKey          = workflow.MetadataKindKey
	WorkflowMetadataChainIDKey       = workflow.MetadataChainIDKey
	WorkflowMetadataChainStepKey     = workflow.MetadataChainStepKey
	WorkflowMetadataGroupIDKey       = workflow.MetadataGroupIDKey
	WorkflowMetadataGroupIndexKey    = workflow.MetadataGroupIndexKey
	WorkflowMetadataChordIDKey       = workflow.MetadataChordIDKey
	WorkflowMetadataChordCallbackKey = workflow.MetadataChordCallbackKey
)

const (
	FailureExecution           = task.FailureExecution
	FailureMalformedMessage    = task.FailureMalformedMessage
	FailureUnknownTask         = task.FailureUnknownTask
	FailureExpired             = task.FailureExpired
	FailureRetryExhausted      = task.FailureRetryExhausted
	FailureRetryScheduleFailed = task.FailureRetryScheduleFailed
)

const (
	FailureMetadataCategoryKey       = task.FailureMetadataCategoryKey
	FailureMetadataAttemptKey        = task.FailureMetadataAttemptKey
	FailureMetadataMaxAttemptsKey    = task.FailureMetadataMaxAttemptsKey
	FailureMetadataRetryableKey      = task.FailureMetadataRetryableKey
	FailureMetadataNextRetryAtKey    = task.FailureMetadataNextRetryAtKey
	FailureMetadataDeadLetteredKey   = task.FailureMetadataDeadLetteredKey
	FailureMetadataDeadLetteredAtKey = task.FailureMetadataDeadLetteredAtKey
	FailureMetadataLastErrorKey      = task.FailureMetadataLastErrorKey
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

func WithProducerDefaultQueue(queue QueueName) ProducerOption {
	return producer.WithProducerDefaultQueue(queue)
}

func WithProducerCodec(codec PayloadCodec) ProducerOption {
	return producer.WithProducerCodec(codec)
}

func WithProducerNow(now func() time.Time) ProducerOption {
	return producer.WithProducerNow(now)
}

func WithApplyQueue(queue QueueName) ApplyOption {
	return producer.WithApplyQueue(queue)
}

func WithApplyTaskID(id TaskID) ApplyOption {
	return producer.WithApplyTaskID(id)
}

func WithApplyMetadata(metadata map[string]string) ApplyOption {
	return producer.WithApplyMetadata(metadata)
}

func WithApplyPriority(priority Priority) ApplyOption {
	return producer.WithApplyPriority(priority)
}

func WithApplyRetryPolicy(policy RetryPolicy) ApplyOption {
	return producer.WithApplyRetryPolicy(policy)
}

func WithApplyCountDown(countDown time.Duration) ApplyOption {
	return producer.WithApplyCountDown(countDown)
}

func WithApplyETA(eta time.Time) ApplyOption {
	return producer.WithApplyETA(eta)
}

func WithApplyExpiresAt(expiresAt time.Time) ApplyOption {
	return producer.WithApplyExpiresAt(expiresAt)
}

func WithApplyAttempt(attempt int) ApplyOption {
	return producer.WithApplyAttempt(attempt)
}

func WithApplyCreatedAt(createdAt time.Time) ApplyOption {
	return producer.WithApplyCreatedAt(createdAt)
}

func Every(interval time.Duration) IntervalSchedule {
	return scheduler.Every(interval)
}

func WithSchedulerIdentity(identity string) SchedulerOption {
	return scheduler.WithSchedulerIdentity(identity)
}

func WithSchedulerDefaultQueue(queue QueueName) SchedulerOption {
	return scheduler.WithSchedulerDefaultQueue(queue)
}

func WithSchedulerPollInterval(interval time.Duration) SchedulerOption {
	return scheduler.WithSchedulerPollInterval(interval)
}

func WithSchedulerBatchSize(size int64) SchedulerOption {
	return scheduler.WithSchedulerBatchSize(size)
}

func WithSchedulerLockTTL(ttl time.Duration) SchedulerOption {
	return scheduler.WithSchedulerLockTTL(ttl)
}

func WithSchedulerCodec(codec PayloadCodec) SchedulerOption {
	return scheduler.WithSchedulerCodec(codec)
}

func WithSchedulerNow(now func() time.Time) SchedulerOption {
	return scheduler.WithSchedulerNow(now)
}

func WithCanvasDefaultQueue(queue QueueName) CanvasOption {
	return workflow.WithCanvasDefaultQueue(queue)
}

func WithCanvasCodec(codec PayloadCodec) CanvasOption {
	return workflow.WithCanvasCodec(codec)
}

func WithCanvasNow(now func() time.Time) CanvasOption {
	return workflow.WithCanvasNow(now)
}

func WithWorkerQueue(queue QueueName) WorkerOption {
	return worker.WithWorkerQueue(queue)
}

func WithWorkerGroup(group string) WorkerOption {
	return worker.WithWorkerGroup(group)
}

func WithWorkerConsumer(consumer string) WorkerOption {
	return worker.WithWorkerConsumer(consumer)
}

func WithWorkerCodec(codec PayloadCodec) WorkerOption {
	return worker.WithWorkerCodec(codec)
}

func WithWorkerConcurrency(concurrency int) WorkerOption {
	return worker.WithWorkerConcurrency(concurrency)
}

func WithWorkerReadBatch(readBatch int64) WorkerOption {
	return worker.WithWorkerReadBatch(readBatch)
}

func WithWorkerBlock(block time.Duration) WorkerOption {
	return worker.WithWorkerBlock(block)
}

func WithWorkerMoveDueLimit(limit int64) WorkerOption {
	return worker.WithWorkerMoveDueLimit(limit)
}

func WithWorkerMoveDueEnabled(enabled bool) WorkerOption {
	return worker.WithWorkerMoveDueEnabled(enabled)
}

func WithWorkerIdleDelay(delay time.Duration) WorkerOption {
	return worker.WithWorkerIdleDelay(delay)
}

func WithWorkerDeadLetterEnabled(enabled bool) WorkerOption {
	return worker.WithWorkerDeadLetterEnabled(enabled)
}

func WithWorkerPendingRecoveryEnabled(enabled bool) WorkerOption {
	return worker.WithWorkerPendingRecoveryEnabled(enabled)
}

func WithWorkerPendingMinIdle(minIdle time.Duration) WorkerOption {
	return worker.WithWorkerPendingMinIdle(minIdle)
}

func WithWorkerPendingClaimBatch(count int64) WorkerOption {
	return worker.WithWorkerPendingClaimBatch(count)
}

func WithWorkerPendingClaimInterval(interval time.Duration) WorkerOption {
	return worker.WithWorkerPendingClaimInterval(interval)
}

func WithWorkerNow(now func() time.Time) WorkerOption {
	return worker.WithWorkerNow(now)
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
