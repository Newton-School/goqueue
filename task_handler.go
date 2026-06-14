package goqueue

// TaskHandler executes a registered task.
type TaskHandler interface {
	HandleTask(HandlerContext, TaskPayload) (TaskResult, error)
}

// TaskHandlerFunc adapts a function into a TaskHandler.
type TaskHandlerFunc func(HandlerContext, TaskPayload) (TaskResult, error)

// HandleTask executes f.
func (f TaskHandlerFunc) HandleTask(ctx HandlerContext, payload TaskPayload) (TaskResult, error) {
	return f(ctx, payload)
}
