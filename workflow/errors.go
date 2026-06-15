package workflow

import "errors"

var (
	// ErrInvalidSignature is returned when a workflow signature is unsafe.
	ErrInvalidSignature = errors.New("goqueue workflow: invalid signature")

	// ErrInvalidWorkflow is returned when a workflow primitive is unsafe.
	ErrInvalidWorkflow = errors.New("goqueue workflow: invalid workflow")

	// ErrNilBackend is returned when a canvas is created without storage.
	ErrNilBackend = errors.New("goqueue workflow: backend is nil")
)
