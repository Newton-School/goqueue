package goqueue

import "testing"

func TestNewTaskMetadataCopiesInput(t *testing.T) {
	values := map[string]string{"trace_id": "trace-1"}

	metadata := NewTaskMetadata(values)
	values["trace_id"] = "trace-2"

	if got := metadata.Values()["trace_id"]; got != "trace-1" {
		t.Fatalf("Values()[trace_id] = %q, want original value", got)
	}
}

func TestTaskMetadataValuesReturnsCopy(t *testing.T) {
	metadata := NewTaskMetadata(map[string]string{"trace_id": "trace-1"})

	values := metadata.Values()
	values["trace_id"] = "trace-2"

	if got := metadata.Values()["trace_id"]; got != "trace-1" {
		t.Fatalf("Values()[trace_id] = %q, want original value", got)
	}
}
