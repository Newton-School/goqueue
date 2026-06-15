package redisbackend

import (
	"errors"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

func TestWorkflowSignatureCodecRoundTripsRecord(t *testing.T) {
	record := testWorkflowSignatureRecord()

	encoded, err := (workflowSignatureCodec{}).encode(record)
	if err != nil {
		t.Fatalf("encode returned error: %v", err)
	}

	decoded, err := (workflowSignatureCodec{}).decode(encoded)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if decoded.Name != record.Name {
		t.Fatalf("Name = %q, want %q", decoded.Name, record.Name)
	}
	if decoded.Queue != record.Queue {
		t.Fatalf("Queue = %q, want %q", decoded.Queue, record.Queue)
	}
	if decoded.Metadata["source"] != "workflow" {
		t.Fatalf("Metadata source = %q, want workflow", decoded.Metadata["source"])
	}
}

func TestWorkflowSignatureCodecRejectsInvalidJSON(t *testing.T) {
	_, err := (workflowSignatureCodec{}).decode([]byte("{"))
	if !errors.Is(err, ErrInvalidRedisMessage) {
		t.Fatalf("decode error = %v, want ErrInvalidRedisMessage", err)
	}
}

func testWorkflowSignatureRecord() backend.WorkflowSignatureRecord {
	return backend.WorkflowSignatureRecord{
		Name:        "email.send",
		Queue:       "default",
		Args:        []any{"u_123"},
		Kwargs:      map[string]any{"template": "welcome"},
		Metadata:    map[string]string{"source": "workflow"},
		Timing:      task.TaskTiming{ExpiresAt: time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC)},
		Priority:    task.DefaultPriority,
		RetryPolicy: task.DefaultRetryPolicy(),
	}
}
