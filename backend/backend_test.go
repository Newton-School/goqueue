package backend

import (
	"context"
	"testing"

	"github.com/Newton-School/goqueue/task"
)

func TestQueueBackendInterfaceAcceptsImplementation(t *testing.T) {
	var backend QueueBackend = noopBackend{}
	if backend == nil {
		t.Fatal("QueueBackend should accept implementations")
	}
}

type noopBackend struct{}

func (noopBackend) EnqueueReady(context.Context, EnqueueRequest) (EnqueueResponse, error) {
	return EnqueueResponse{}, nil
}
func (noopBackend) EnqueueScheduled(context.Context, EnqueueRequest) (EnqueueResponse, error) {
	return EnqueueResponse{}, nil
}
func (noopBackend) MoveDueScheduled(context.Context, MoveDueScheduledRequest) ([]MovedScheduledMessage, error) {
	return nil, nil
}
func (noopBackend) ReadReady(context.Context, ReadReadyRequest) ([]ReadyMessage, error) {
	return nil, nil
}
func (noopBackend) Ack(context.Context, AckRequest) error { return nil }
func (noopBackend) SetTaskState(context.Context, TaskStateRecord) error {
	return nil
}
func (noopBackend) GetTaskState(context.Context, task.TaskID) (TaskStateRecord, error) {
	return TaskStateRecord{}, nil
}
func (noopBackend) SaveTaskResult(context.Context, TaskResultRecord) error {
	return nil
}
func (noopBackend) GetTaskResult(context.Context, task.TaskID) (TaskResultRecord, error) {
	return TaskResultRecord{}, nil
}
func (noopBackend) ForgetTaskResult(context.Context, task.TaskID) error {
	return nil
}
func (noopBackend) QueueStats(context.Context, QueueStatsRequest) (QueueStats, error) {
	return QueueStats{}, nil
}
func (noopBackend) Ping(context.Context) error { return nil }
func (noopBackend) Close() error               { return nil }
