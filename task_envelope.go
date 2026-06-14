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

	envelope := TaskEnvelope{
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
	}

	if err := envelope.Validate(); err != nil {
		return TaskEnvelope{}, err
	}

	return envelope, nil
}

// Validate verifies that the envelope is safe for queue storage and execution.
func (e TaskEnvelope) Validate() error {
	if err := ValidateTaskID(e.ID.String()); err != nil {
		return err
	}

	if err := ValidateTaskName(e.Name.String()); err != nil {
		return err
	}

	if err := ValidateQueueName(e.Queue.String()); err != nil {
		return err
	}

	if err := ValidatePriority(e.Priority); err != nil {
		return err
	}

	if err := e.RetryPolicy.Validate(); err != nil {
		return err
	}

	if err := e.Timing.Validate(); err != nil {
		return err
	}

	return nil
}

// Clone returns a copy of the envelope with copied mutable payload and metadata.
func (e TaskEnvelope) Clone() TaskEnvelope {
	e.Payload = NewTaskPayload(e.Payload.Args(), e.Payload.Kwargs())
	e.Metadata = NewTaskMetadata(e.Metadata.Values())
	return e
}
