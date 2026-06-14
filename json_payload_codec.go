package goqueue

import (
	"encoding/json"
	"fmt"
)

// JSONPayloadCodec encodes task payloads as JSON.
type JSONPayloadCodec struct{}

type jsonPayload struct {
	Args   []any          `json:"args"`
	Kwargs map[string]any `json:"kwargs"`
}

// EncodePayload serializes payload into a stable JSON object.
func (JSONPayloadCodec) EncodePayload(payload TaskPayload) ([]byte, error) {
	encoded, err := json.Marshal(jsonPayload{
		Args:   payload.Args(),
		Kwargs: payload.Kwargs(),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: encode json payload: %v", ErrInvalidPayload, err)
	}

	return encoded, nil
}

// DecodePayload deserializes a JSON payload object.
func (JSONPayloadCodec) DecodePayload(data []byte) (TaskPayload, error) {
	var decoded jsonPayload
	if err := json.Unmarshal(data, &decoded); err != nil {
		return TaskPayload{}, fmt.Errorf("%w: decode json payload: %v", ErrInvalidPayload, err)
	}

	return NewTaskPayload(decoded.Args, decoded.Kwargs), nil
}
