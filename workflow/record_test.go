package workflow

import (
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
)

func TestSignatureToBackendRecordPreservesFields(t *testing.T) {
	signature := validSignature()

	record, err := signature.toBackendRecord("critical")
	if err != nil {
		t.Fatalf("toBackendRecord returned error: %v", err)
	}

	if record.Name != signature.Name {
		t.Fatalf("Name = %q, want %q", record.Name, signature.Name)
	}
	if record.Queue != signature.Queue {
		t.Fatalf("Queue = %q, want %q", record.Queue, signature.Queue)
	}
	if record.Metadata["source"] != "workflow" {
		t.Fatalf("Metadata source = %q, want workflow", record.Metadata["source"])
	}
}

func TestBackendRecordToSignatureRestoresSignature(t *testing.T) {
	record := backend.WorkflowSignatureRecord{
		Name:        "email.send",
		Queue:       "critical",
		Args:        []any{"u_123"},
		Kwargs:      map[string]any{"template": "welcome"},
		Metadata:    map[string]string{"source": "workflow"},
		Timing:      validSignature().Timing,
		Priority:    validSignature().Priority,
		RetryPolicy: validSignature().RetryPolicy,
	}

	signature, err := signatureFromBackendRecord(record)
	if err != nil {
		t.Fatalf("signatureFromBackendRecord returned error: %v", err)
	}

	if signature.Name != "email.send" {
		t.Fatalf("Name = %q, want email.send", signature.Name)
	}
	if signature.Queue != "critical" {
		t.Fatalf("Queue = %q, want critical", signature.Queue)
	}
	if signature.Kwargs["template"] != "welcome" {
		t.Fatalf("Kwargs template = %v, want welcome", signature.Kwargs["template"])
	}
}

func TestChainToBackendRecordUsesProvidedID(t *testing.T) {
	chain := Chain{Signatures: []Signature{validSignature()}}
	now := time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC)

	record, err := chain.toBackendRecord("chain-1", "default", now)
	if err != nil {
		t.Fatalf("toBackendRecord returned error: %v", err)
	}

	if record.ID != "chain-1" {
		t.Fatalf("ID = %q, want chain-1", record.ID)
	}
	if len(record.Signatures) != 1 {
		t.Fatalf("signature count = %d, want 1", len(record.Signatures))
	}
	if !record.CreatedAt.Equal(now) {
		t.Fatalf("CreatedAt = %v, want %v", record.CreatedAt, now)
	}
}
