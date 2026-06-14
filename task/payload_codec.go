package task

// PayloadCodec serializes and deserializes task payloads.
type PayloadCodec interface {
	EncodePayload(TaskPayload) ([]byte, error)
	DecodePayload([]byte) (TaskPayload, error)
}
