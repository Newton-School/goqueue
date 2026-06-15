package redisbackend

import (
	"context"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

func TestWorkflowChainRedisLifecycle(t *testing.T) {
	ctx := context.Background()
	options := redisIntegrationOptions(t)
	b, err := New(options)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	defer b.Close()
	cleanupIntegrationNamespace(ctx, t, b)

	record := testWorkflowChainRecord()
	if err := b.SaveWorkflowChain(ctx, record); err != nil {
		t.Fatalf("SaveWorkflowChain returned error: %v", err)
	}

	response, err := b.AdvanceWorkflowChain(ctx, backend.AdvanceWorkflowChainRequest{
		WorkflowID:      record.ID,
		CompletedTaskID: "11111111-1111-4111-8111-111111111111",
		CompletedIndex:  0,
		CompletedAt:     time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("AdvanceWorkflowChain returned error: %v", err)
	}
	if !response.Advanced || response.Next == nil {
		t.Fatalf("response = %+v, want advanced with next signature", response)
	}

	duplicate, err := b.AdvanceWorkflowChain(ctx, backend.AdvanceWorkflowChainRequest{
		WorkflowID:      record.ID,
		CompletedTaskID: "11111111-1111-4111-8111-111111111111",
		CompletedIndex:  0,
		CompletedAt:     time.Date(2026, time.June, 15, 12, 1, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("duplicate AdvanceWorkflowChain returned error: %v", err)
	}
	if duplicate.Advanced || duplicate.Next != nil {
		t.Fatalf("duplicate response = %+v, want no advancement", duplicate)
	}
}

func TestWorkflowGroupRedisLifecycle(t *testing.T) {
	ctx := context.Background()
	options := redisIntegrationOptions(t)
	b, err := New(options)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	defer b.Close()
	cleanupIntegrationNamespace(ctx, t, b)

	callback := testWorkflowSignatureRecord()
	record := testWorkflowGroupRecord(&callback)
	if err := b.SaveWorkflowGroup(ctx, record); err != nil {
		t.Fatalf("SaveWorkflowGroup returned error: %v", err)
	}

	first, err := b.RecordWorkflowTaskCompleted(ctx, backend.RecordWorkflowTaskCompletedRequest{
		GroupID:     record.ID,
		TaskID:      record.TaskIDs[0],
		State:       task.TaskSucceeded,
		CompletedAt: time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("RecordWorkflowTaskCompleted first returned error: %v", err)
	}
	if first.Done || first.Callback != nil {
		t.Fatalf("first progress = %+v, want not done", first)
	}

	second, err := b.RecordWorkflowTaskCompleted(ctx, backend.RecordWorkflowTaskCompletedRequest{
		GroupID:     record.ID,
		TaskID:      record.TaskIDs[1],
		State:       task.TaskSucceeded,
		CompletedAt: time.Date(2026, time.June, 15, 12, 1, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("RecordWorkflowTaskCompleted second returned error: %v", err)
	}
	if !second.Done || !second.Succeeded || second.Callback == nil {
		t.Fatalf("second progress = %+v, want successful done with callback", second)
	}
}
