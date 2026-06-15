package redisbackend

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/redis/go-redis/v9"
)

// SaveWorkflowChain stores a chain workflow definition.
func (b *Backend) SaveWorkflowChain(ctx context.Context, record backend.WorkflowChainRecord) error {
	if b.client == nil {
		return fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := record.Validate(); err != nil {
		return err
	}

	codec := workflowSignatureCodec{}
	signatures := make(map[string]any, len(record.Signatures))
	for index, signature := range record.Signatures {
		encoded, err := codec.encode(signature)
		if err != nil {
			return err
		}
		signatures[strconv.Itoa(index)] = string(encoded)
	}

	metaKey := b.keys.workflowChainMeta(record.ID)
	signaturesKey := b.keys.workflowChainSignatures(record.ID)
	ttl := b.options.MessageTTL
	pipe := b.client.TxPipeline()
	pipe.HSet(ctx, metaKey, map[string]any{
		"total":            len(record.Signatures),
		"completed_index":  -1,
		"dispatched_index": 0,
		"created_at":       record.CreatedAt.UTC().Format(time.RFC3339Nano),
	})
	pipe.HSet(ctx, signaturesKey, signatures)
	pipe.Expire(ctx, metaKey, ttl)
	pipe.Expire(ctx, signaturesKey, ttl)

	_, err := pipe.Exec(ctx)
	return err
}

// AdvanceWorkflowChain records a completed chain step and returns the next signature.
func (b *Backend) AdvanceWorkflowChain(ctx context.Context, request backend.AdvanceWorkflowChainRequest) (backend.AdvanceWorkflowChainResponse, error) {
	if b.client == nil {
		return backend.AdvanceWorkflowChainResponse{}, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return backend.AdvanceWorkflowChainResponse{}, err
	}

	values, err := redis.NewScript(advanceWorkflowChainScript()).Run(
		ctx,
		b.client,
		[]string{
			b.keys.workflowChainMeta(request.WorkflowID),
			b.keys.workflowChainSignatures(request.WorkflowID),
		},
		request.CompletedIndex,
		request.CompletedTaskID.String(),
		request.CompletedAt.UTC().Format(time.RFC3339Nano),
	).Slice()
	if err != nil {
		return backend.AdvanceWorkflowChainResponse{}, err
	}

	return parseAdvanceWorkflowChainResponse(values)
}

// SaveWorkflowGroup stores a group or chord header workflow definition.
func (b *Backend) SaveWorkflowGroup(ctx context.Context, record backend.WorkflowGroupRecord) error {
	if b.client == nil {
		return fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := record.Validate(); err != nil {
		return err
	}

	metaKey := b.keys.workflowGroupMeta(record.ID)
	completedKey := b.keys.workflowGroupCompleted(record.ID)
	callbackKey := b.keys.workflowGroupCallback(record.ID)
	ttl := b.options.MessageTTL
	pipe := b.client.TxPipeline()
	pipe.HSet(ctx, metaKey, map[string]any{
		"total":               len(record.TaskIDs),
		"completed":           0,
		"failed":              0,
		"callback_dispatched": 0,
		"created_at":          record.CreatedAt.UTC().Format(time.RFC3339Nano),
	})
	if record.Callback != nil {
		encoded, err := (workflowSignatureCodec{}).encode(*record.Callback)
		if err != nil {
			return err
		}
		pipe.Set(ctx, callbackKey, string(encoded), ttl)
	} else {
		pipe.Del(ctx, callbackKey)
	}
	pipe.Del(ctx, completedKey)
	pipe.Expire(ctx, metaKey, ttl)
	pipe.Expire(ctx, completedKey, ttl)

	_, err := pipe.Exec(ctx)
	return err
}

// RecordWorkflowTaskCompleted records terminal progress for a group child.
func (b *Backend) RecordWorkflowTaskCompleted(ctx context.Context, request backend.RecordWorkflowTaskCompletedRequest) (backend.WorkflowGroupProgress, error) {
	if b.client == nil {
		return backend.WorkflowGroupProgress{}, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return backend.WorkflowGroupProgress{}, err
	}

	values, err := redis.NewScript(recordWorkflowGroupCompletedScript()).Run(
		ctx,
		b.client,
		[]string{
			b.keys.workflowGroupMeta(request.GroupID),
			b.keys.workflowGroupCompleted(request.GroupID),
			b.keys.workflowGroupCallback(request.GroupID),
		},
		request.TaskID.String(),
		request.State.String(),
		request.CompletedAt.UTC().Format(time.RFC3339Nano),
	).Slice()
	if err != nil {
		return backend.WorkflowGroupProgress{}, err
	}

	return parseWorkflowGroupProgress(request.GroupID, values)
}
