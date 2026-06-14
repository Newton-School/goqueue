package goqueue

// TaskPayload contains task call arguments.
type TaskPayload struct {
	args   []any
	kwargs map[string]any
}

// NewTaskPayload copies args and kwargs into a new task payload.
func NewTaskPayload(args []any, kwargs map[string]any) TaskPayload {
	return TaskPayload{
		args:   cloneAnySlice(args),
		kwargs: cloneAnyMap(kwargs),
	}
}

// Args returns a copy of positional arguments.
func (p TaskPayload) Args() []any {
	return cloneAnySlice(p.args)
}

// Kwargs returns a copy of keyword arguments.
func (p TaskPayload) Kwargs() map[string]any {
	return cloneAnyMap(p.kwargs)
}

func cloneAnySlice(values []any) []any {
	if len(values) == 0 {
		return nil
	}

	cloned := make([]any, len(values))
	copy(cloned, values)
	return cloned
}

func cloneAnyMap(values map[string]any) map[string]any {
	if len(values) == 0 {
		return nil
	}

	cloned := make(map[string]any, len(values))
	for key, value := range values {
		cloned[key] = value
	}

	return cloned
}
