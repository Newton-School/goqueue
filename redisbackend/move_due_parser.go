package redisbackend

import (
	"fmt"

	"github.com/Newton-School/goqueue/backend"
)

func parseMovedScheduledMessages(values []any) ([]backend.MovedScheduledMessage, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("%w: moved response must contain stream/message pairs", ErrInvalidRedisMessage)
	}

	moved := make([]backend.MovedScheduledMessage, 0, len(values)/2)
	codec := messageCodec{}
	for i := 0; i < len(values); i += 2 {
		streamID, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("%w: moved stream id must be string", ErrInvalidRedisMessage)
		}

		encoded, ok := values[i+1].(string)
		if !ok {
			return nil, fmt.Errorf("%w: moved message must be string", ErrInvalidRedisMessage)
		}

		message, err := codec.decode([]byte(encoded))
		if err != nil {
			return nil, err
		}
		moved = append(moved, backend.MovedScheduledMessage{
			StreamID: streamID,
			Message:  message,
		})
	}

	return moved, nil
}
