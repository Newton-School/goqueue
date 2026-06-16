package inspect

import "context"

// Ping checks Redis/backend connectivity without mutating queue state.
func (i *Inspector) Ping(ctx context.Context) error {
	if i == nil {
		return ErrNilInspector
	}
	if i.backend == nil {
		return ErrInspectorBackend
	}

	return i.backend.Ping(ctx)
}
