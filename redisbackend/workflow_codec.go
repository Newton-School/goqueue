package redisbackend

import (
	"encoding/json"
	"fmt"

	"github.com/Newton-School/goqueue/backend"
)

type workflowSignatureCodec struct{}

func (workflowSignatureCodec) encode(record backend.WorkflowSignatureRecord) ([]byte, error) {
	if err := record.Validate(); err != nil {
		return nil, err
	}

	encoded, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("%w: encode workflow signature: %v", ErrInvalidRedisMessage, err)
	}

	return encoded, nil
}

func (workflowSignatureCodec) decode(encoded []byte) (backend.WorkflowSignatureRecord, error) {
	var record backend.WorkflowSignatureRecord
	if err := json.Unmarshal(encoded, &record); err != nil {
		return backend.WorkflowSignatureRecord{}, fmt.Errorf("%w: decode workflow signature: %v", ErrInvalidRedisMessage, err)
	}
	if err := record.Validate(); err != nil {
		return backend.WorkflowSignatureRecord{}, err
	}

	return record, nil
}
