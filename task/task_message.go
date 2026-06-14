package task

import "time"

// TaskMessage is the serialized form future backends store and deliver.
type TaskMessage struct {
	ID          string
	Name        string
	Queue       string
	Payload     []byte
	Metadata    map[string]string
	Timing      TaskTiming
	Priority    Priority
	RetryPolicy RetryPolicy
	CreatedAt   time.Time
	Attempt     int
}

// TaskEnvelopeToMessage serializes an envelope payload for backend storage.
func TaskEnvelopeToMessage(envelope TaskEnvelope, codec PayloadCodec) (TaskMessage, error) {
	payload, err := codec.EncodePayload(envelope.Payload)
	if err != nil {
		return TaskMessage{}, err
	}

	return TaskMessage{
		ID:          envelope.ID.String(),
		Name:        envelope.Name.String(),
		Queue:       envelope.Queue.String(),
		Payload:     cloneBytes(payload),
		Metadata:    envelope.Metadata.Values(),
		Timing:      envelope.Timing,
		Priority:    envelope.Priority,
		RetryPolicy: envelope.RetryPolicy,
		CreatedAt:   envelope.CreatedAt,
		Attempt:     envelope.Attempt,
	}, nil
}

// TaskMessageToEnvelope decodes a backend message into a validated envelope.
func TaskMessageToEnvelope(message TaskMessage, codec PayloadCodec) (TaskEnvelope, error) {
	payload, err := codec.DecodePayload(message.Payload)
	if err != nil {
		return TaskEnvelope{}, err
	}

	return NewTaskEnvelope(TaskEnvelopeInput{
		ID:          TaskID(message.ID),
		Name:        TaskName(message.Name),
		Queue:       QueueName(message.Queue),
		Args:        payload.Args(),
		Kwargs:      payload.Kwargs(),
		Metadata:    message.Metadata,
		Timing:      message.Timing,
		Priority:    message.Priority,
		RetryPolicy: message.RetryPolicy,
		CreatedAt:   message.CreatedAt,
		Attempt:     message.Attempt,
	})
}

func cloneBytes(values []byte) []byte {
	if len(values) == 0 {
		return nil
	}

	cloned := make([]byte, len(values))
	copy(cloned, values)
	return cloned
}
