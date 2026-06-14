package redisbackend

import (
	"fmt"

	"github.com/Newton-School/goqueue/backend"
	"github.com/redis/go-redis/v9"
)

func parseReadyStreamMessages(streams []redis.XStream) ([]backend.ReadyMessage, error) {
	codec := messageCodec{}
	ready := make([]backend.ReadyMessage, 0)

	for _, stream := range streams {
		for _, message := range stream.Messages {
			raw, exists := message.Values["message"]
			if !exists {
				return nil, fmt.Errorf("%w: stream entry missing message field", ErrInvalidRedisMessage)
			}
			encoded, ok := raw.(string)
			if !ok {
				return nil, fmt.Errorf("%w: stream message field must be string", ErrInvalidRedisMessage)
			}

			decoded, err := codec.decode([]byte(encoded))
			if err != nil {
				return nil, err
			}
			ready = append(ready, backend.ReadyMessage{
				StreamID: message.ID,
				Message:  decoded,
			})
		}
	}

	return ready, nil
}
