package redisbackend

import (
	"encoding/json"
	"fmt"

	"github.com/Newton-School/goqueue/backend"
)

type resultCodec struct{}

func (resultCodec) encode(record backend.TaskResultRecord) ([]byte, error) {
	encoded, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("%w: encode result: %v", ErrInvalidRedisMessage, err)
	}

	return encoded, nil
}

func (resultCodec) decode(data []byte) (backend.TaskResultRecord, error) {
	var record backend.TaskResultRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return backend.TaskResultRecord{}, fmt.Errorf("%w: decode result: %v", ErrInvalidRedisMessage, err)
	}

	return record, nil
}
