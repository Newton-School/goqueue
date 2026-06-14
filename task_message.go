package goqueue

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

func cloneBytes(values []byte) []byte {
	if len(values) == 0 {
		return nil
	}

	cloned := make([]byte, len(values))
	copy(cloned, values)
	return cloned
}
