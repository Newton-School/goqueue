package redisbackend

import (
	"context"
	"fmt"

	"github.com/Newton-School/goqueue/backend"
)

// SaveWorkflowChain stores a chain workflow definition.
func (b *Backend) SaveWorkflowChain(ctx context.Context, record backend.WorkflowChainRecord) error {
	if b.client == nil {
		return fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := record.Validate(); err != nil {
		return err
	}

	return fmt.Errorf("%w: workflow chain save not implemented", ErrInvalidRedisMessage)
}

// AdvanceWorkflowChain records a completed chain step and returns the next signature.
func (b *Backend) AdvanceWorkflowChain(ctx context.Context, request backend.AdvanceWorkflowChainRequest) (backend.AdvanceWorkflowChainResponse, error) {
	if b.client == nil {
		return backend.AdvanceWorkflowChainResponse{}, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return backend.AdvanceWorkflowChainResponse{}, err
	}

	return backend.AdvanceWorkflowChainResponse{}, fmt.Errorf("%w: workflow chain advance not implemented", ErrInvalidRedisMessage)
}

// SaveWorkflowGroup stores a group or chord header workflow definition.
func (b *Backend) SaveWorkflowGroup(ctx context.Context, record backend.WorkflowGroupRecord) error {
	if b.client == nil {
		return fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := record.Validate(); err != nil {
		return err
	}

	return fmt.Errorf("%w: workflow group save not implemented", ErrInvalidRedisMessage)
}

// RecordWorkflowTaskCompleted records terminal progress for a group child.
func (b *Backend) RecordWorkflowTaskCompleted(ctx context.Context, request backend.RecordWorkflowTaskCompletedRequest) (backend.WorkflowGroupProgress, error) {
	if b.client == nil {
		return backend.WorkflowGroupProgress{}, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return backend.WorkflowGroupProgress{}, err
	}

	return backend.WorkflowGroupProgress{}, fmt.Errorf("%w: workflow group progress not implemented", ErrInvalidRedisMessage)
}
