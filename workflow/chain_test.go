package workflow

import (
	"errors"
	"testing"
)

func TestChainValidateRequiresAtLeastOneSignature(t *testing.T) {
	chain := Chain{}

	if err := chain.Validate(); !errors.Is(err, ErrInvalidWorkflow) {
		t.Fatalf("Validate error = %v, want ErrInvalidWorkflow", err)
	}
}

func TestChainValidateRejectsInvalidSignature(t *testing.T) {
	chain := Chain{Signatures: []Signature{{}}}

	if err := chain.Validate(); err == nil {
		t.Fatal("Validate expected error")
	}
}

func TestChainNormalizeAppliesSignatureDefaults(t *testing.T) {
	first := validSignature()
	first.Queue = ""
	second := validSignature()
	second.Name = "email.audit"
	second.Queue = ""

	chain := Chain{Signatures: []Signature{first, second}}
	normalized, err := chain.Normalize("critical")
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

func TestChainNormalizeCopiesSignatures(t *testing.T) {
	signature := validSignature()
	chain := Chain{Signatures: []Signature{signature}}

	normalized, err := chain.Normalize("default")
	if err != nil {
		t.Fatalf("Normalize returned error: %v", err)
	}

	chain.Signatures[0].Args[0] = "changed"
	if normalized.Signatures[0].Args[0] != "u_123" {
		t.Fatalf("normalized arg = %v, want copied value", normalized.Signatures[0].Args[0])
	}
}
