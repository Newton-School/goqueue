package redisbackend

import (
	"encoding/json"
	"fmt"

	"github.com/Newton-School/goqueue/backend"
)

type periodicTaskCodec struct{}

func (periodicTaskCodec) encode(record backend.PeriodicTaskRecord) ([]byte, error) {
	if err := record.Validate(); err != nil {
		return nil, err
	}

	encoded, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("%w: encode periodic task: %v", ErrInvalidRedisMessage, err)
	}

	return encoded, nil
}

func (periodicTaskCodec) decode(encoded []byte) (backend.PeriodicTaskRecord, error) {
	var record backend.PeriodicTaskRecord
	if err := json.Unmarshal(encoded, &record); err != nil {
		return backend.PeriodicTaskRecord{}, fmt.Errorf("%w: decode periodic task: %v", ErrInvalidRedisMessage, err)
	}
	if err := record.Validate(); err != nil {
		return backend.PeriodicTaskRecord{}, err
	}

	return record, nil
}
