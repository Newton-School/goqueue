package admin

import (
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

// RetryTaskOptions control how an existing persisted task message is re-queued.
type RetryTaskOptions struct {
	Queue           task.QueueName
	ScheduledAt     time.Time
	CountDown       time.Duration
	PreserveAttempt bool
	ClearState      bool
	ClearResult     bool
	Now             func() time.Time
}

// RetryTaskResult summarizes retry dispatch.
type RetryTaskResult struct {
	TaskID        task.TaskID
	Queue         task.QueueName
	OriginalQueue task.QueueName
	ScheduledAt   time.Time
	EnqueueResult backend.EnqueueResponse
	Attempt       int
}

// RevokeTaskResult summarizes revoke operation output.
type RevokeTaskResult struct {
	TaskID task.TaskID
	State  task.TaskState
}

// ReplayDeadLetterOptions controls DLQ replay behavior.
type ReplayDeadLetterOptions struct {
	DestinationQueue task.QueueName
	DeleteSource     bool
}

// ReplayDeadLetterResult summarizes replay action output.
type ReplayDeadLetterResult struct {
	StreamID       string
	Queue          task.QueueName
	Destination    task.QueueName
	EnqueueResult  backend.EnqueueResponse
	SourceDeleted  bool
	OriginalTaskID task.TaskID
}

// PurgeQueueOptions controls optional deletes during queue purge.
type PurgeQueueOptions struct {
	Queue         task.QueueName
	DeleteMessage bool
	DeleteState   bool
	DeleteResult  bool
}

// PurgeQueueResult reports queue purge status.
type PurgeQueueResult struct {
	Queue            task.QueueName
	ReadyStream      int64
	ScheduledSet     int64
	DeadLetterStream int64
	TaskMessages     int64
	TaskStates       int64
	TaskResults      int64
}

// DeleteDeadLettersResult summarizes delete action output.
type DeleteDeadLettersResult struct {
	Queue   task.QueueName
	Deleted int64
}
