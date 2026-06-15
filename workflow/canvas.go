package workflow

import (
	"context"
	"strconv"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/producer"
	"github.com/Newton-School/goqueue/task"
)

// Canvas dispatches workflow primitives through a queue backend.
type Canvas struct {
	backend      backend.QueueBackend
	producer     *producer.Producer
	defaultQueue task.QueueName
	codec        task.PayloadCodec
	now          func() time.Time
}

// NewCanvas creates a workflow canvas bound to a backend.
func NewCanvas(queueBackend backend.QueueBackend, opts ...CanvasOption) (*Canvas, error) {
	if queueBackend == nil {
		return nil, ErrNilBackend
	}

	config := defaultCanvasConfig()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(&config); err != nil {
			return nil, err
		}
	}

	dispatchProducer, err := producer.NewProducer(
		queueBackend,
		producer.WithProducerDefaultQueue(config.defaultQueue),
		producer.WithProducerCodec(config.codec),
		producer.WithProducerNow(config.now),
	)
	if err != nil {
		return nil, err
	}

	return &Canvas{
		backend:      queueBackend,
		producer:     dispatchProducer,
		defaultQueue: config.defaultQueue,
		codec:        config.codec,
		now:          config.now,
	}, nil
}

// ApplySignature dispatches one signature through the producer path.
func (c *Canvas) ApplySignature(ctx context.Context, signature Signature) (*producer.AsyncResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	normalized, err := signature.Normalize(c.defaultQueue)
	if err != nil {
		return nil, err
	}

	return c.producer.ApplyAsync(
		ctx,
		normalized.Name,
		copyAnySlice(normalized.Args),
		copyAnyMap(normalized.Kwargs),
		producer.WithApplyQueue(normalized.Queue),
		producer.WithApplyMetadata(copyStringMap(normalized.Metadata)),
		producer.WithApplyPriority(normalized.Priority),
		producer.WithApplyRetryPolicy(normalized.RetryPolicy),
		producer.WithApplyETA(normalized.Timing.ETA),
		producer.WithApplyExpiresAt(normalized.Timing.ExpiresAt),
		producer.WithApplyCreatedAt(c.now()),
	)
}

// ApplyChain stores a chain workflow and dispatches the first signature.
func (c *Canvas) ApplyChain(ctx context.Context, chain Chain) (ChainResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	workflowID, err := newWorkflowID()
	if err != nil {
		return ChainResult{}, err
	}
	firstTaskID, err := task.NewTaskID()
	if err != nil {
		return ChainResult{}, err
	}

	record, err := chain.toBackendRecord(workflowID.String(), c.defaultQueue, c.now())
	if err != nil {
		return ChainResult{}, err
	}
	if err := c.backend.SaveWorkflowChain(ctx, record); err != nil {
		return ChainResult{}, err
	}

	first := record.Signatures[0]
	_, err = c.applyRecord(ctx, first, firstTaskID, map[string]string{
		MetadataKindKey:      WorkflowKindChain,
		MetadataChainIDKey:   workflowID.String(),
		MetadataChainStepKey: workflowIndexMetadata(0),
	})
	if err != nil {
		return ChainResult{}, err
	}

	return ChainResult{WorkflowID: workflowID, FirstTask: firstTaskID}, nil
}

func (c *Canvas) applyRecord(
	ctx context.Context,
	record backend.WorkflowSignatureRecord,
	taskID task.TaskID,
	reserved map[string]string,
) (*producer.AsyncResult, error) {
	metadata := MergeMetadata(record.Metadata, reserved)
	return c.producer.ApplyAsync(
		ctx,
		record.Name,
		copyAnySlice(record.Args),
		copyAnyMap(record.Kwargs),
		producer.WithApplyTaskID(taskID),
		producer.WithApplyQueue(record.Queue),
		producer.WithApplyMetadata(metadata),
		producer.WithApplyPriority(record.Priority),
		producer.WithApplyRetryPolicy(record.RetryPolicy),
		producer.WithApplyETA(record.Timing.ETA),
		producer.WithApplyExpiresAt(record.Timing.ExpiresAt),
		producer.WithApplyCreatedAt(c.now()),
	)
}

func workflowIndexMetadata(index int) string {
	return strconv.Itoa(index)
}

func newWorkflowID() (task.TaskID, error) {
	return task.NewTaskID()
}
