package redisbackend

import (
	"encoding/json"
	"fmt"

	"github.com/Newton-School/goqueue/task"
)

type messageCodec struct{}

func (messageCodec) encode(message task.TaskMessage) ([]byte, error) {
	encoded, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("%w: encode task message: %v", ErrInvalidRedisMessage, err)
	}

	return encoded, nil
}

func (messageCodec) decode(data []byte) (task.TaskMessage, error) {
	var message task.TaskMessage
	if err := json.Unmarshal(data, &message); err != nil {
		return task.TaskMessage{}, fmt.Errorf("%w: decode task message: %v", ErrInvalidRedisMessage, err)
	}
	if message.ID == "" || message.Name == "" || message.Queue == "" {
		return task.TaskMessage{}, fmt.Errorf("%w: decoded message requires id, name, and queue", ErrInvalidRedisMessage)
	}

	return message, nil
}
