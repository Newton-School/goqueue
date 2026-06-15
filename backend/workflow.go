package backend

import (
	"fmt"
	"time"

	"github.com/Newton-School/goqueue/task"
)

// WorkflowChainRecord stores an ordered chain of workflow signatures.
type WorkflowChainRecord struct {
	ID         string
	Signatures []WorkflowSignatureRecord
	CreatedAt  time.Time
}

// AdvanceWorkflowChainRequest records a completed chain step.
type AdvanceWorkflowChainRequest struct {
	WorkflowID      string
	CompletedTaskID task.TaskID
	CompletedIndex  int
	CompletedAt     time.Time
}

// AdvanceWorkflowChainResponse describes an idempotent chain advancement.
type AdvanceWorkflowChainResponse struct {
	Advanced  bool
	Completed bool
	Next      *WorkflowSignatureRecord
}

// WorkflowGroupRecord stores task membership for a group or chord header.
type WorkflowGroupRecord struct {
	ID        string
	TaskIDs   []task.TaskID
	Callback  *WorkflowSignatureRecord
	CreatedAt time.Time
}

// RecordWorkflowTaskCompletedRequest records a terminal group child state.
type RecordWorkflowTaskCompletedRequest struct {
	GroupID     string
	TaskID      task.TaskID
	State       task.TaskState
	CompletedAt time.Time
}

// WorkflowGroupProgress describes an idempotent group progress update.
type WorkflowGroupProgress struct {
	GroupID   string
	Total     int
	Completed int
	Failed    int
	Duplicate bool
	Done      bool
	Succeeded bool
	Callback  *WorkflowSignatureRecord
}

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

// Validate verifies that a chain record is safe to store.
func (r WorkflowChainRecord) Validate() error {
	if err := validateWorkflowID(r.ID); err != nil {
		return err
	}
	if len(r.Signatures) == 0 {
		return fmt.Errorf("%w: chain requires at least one signature", ErrInvalidBackendRequest)
	}
	for index, signature := range r.Signatures {
		if err := signature.Validate(); err != nil {
			return fmt.Errorf("%w: chain signature %d: %v", ErrInvalidBackendRequest, index, err)
		}
	}
	if r.CreatedAt.IsZero() {
		return fmt.Errorf("%w: chain created at is required", ErrInvalidBackendRequest)
	}

	return nil
}

// Validate verifies that chain advancement identifies a completed step.
func (r AdvanceWorkflowChainRequest) Validate() error {
	if err := validateWorkflowID(r.WorkflowID); err != nil {
		return err
	}
	if err := task.ValidateTaskID(r.CompletedTaskID.String()); err != nil {
		return err
	}
	if r.CompletedIndex < 0 {
		return fmt.Errorf("%w: completed index cannot be negative", ErrInvalidBackendRequest)
	}
	if r.CompletedAt.IsZero() {
		return fmt.Errorf("%w: completed at is required", ErrInvalidBackendRequest)
	}

	return nil
}

// Validate verifies that a group record is safe to store.
func (r WorkflowGroupRecord) Validate() error {
	if err := validateWorkflowID(r.ID); err != nil {
		return err
	}
	if len(r.TaskIDs) == 0 {
		return fmt.Errorf("%w: group requires at least one task id", ErrInvalidBackendRequest)
	}
	for index, taskID := range r.TaskIDs {
		if err := task.ValidateTaskID(taskID.String()); err != nil {
			return fmt.Errorf("%w: group task id %d: %v", ErrInvalidBackendRequest, index, err)
		}
	}
	if r.Callback != nil {
		if err := r.Callback.Validate(); err != nil {
			return fmt.Errorf("%w: group callback: %v", ErrInvalidBackendRequest, err)
		}
	}
	if r.CreatedAt.IsZero() {
		return fmt.Errorf("%w: group created at is required", ErrInvalidBackendRequest)
	}

	return nil
}

// Validate verifies that a group progress request records a terminal state.
func (r RecordWorkflowTaskCompletedRequest) Validate() error {
	if err := validateWorkflowID(r.GroupID); err != nil {
		return err
	}
	if err := task.ValidateTaskID(r.TaskID.String()); err != nil {
		return err
	}
	if !isTerminalWorkflowState(r.State) {
		return fmt.Errorf("%w: state %q is not terminal", ErrInvalidBackendRequest, r.State)
	}
	if r.CompletedAt.IsZero() {
		return fmt.Errorf("%w: completed at is required", ErrInvalidBackendRequest)
	}

	return nil
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

func isTerminalWorkflowState(state task.TaskState) bool {
	switch state {
	case task.TaskSucceeded, task.TaskFailed, task.TaskRevoked, task.TaskExpired, task.TaskDeadLettered:
		return true
	default:
		return false
	}
}

func validateWorkflowID(id string) error {
	if id == "" {
		return fmt.Errorf("%w: workflow id is required", ErrInvalidBackendRequest)
	}
	if err := task.ValidateTaskName(id); err != nil {
		return fmt.Errorf("%w: workflow id is invalid: %v", ErrInvalidBackendRequest, err)
	}

	return nil
}
