package workflow

import (
	"github.com/Newton-School/goqueue/task"
)

// Signature is a reusable task invocation for workflow composition.
type Signature struct {
	Name        task.TaskName
	Queue       task.QueueName
	Args        []any
	Kwargs      map[string]any
	Metadata    map[string]string
	Timing      task.TaskTiming
	Priority    task.Priority
	RetryPolicy task.RetryPolicy
}

// Validate verifies that the signature can be dispatched as a task.
func (s Signature) Validate() error {
	if err := task.ValidateTaskName(s.Name.String()); err != nil {
		return err
	}
	if s.Queue != "" {
		if err := task.ValidateQueueName(s.Queue.String()); err != nil {
			return err
		}
	}
	if err := s.Timing.Validate(); err != nil {
		return err
	}
	if s.Priority != 0 {
		if err := task.ValidatePriority(s.Priority); err != nil {
			return err
		}
	}
	if s.RetryPolicy.MaxAttempts != 0 {
		if err := s.RetryPolicy.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Normalize applies default queue, priority, and retry policy and copies mutable fields.
func (s Signature) Normalize(defaultQueue task.QueueName) (Signature, error) {
	normalized := s
	if normalized.Queue == "" {
		normalized.Queue = defaultQueue
	}
	if normalized.Priority == 0 {
		normalized.Priority = task.DefaultPriority
	}
	if normalized.RetryPolicy.MaxAttempts == 0 {
		normalized.RetryPolicy = task.DefaultRetryPolicy()
	}
	normalized.Args = copyAnySlice(s.Args)
	normalized.Kwargs = copyAnyMap(s.Kwargs)
	normalized.Metadata = copyStringMap(s.Metadata)

	if err := normalized.Validate(); err != nil {
		return Signature{}, err
	}

	return normalized, nil
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
