package backend

import (
	"github.com/Newton-School/goqueue/task"
)

// WorkflowSignatureRecord is the backend storage form of a workflow task signature.
type WorkflowSignatureRecord struct {
	Name        task.TaskName
	Queue       task.QueueName
	Args        []any
	Kwargs      map[string]any
	Metadata    map[string]string
	Timing      task.TaskTiming
	Priority    task.Priority
	RetryPolicy task.RetryPolicy
}

// Validate verifies that a workflow signature record is dispatchable.
func (r WorkflowSignatureRecord) Validate() error {
	if err := task.ValidateTaskName(r.Name.String()); err != nil {
		return err
	}
	if err := task.ValidateQueueName(r.Queue.String()); err != nil {
		return err
	}
	if err := r.Timing.Validate(); err != nil {
		return err
	}
	if err := task.ValidatePriority(r.Priority); err != nil {
		return err
	}
	if err := r.RetryPolicy.Validate(); err != nil {
		return err
	}

	return nil
}
