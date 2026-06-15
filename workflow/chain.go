package workflow

import (
	"fmt"

	"github.com/Newton-School/goqueue/task"
)

// Chain dispatches signatures sequentially.
type Chain struct {
	Signatures []Signature
}

// Validate verifies that the chain has at least one valid signature.
func (c Chain) Validate() error {
	if len(c.Signatures) == 0 {
		return fmt.Errorf("%w: chain requires at least one signature", ErrInvalidWorkflow)
	}
	for index, signature := range c.Signatures {
		if err := signature.Validate(); err != nil {
			return fmt.Errorf("%w: chain signature %d: %v", ErrInvalidWorkflow, index, err)
		}
	}

	return nil
}

// Normalize applies defaults to every chain signature and copies the slice.
func (c Chain) Normalize(defaultQueue task.QueueName) (Chain, error) {
	normalized := Chain{Signatures: make([]Signature, len(c.Signatures))}
	for index, signature := range c.Signatures {
		normalizedSignature, err := signature.Normalize(defaultQueue)
		if err != nil {
			return Chain{}, fmt.Errorf("%w: chain signature %d: %v", ErrInvalidWorkflow, index, err)
		}
		normalized.Signatures[index] = normalizedSignature
	}

	if err := normalized.Validate(); err != nil {
		return Chain{}, err
	}

	return normalized, nil
}
