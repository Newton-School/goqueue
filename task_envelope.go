package goqueue

import "time"

// TaskEnvelope is the complete SDK-level representation of a task invocation.
type TaskEnvelope struct {
	ID          TaskID
	Name        TaskName
	Queue       QueueName
	Payload     TaskPayload
	Metadata    TaskMetadata
	Timing      TaskTiming
	Priority    Priority
	RetryPolicy RetryPolicy
	CreatedAt   time.Time
	Attempt     int
}

// TaskEnvelopeInput contains fields used to create a task envelope.
type TaskEnvelopeInput struct {
	ID          TaskID
	Name        TaskName
	Queue       QueueName
	Args        []any
	Kwargs      map[string]any
	Metadata    map[string]string
	Timing      TaskTiming
	Priority    Priority
	RetryPolicy RetryPolicy
	CreatedAt   time.Time
	Attempt     int
}

// NewTaskEnvelope creates a task envelope with generated IDs and safe defaults.
func NewTaskEnvelope(input TaskEnvelopeInput) (TaskEnvelope, error) {
	id := input.ID
	if id == "" {
		generated, err := NewTaskID()
		if err != nil {
			return TaskEnvelope{}, err
		}
		id = generated
	}

	priority := input.Priority
	if priority == 0 {
		priority = DefaultPriority
	}

	retryPolicy := input.RetryPolicy
	if retryPolicy.MaxAttempts == 0 {
		retryPolicy = DefaultRetryPolicy()
	}

	createdAt := input.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}

	return TaskEnvelope{
		ID:          id,
		Name:        input.Name,
		Queue:       input.Queue,
		Payload:     NewTaskPayload(input.Args, input.Kwargs),
		Metadata:    NewTaskMetadata(input.Metadata),
		Timing:      input.Timing,
		Priority:    priority,
		RetryPolicy: retryPolicy,
		CreatedAt:   createdAt,
		Attempt:     input.Attempt,
	}, nil
}
