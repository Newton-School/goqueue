package redisbackend

import (
	"fmt"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

type deadLetterCodec struct{}

func (c deadLetterCodec) encode(record backend.DeadLetterRecord) (map[string]any, error) {
	encoded, err := (messageCodec{}).encode(record.Message)
	if err != nil {
		return nil, err
	}

	values := map[string]any{
		"message":          string(encoded),
		"reason":           string(record.Reason),
		"error":            record.Error,
		"source_stream_id": record.SourceStreamID,
		"group":            record.Group,
		"consumer":         record.Consumer,
		"failed_at":        record.FailedAt.UTC().Format(time.RFC3339Nano),
	}
	return values, nil
}

func (c deadLetterCodec) decode(streamID string, values map[string]any) (backend.DeadLetterRecord, error) {
	rawMessage, ok := values["message"]
	if !ok {
		return backend.DeadLetterRecord{}, fmt.Errorf("%w: dead letter missing message field", ErrInvalidRedisMessage)
	}
	encodedMessage, ok := rawMessage.(string)
	if !ok {
		return backend.DeadLetterRecord{}, fmt.Errorf("%w: dead letter message field must be string", ErrInvalidRedisMessage)
	}
	message, err := (messageCodec{}).decode([]byte(encodedMessage))
	if err != nil {
		return backend.DeadLetterRecord{}, err
	}

	failedAt, err := parseDeadLetterTime(values["failed_at"])
	if err != nil {
		return backend.DeadLetterRecord{}, err
	}

	return backend.DeadLetterRecord{
		StreamID:       streamID,
		Message:        message,
		Reason:         task.FailureCategory(stringValue(values["reason"])),
		Error:          stringValue(values["error"]),
		SourceStreamID: stringValue(values["source_stream_id"]),
		Group:          stringValue(values["group"]),
		Consumer:       stringValue(values["consumer"]),
		FailedAt:       failedAt,
	}, nil
}

func parseDeadLetterTime(value any) (time.Time, error) {
	raw := stringValue(value)
	if raw == "" {
		return time.Time{}, nil
	}
	parsed, err := time.Parse(time.RFC3339Nano, raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: invalid dead letter timestamp", ErrInvalidRedisMessage)
	}
	return parsed, nil
}

func stringValue(value any) string {
	typed, ok := value.(string)
	if !ok {
		return ""
	}
	return typed
}
