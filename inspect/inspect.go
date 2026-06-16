package inspect

import (
	"fmt"

	"github.com/Newton-School/goqueue/backend"
)

// Inspector provides read-only visibility APIs over task lifecycle data.
type Inspector struct {
	backend backend.QueueBackend
}

// NewInspector creates an inspector bound to a queue backend.
func NewInspector(queueBackend backend.QueueBackend) (*Inspector, error) {
	if queueBackend == nil {
		return nil, fmt.Errorf("inspect: queue backend is required")
	}

	return &Inspector{backend: queueBackend}, nil
}
