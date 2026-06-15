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

// Normalize applies scheduler defaults and copies mutable fields.
func (p PeriodicTask) Normalize(defaultQueue task.QueueName) (PeriodicTask, error) {
	normalized := p
	if normalized.Queue == "" {
		normalized.Queue = defaultQueue
	}
	if normalized.Priority == 0 {
		normalized.Priority = task.DefaultPriority
	}
	if normalized.RetryPolicy.MaxAttempts == 0 {
		normalized.RetryPolicy = task.DefaultRetryPolicy()
	}

	normalized.Args = copyAnySlice(p.Args)
	normalized.Kwargs = copyAnyMap(p.Kwargs)
	normalized.Metadata = copyStringMap(p.Metadata)

	if err := normalized.Validate(); err != nil {
		return PeriodicTask{}, err
	}

	return normalized, nil
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

func copyAnySlice(values []any) []any {
	if values == nil {
		return nil
	}

	copied := make([]any, len(values))
	copy(copied, values)
	return copied
}

func copyAnyMap(values map[string]any) map[string]any {
	if values == nil {
		return nil
	}

	copied := make(map[string]any, len(values))
	for key, value := range values {
		copied[key] = value
	}

	return copied
}

func copyStringMap(values map[string]string) map[string]string {
	if values == nil {
		return nil
	}

	copied := make(map[string]string, len(values))
	for key, value := range values {
		copied[key] = value
	}

	return copied
}

func validatePeriodicTaskName(name string) error {
	if err := task.ValidateTaskName(name); err != nil {
		return fmt.Errorf("%w: name must use 1-128 characters from [A-Za-z0-9._:-]", ErrInvalidPeriodicTask)
	}

	return nil
}
