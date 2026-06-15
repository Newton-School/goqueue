package workflow

import (
	"errors"
	"testing"
)

func TestGroupValidateRequiresAtLeastOneSignature(t *testing.T) {
	group := Group{}

	if err := group.Validate(); !errors.Is(err, ErrInvalidWorkflow) {
		t.Fatalf("Validate error = %v, want ErrInvalidWorkflow", err)
	}
}

func TestGroupValidateRejectsInvalidSignature(t *testing.T) {
	group := Group{Signatures: []Signature{{}}}

	if err := group.Validate(); err == nil {
		t.Fatal("Validate expected error")
	}
}

func TestGroupNormalizeAppliesSignatureDefaults(t *testing.T) {
	first := validSignature()
	first.Queue = ""
	second := validSignature()
	second.Name = "email.audit"
	second.Queue = ""

	group := Group{Signatures: []Signature{first, second}}
	normalized, err := group.Normalize("critical")
	if err != nil {
		t.Fatalf("Normalize returned error: %v", err)
	}

	if len(normalized.Signatures) != 2 {
		t.Fatalf("signature count = %d, want 2", len(normalized.Signatures))
	}
	if normalized.Signatures[0].Queue != "critical" {
		t.Fatalf("first queue = %q, want critical", normalized.Signatures[0].Queue)
	}
	if normalized.Signatures[1].Queue != "critical" {
		t.Fatalf("second queue = %q, want critical", normalized.Signatures[1].Queue)
	}
}
