package redisbackend

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

func TestSaveWorkflowChainRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	err := b.SaveWorkflowChain(context.Background(), testWorkflowChainRecord())
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("SaveWorkflowChain error = %v, want ErrInvalidRedisOptions", err)
	}
}

func TestAdvanceWorkflowChainRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	_, err := b.AdvanceWorkflowChain(context.Background(), backend.AdvanceWorkflowChainRequest{
		WorkflowID:      "chain-1",
		CompletedTaskID: "11111111-1111-4111-8111-111111111111",
		CompletedIndex:  0,
		CompletedAt:     time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("AdvanceWorkflowChain error = %v, want ErrInvalidRedisOptions", err)
	}
}

func TestSaveWorkflowGroupRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	err := b.SaveWorkflowGroup(context.Background(), testWorkflowGroupRecord(nil))
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("SaveWorkflowGroup error = %v, want ErrInvalidRedisOptions", err)
	}
}

func TestRecordWorkflowTaskCompletedRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	_, err := b.RecordWorkflowTaskCompleted(context.Background(), backend.RecordWorkflowTaskCompletedRequest{
		GroupID:     "group-1",
		TaskID:      "11111111-1111-4111-8111-111111111111",
		State:       task.TaskSucceeded,
		CompletedAt: time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("RecordWorkflowTaskCompleted error = %v, want ErrInvalidRedisOptions", err)
	}
}

func testWorkflowChainRecord() backend.WorkflowChainRecord {
	return backend.WorkflowChainRecord{
		ID:         "chain-1",
		Signatures: []backend.WorkflowSignatureRecord{testWorkflowSignatureRecord(), testWorkflowSignatureRecord()},
		CreatedAt:  time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC),
	}
}

func testWorkflowGroupRecord(callback *backend.WorkflowSignatureRecord) backend.WorkflowGroupRecord {
	return backend.WorkflowGroupRecord{
		ID: "group-1",
		TaskIDs: []task.TaskID{
			"11111111-1111-4111-8111-111111111111",
			"22222222-2222-4222-8222-222222222222",
		},
		Callback:  callback,
		CreatedAt: time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC),
	}
}
