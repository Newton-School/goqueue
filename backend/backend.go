package backend

import (
	"context"

	"github.com/Newton-School/goqueue/task"
)

// QueueBackend is the storage boundary used by producers, workers, and schedulers.
type QueueBackend interface {
	EnqueueReady(context.Context, EnqueueRequest) (EnqueueResponse, error)
	EnqueueScheduled(context.Context, EnqueueRequest) (EnqueueResponse, error)
	MoveDueScheduled(context.Context, MoveDueScheduledRequest) ([]MovedScheduledMessage, error)
	EnsureConsumerGroup(context.Context, ConsumerGroupRequest) error
	ReadReady(context.Context, ReadReadyRequest) ([]ReadyMessage, error)
	ClaimStaleReady(context.Context, ClaimStaleReadyRequest) ([]ReadyMessage, error)
	Ack(context.Context, AckRequest) error
	EnqueueDeadLetter(context.Context, DeadLetterRequest) (DeadLetterRecord, error)
	ReadDeadLetters(context.Context, ReadDeadLettersRequest) ([]DeadLetterRecord, error)
	UpsertPeriodicTask(context.Context, UpsertPeriodicTaskRequest) error
	DeletePeriodicTask(context.Context, DeletePeriodicTaskRequest) error
	ListDuePeriodicTasks(context.Context, ListDuePeriodicTasksRequest) ([]DuePeriodicTask, error)
	MarkPeriodicTaskDispatched(context.Context, MarkPeriodicTaskDispatchedRequest) error
	SetTaskState(context.Context, TaskStateRecord) error
	GetTaskState(context.Context, task.TaskID) (TaskStateRecord, error)
	SaveTaskResult(context.Context, TaskResultRecord) error
	GetTaskResult(context.Context, task.TaskID) (TaskResultRecord, error)
	ForgetTaskResult(context.Context, task.TaskID) error
	QueueStats(context.Context, QueueStatsRequest) (QueueStats, error)
	Ping(context.Context) error
	Close() error
}
