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

// Validate verifies that the upsert request contains a complete record.
func (r UpsertPeriodicTaskRequest) Validate() error {
	return r.Record.Validate()
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
