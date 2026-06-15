package workflow

import (
	"fmt"

	"github.com/Newton-School/goqueue/task"
)

// Chord dispatches a group and then a callback after the group succeeds.
type Chord struct {
	Header   Group
	Callback Signature
}

// Validate verifies that the chord header and callback are safe.
func (c Chord) Validate() error {
	if err := c.Header.Validate(); err != nil {
		return fmt.Errorf("%w: chord header: %v", ErrInvalidWorkflow, err)
	}
	if err := c.Callback.Validate(); err != nil {
		return fmt.Errorf("%w: chord callback: %v", ErrInvalidWorkflow, err)
	}

	return nil
}

// Normalize applies defaults to the header group and callback.
func (c Chord) Normalize(defaultQueue task.QueueName) (Chord, error) {
	header, err := c.Header.Normalize(defaultQueue)
	if err != nil {
		return Chord{}, fmt.Errorf("%w: chord header: %v", ErrInvalidWorkflow, err)
	}
	callback, err := c.Callback.Normalize(defaultQueue)
	if err != nil {
		return Chord{}, fmt.Errorf("%w: chord callback: %v", ErrInvalidWorkflow, err)
	}

	normalized := Chord{
		Header:   header,
		Callback: callback,
	}
	if err := normalized.Validate(); err != nil {
		return Chord{}, err
	}

	return normalized, nil
}
