package task

// TaskResult is returned by task handlers.
type TaskResult struct {
	State    TaskState
	Value    any
	Error    string
	Metadata map[string]string
}

// SucceededResult returns a successful task result.
func SucceededResult(value any) TaskResult {
	return TaskResult{
		State: TaskSucceeded,
		Value: value,
	}
}

// FailedResult returns a failed task result.
func FailedResult(err error) TaskResult {
	message := ""
	if err != nil {
		message = err.Error()
	}

	return TaskResult{
		State: TaskFailed,
		Error: message,
	}
}

// Validate verifies result state.
func (r TaskResult) Validate() error {
	return ValidateTaskState(r.State)
}
