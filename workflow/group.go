package workflow

import (
	"fmt"

	"github.com/Newton-School/goqueue/task"
)

// Group dispatches signatures in parallel and tracks their combined progress.
type Group struct {
	Signatures []Signature
}

// Validate verifies that the group has at least one valid signature.
func (g Group) Validate() error {
	if len(g.Signatures) == 0 {
		return fmt.Errorf("%w: group requires at least one signature", ErrInvalidWorkflow)
	}
	for index, signature := range g.Signatures {
		if err := signature.Validate(); err != nil {
			return fmt.Errorf("%w: group signature %d: %v", ErrInvalidWorkflow, index, err)
		}
	}

	return nil
}

// Normalize applies defaults to every group signature and copies the slice.
func (g Group) Normalize(defaultQueue task.QueueName) (Group, error) {
	normalized := Group{Signatures: make([]Signature, len(g.Signatures))}
	for index, signature := range g.Signatures {
		normalizedSignature, err := signature.Normalize(defaultQueue)
		if err != nil {
			return Group{}, fmt.Errorf("%w: group signature %d: %v", ErrInvalidWorkflow, index, err)
		}
		normalized.Signatures[index] = normalizedSignature
	}

	if err := normalized.Validate(); err != nil {
		return Group{}, err
	}

	return normalized, nil
}
