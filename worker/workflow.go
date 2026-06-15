package worker

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
	"github.com/Newton-School/goqueue/workflow"
)

func (w *Worker) advanceWorkflow(ctx context.Context, envelope task.TaskEnvelope, finalState task.TaskState) error {
	if finalState != task.TaskSucceeded {
		return w.recordGroupProgress(ctx, envelope, finalState)
	}

	if err := w.advanceChain(ctx, envelope); err != nil {
		return err
	}
	return w.recordGroupProgress(ctx, envelope, finalState)
}

func (w *Worker) advanceChain(ctx context.Context, envelope task.TaskEnvelope) error {
	metadata := envelope.Metadata.Values()
	if metadata[workflow.MetadataKindKey] != workflow.WorkflowKindChain {
		return nil
	}

	chainID := metadata[workflow.MetadataChainIDKey]
	stepValue := metadata[workflow.MetadataChainStepKey]
	if chainID == "" || stepValue == "" {
		return nil
	}
	step, err := strconv.Atoi(stepValue)
	if err != nil {
		return fmt.Errorf("goqueue worker: invalid chain step %q: %w", stepValue, err)
	}

	response, err := w.backend.AdvanceWorkflowChain(ctx, backend.AdvanceWorkflowChainRequest{
		WorkflowID:      chainID,
		CompletedTaskID: envelope.ID,
		CompletedIndex:  step,
		CompletedAt:     w.now(),
	})
	if err != nil {
		return fmt.Errorf("goqueue worker: advance chain workflow: %w", err)
	}
	if response.Next == nil {
		return nil
	}

	return w.enqueueWorkflowSignature(ctx, *response.Next, map[string]string{
		workflow.MetadataKindKey:      workflow.WorkflowKindChain,
		workflow.MetadataChainIDKey:   chainID,
		workflow.MetadataChainStepKey: strconv.Itoa(step + 1),
	})
}

func (w *Worker) recordGroupProgress(ctx context.Context, envelope task.TaskEnvelope, finalState task.TaskState) error {
	metadata := envelope.Metadata.Values()
	kind := metadata[workflow.MetadataKindKey]
	if kind != workflow.WorkflowKindGroup && kind != workflow.WorkflowKindChord {
		return nil
	}

	groupID := metadata[workflow.MetadataGroupIDKey]
	if groupID == "" {
		groupID = metadata[workflow.MetadataChordIDKey]
	}
	if groupID == "" {
		return nil
	}

	progress, err := w.backend.RecordWorkflowTaskCompleted(ctx, backend.RecordWorkflowTaskCompletedRequest{
		GroupID:     groupID,
		TaskID:      envelope.ID,
		State:       finalState,
		CompletedAt: w.now(),
	})
	if err != nil {
		return fmt.Errorf("goqueue worker: record group workflow progress: %w", err)
	}
	if progress.Callback == nil {
		return nil
	}

	return w.enqueueWorkflowSignature(ctx, *progress.Callback, map[string]string{
		workflow.MetadataKindKey:          workflow.WorkflowKindChord,
		workflow.MetadataGroupIDKey:       groupID,
		workflow.MetadataChordIDKey:       groupID,
		workflow.MetadataChordCallbackKey: "true",
	})
}

func (w *Worker) enqueueWorkflowSignature(
	ctx context.Context,
	record backend.WorkflowSignatureRecord,
	reserved map[string]string,
) error {
	id, err := task.NewTaskID()
	if err != nil {
		return fmt.Errorf("goqueue worker: generate workflow task id: %w", err)
	}

	metadata := workflow.MergeMetadata(record.Metadata, reserved)
	envelope, err := task.NewTaskEnvelope(task.TaskEnvelopeInput{
		ID:          id,
		Name:        record.Name,
		Queue:       record.Queue,
		Args:        record.Args,
		Kwargs:      record.Kwargs,
		Metadata:    metadata,
		Timing:      record.Timing,
		Priority:    record.Priority,
		RetryPolicy: record.RetryPolicy,
		CreatedAt:   w.now(),
	})
	if err != nil {
		return fmt.Errorf("goqueue worker: build workflow envelope: %w", err)
	}

	message, err := task.TaskEnvelopeToMessage(envelope, w.codec)
	if err != nil {
		return fmt.Errorf("goqueue worker: encode workflow message: %w", err)
	}

	initialState := task.TaskPending
	if envelope.Timing.Scheduled() {
		initialState = task.TaskScheduled
	}
	if err := w.writeState(ctx, envelope.ID, initialState, ""); err != nil {
		return err
	}
	if envelope.Timing.Scheduled() {
		if _, err := w.backend.EnqueueScheduled(ctx, backend.EnqueueRequest{Message: message}); err != nil {
			return fmt.Errorf("goqueue worker: enqueue scheduled workflow task: %w", err)
		}
		return nil
	}
	if _, err := w.backend.EnqueueReady(ctx, backend.EnqueueRequest{Message: message}); err != nil {
		return fmt.Errorf("goqueue worker: enqueue ready workflow task: %w", err)
	}
	return nil
}
