package redisbackend

import (
	"encoding/json"
	"fmt"

	"github.com/Newton-School/goqueue/backend"
)

type stateCodec struct{}

func (stateCodec) encode(record backend.TaskStateRecord) ([]byte, error) {
	encoded, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("%w: encode state: %v", ErrInvalidRedisMessage, err)
	}

	return encoded, nil
}

func (stateCodec) decode(data []byte) (backend.TaskStateRecord, error) {
	var record backend.TaskStateRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return backend.TaskStateRecord{}, fmt.Errorf("%w: decode state: %v", ErrInvalidRedisMessage, err)
	}

	return record, nil
}
