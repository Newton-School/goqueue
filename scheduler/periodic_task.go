package scheduler

import (
	"fmt"
	"time"

	"github.com/Newton-School/goqueue/task"
)

// PeriodicTaskName uniquely identifies a scheduler entry.
type PeriodicTaskName string

// String returns the periodic task name as a string.
func (n PeriodicTaskName) String() string {
	return string(n)
}

// PeriodicTask describes a task the scheduler should enqueue repeatedly.
type PeriodicTask struct {
	Name        PeriodicTaskName
	TaskName    task.TaskName
	Queue       task.QueueName
	Args        []any
	Kwargs      map[string]any
	Metadata    map[string]string
	Schedule    IntervalSchedule
	StartAt     time.Time
	Priority    task.Priority
	RetryPolicy task.RetryPolicy
}

// Validate verifies that the periodic task can be stored and dispatched.
func (p PeriodicTask) Validate() error {
	if err := validatePeriodicTaskName(p.Name.String()); err != nil {
		return err
	}
	if err := task.ValidateTaskName(p.TaskName.String()); err != nil {
		return err
	}
	if p.Queue != "" {
		if err := task.ValidateQueueName(p.Queue.String()); err != nil {
			return err
		}
	}
	if err := p.Schedule.Validate(); err != nil {
		return err
	}
	if p.Priority != 0 {
		if err := task.ValidatePriority(p.Priority); err != nil {
			return err
		}
	}
	if p.RetryPolicy.MaxAttempts != 0 {
		if err := p.RetryPolicy.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func validatePeriodicTaskName(name string) error {
	if err := task.ValidateTaskName(name); err != nil {
		return fmt.Errorf("%w: name must use 1-128 characters from [A-Za-z0-9._:-]", ErrInvalidPeriodicTask)
	}

	return nil
}
