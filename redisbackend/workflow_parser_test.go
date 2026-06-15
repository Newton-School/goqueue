package redisbackend

import (
	"errors"
	"testing"
)

func TestParseAdvanceWorkflowChainResponseWithNextSignature(t *testing.T) {
	encoded, err := (workflowSignatureCodec{}).encode(testWorkflowSignatureRecord())
	if err != nil {
		t.Fatalf("encode returned error: %v", err)
	}

	response, err := parseAdvanceWorkflowChainResponse([]any{int64(1), int64(0), string(encoded)})
	if err != nil {
		t.Fatalf("parseAdvanceWorkflowChainResponse returned error: %v", err)
	}

	if !response.Advanced {
		t.Fatal("Advanced should be true")
	}
	if response.Completed {
		t.Fatal("Completed should be false")
	}
	if response.Next == nil {
		t.Fatal("Next should be set")
	}
	if response.Next.Name != "email.send" {
		t.Fatalf("Next name = %q, want email.send", response.Next.Name)
	}
}

func TestParseAdvanceWorkflowChainResponseWithCompletedWorkflow(t *testing.T) {
	response, err := parseAdvanceWorkflowChainResponse([]any{int64(1), int64(1), ""})
	if err != nil {
		t.Fatalf("parseAdvanceWorkflowChainResponse returned error: %v", err)
	}

	if !response.Advanced {
		t.Fatal("Advanced should be true")
	}
	if !response.Completed {
		t.Fatal("Completed should be true")
	}
	if response.Next != nil {
		t.Fatal("Next should be nil")
	}
}

func TestParseAdvanceWorkflowChainResponseRejectsInvalidShape(t *testing.T) {
	_, err := parseAdvanceWorkflowChainResponse([]any{int64(1), ""})
	if !errors.Is(err, ErrInvalidRedisMessage) {
		t.Fatalf("parseAdvanceWorkflowChainResponse error = %v, want ErrInvalidRedisMessage", err)
	}
}

func TestParseWorkflowGroupProgressWithCallback(t *testing.T) {
	encoded, err := (workflowSignatureCodec{}).encode(testWorkflowSignatureRecord())
	if err != nil {
		t.Fatalf("encode returned error: %v", err)
	}

	progress, err := parseWorkflowGroupProgress("group-1", []any{int64(2), int64(2), int64(0), int64(0), int64(1), string(encoded)})
	if err != nil {
		t.Fatalf("parseWorkflowGroupProgress returned error: %v", err)
	}

	if progress.GroupID != "group-1" {
		t.Fatalf("GroupID = %q, want group-1", progress.GroupID)
	}
	if progress.Total != 2 || progress.Completed != 2 || progress.Failed != 0 {
		t.Fatalf("progress = %+v, want total=2 completed=2 failed=0", progress)
	}
	if !progress.Done {
		t.Fatal("Done should be true")
	}
	if !progress.Succeeded {
		t.Fatal("Succeeded should be true")
	}
	if progress.Callback == nil {
		t.Fatal("Callback should be set")
	}
}

func TestParseWorkflowGroupProgressWithoutCallback(t *testing.T) {
	progress, err := parseWorkflowGroupProgress("group-1", []any{int64(2), int64(1), int64(1), int64(0), int64(0), ""})
	if err != nil {
		t.Fatalf("parseWorkflowGroupProgress returned error: %v", err)
	}

	if !progress.Done {
		t.Fatal("Done should be true")
	}
	if progress.Succeeded {
		t.Fatal("Succeeded should be false")
	}
	if progress.Callback != nil {
		t.Fatal("Callback should be nil")
	}
}

func TestParseWorkflowGroupProgressRejectsInvalidShape(t *testing.T) {
	_, err := parseWorkflowGroupProgress("group-1", []any{int64(1), int64(1), ""})
	if !errors.Is(err, ErrInvalidRedisMessage) {
		t.Fatalf("parseWorkflowGroupProgress error = %v, want ErrInvalidRedisMessage", err)
	}
}
