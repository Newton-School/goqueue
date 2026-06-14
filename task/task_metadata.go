package task

// TaskMetadata stores task-level headers and tracing fields.
type TaskMetadata struct {
	values map[string]string
}

// NewTaskMetadata copies values into task metadata.
func NewTaskMetadata(values map[string]string) TaskMetadata {
	return TaskMetadata{values: cloneStringMap(values)}
}

// Values returns a copy of the metadata map.
func (m TaskMetadata) Values() map[string]string {
	return cloneStringMap(m.values)
}

func cloneStringMap(values map[string]string) map[string]string {
	if len(values) == 0 {
		return nil
	}

	cloned := make(map[string]string, len(values))
	for key, value := range values {
		cloned[key] = value
	}

	return cloned
}
