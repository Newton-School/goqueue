package backend

import (
	"fmt"
	"time"

	"github.com/Newton-School/goqueue/task"
)

const (
	// PeriodicScheduleInterval identifies fixed interval schedule records.
	PeriodicScheduleInterval = "interval"
)

// PeriodicTaskRecord is the backend storage form of a periodic task definition.
type PeriodicTaskRecord struct {
	Name         string
	TaskName     task.TaskName
	Queue        task.QueueName
	Args         []any
	Kwargs       map[string]any
	Metadata     map[string]string
	ScheduleKind string
	Interval     time.Duration
	StartAt      time.Time
	NextDueAt    time.Time
	Priority     task.Priority
	RetryPolicy  task.RetryPolicy
	UpdatedAt    time.Time
}

// UpsertPeriodicTaskRequest stores or replaces a periodic task definition.
type UpsertPeriodicTaskRequest struct {
	Record PeriodicTaskRecord
}

// ListDuePeriodicTasksRequest leases due periodic definitions for a scheduler.
type ListDuePeriodicTasksRequest struct {
	Now         time.Time
	Limit       int64
	SchedulerID string
	LockTTL     time.Duration
}

// DuePeriodicTask is a due periodic definition leased by one scheduler.
type DuePeriodicTask struct {
	Record      PeriodicTaskRecord
	LockToken   string
	LockedUntil time.Time
}

// MarkPeriodicTaskDispatchedRequest advances a periodic definition after dispatch.
type MarkPeriodicTaskDispatchedRequest struct {
	Name             string
	LockToken        string
	DispatchedTaskID task.TaskID
	DispatchedAt     time.Time
	NextDueAt        time.Time
}

// Validate verifies that the upsert request contains a complete record.
func (r UpsertPeriodicTaskRequest) Validate() error {
	return r.Record.Validate()
}

// Validate verifies that due periodic scans can be bounded and leased.
func (r ListDuePeriodicTasksRequest) Validate() error {
	if r.Now.IsZero() {
		return fmt.Errorf("%w: scan time is required", ErrInvalidBackendRequest)
	}
	if r.Limit <= 0 {
		return fmt.Errorf("%w: due scan limit must be positive", ErrInvalidBackendRequest)
	}
	if r.SchedulerID == "" {
		return fmt.Errorf("%w: scheduler id is required", ErrInvalidBackendRequest)
	}
	if r.LockTTL <= 0 {
		return fmt.Errorf("%w: lock ttl must be positive", ErrInvalidBackendRequest)
	}

	return nil
}

// Validate verifies that the due task is leased and dispatchable.
func (d DuePeriodicTask) Validate() error {
	if err := d.Record.Validate(); err != nil {
		return err
	}
	if d.LockToken == "" {
		return fmt.Errorf("%w: lock token is required", ErrInvalidBackendRequest)
	}
	if d.LockedUntil.IsZero() {
		return fmt.Errorf("%w: locked until is required", ErrInvalidBackendRequest)
	}

	return nil
}

// Validate verifies that dispatch marking is tied to a held lease and task ID.
func (r MarkPeriodicTaskDispatchedRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("%w: periodic task name is required", ErrInvalidBackendRequest)
	}
	if err := task.ValidateTaskName(r.Name); err != nil {
		return fmt.Errorf("%w: periodic task name is invalid: %v", ErrInvalidBackendRequest, err)
	}
	if r.LockToken == "" {
		return fmt.Errorf("%w: lock token is required", ErrInvalidBackendRequest)
	}
	if err := task.ValidateTaskID(r.DispatchedTaskID.String()); err != nil {
		return err
	}
	if r.DispatchedAt.IsZero() {
		return fmt.Errorf("%w: dispatched at is required", ErrInvalidBackendRequest)
	}
	if r.NextDueAt.IsZero() {
		return fmt.Errorf("%w: next due time is required", ErrInvalidBackendRequest)
	}

	return nil
}

// Validate verifies that a periodic task record is safe for backend storage.
func (r PeriodicTaskRecord) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("%w: periodic task name is required", ErrInvalidBackendRequest)
	}
	if err := task.ValidateTaskName(r.Name); err != nil {
		return fmt.Errorf("%w: periodic task name is invalid: %v", ErrInvalidBackendRequest, err)
	}
	if err := task.ValidateTaskName(r.TaskName.String()); err != nil {
		return err
	}
	if err := task.ValidateQueueName(r.Queue.String()); err != nil {
		return err
	}
	if r.ScheduleKind != PeriodicScheduleInterval {
		return fmt.Errorf("%w: unsupported periodic schedule kind %q", ErrInvalidBackendRequest, r.ScheduleKind)
	}
	if r.Interval <= 0 {
		return fmt.Errorf("%w: periodic interval must be positive", ErrInvalidBackendRequest)
	}
	if r.NextDueAt.IsZero() {
		return fmt.Errorf("%w: next due time is required", ErrInvalidBackendRequest)
	}
	if err := task.ValidatePriority(r.Priority); err != nil {
		return err
	}
	if err := r.RetryPolicy.Validate(); err != nil {
		return err
	}

	return nil
}
