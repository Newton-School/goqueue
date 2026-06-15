package workflow

import (
	"errors"
	"testing"
)

func TestChordValidateRequiresHeaderGroup(t *testing.T) {
	chord := Chord{Callback: validSignature()}

	if err := chord.Validate(); !errors.Is(err, ErrInvalidWorkflow) {
		t.Fatalf("Validate error = %v, want ErrInvalidWorkflow", err)
	}
}

func TestChordValidateRequiresCallback(t *testing.T) {
	chord := Chord{Header: Group{Signatures: []Signature{validSignature()}}}

	if err := chord.Validate(); !errors.Is(err, ErrInvalidWorkflow) {
		t.Fatalf("Validate error = %v, want ErrInvalidWorkflow", err)
	}
}

func TestChordNormalizeAppliesDefaults(t *testing.T) {
	header := validSignature()
	header.Queue = ""
	callback := validSignature()
	callback.Name = "email.complete"
	callback.Queue = ""

	chord := Chord{
		Header:   Group{Signatures: []Signature{header}},
		Callback: callback,
	}

	normalized, err := chord.Normalize("critical")
	if err != nil {
		t.Fatalf("Normalize returned error: %v", err)
	}

	if normalized.Header.Signatures[0].Queue != "critical" {
		t.Fatalf("header queue = %q, want critical", normalized.Header.Signatures[0].Queue)
	}
	if normalized.Callback.Queue != "critical" {
		t.Fatalf("callback queue = %q, want critical", normalized.Callback.Queue)
	}
}
