package workflow

import (
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

func (s Signature) toBackendRecord(defaultQueue task.QueueName) (backend.WorkflowSignatureRecord, error) {
	normalized, err := s.Normalize(defaultQueue)
	if err != nil {
		return backend.WorkflowSignatureRecord{}, err
	}

	return backend.WorkflowSignatureRecord{
		Name:        normalized.Name,
		Queue:       normalized.Queue,
		Args:        copyAnySlice(normalized.Args),
		Kwargs:      copyAnyMap(normalized.Kwargs),
		Metadata:    copyStringMap(normalized.Metadata),
		Timing:      normalized.Timing,
		Priority:    normalized.Priority,
		RetryPolicy: normalized.RetryPolicy,
	}, nil
}

func signatureFromBackendRecord(record backend.WorkflowSignatureRecord) (Signature, error) {
	if err := record.Validate(); err != nil {
		return Signature{}, err
	}

	signature := Signature{
		Name:        record.Name,
		Queue:       record.Queue,
		Args:        copyAnySlice(record.Args),
		Kwargs:      copyAnyMap(record.Kwargs),
		Metadata:    copyStringMap(record.Metadata),
		Timing:      record.Timing,
		Priority:    record.Priority,
		RetryPolicy: record.RetryPolicy,
	}
	if err := signature.Validate(); err != nil {
		return Signature{}, err
	}

	return signature, nil
}

func (c Chain) toBackendRecord(id string, defaultQueue task.QueueName, now time.Time) (backend.WorkflowChainRecord, error) {
	normalized, err := c.Normalize(defaultQueue)
	if err != nil {
		return backend.WorkflowChainRecord{}, err
	}

	record := backend.WorkflowChainRecord{
		ID:         id,
		Signatures: make([]backend.WorkflowSignatureRecord, len(normalized.Signatures)),
		CreatedAt:  now.UTC(),
	}
	for index, signature := range normalized.Signatures {
		signatureRecord, err := signature.toBackendRecord(defaultQueue)
		if err != nil {
			return backend.WorkflowChainRecord{}, err
		}
		record.Signatures[index] = signatureRecord
	}
	if err := record.Validate(); err != nil {
		return backend.WorkflowChainRecord{}, err
	}

	return record, nil
}
