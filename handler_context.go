package goqueue

import "context"

// HandlerContext carries task execution context for a handler invocation.
type HandlerContext struct {
	ctx      context.Context
	envelope TaskEnvelope
}

// NewHandlerContext creates a handler context with a cloned task envelope.
func NewHandlerContext(ctx context.Context, envelope TaskEnvelope) HandlerContext {
	if ctx == nil {
		ctx = context.Background()
	}

	return HandlerContext{
		ctx:      ctx,
		envelope: envelope.Clone(),
	}
}

// Context returns the underlying cancellation/deadline context.
func (c HandlerContext) Context() context.Context {
	return c.ctx
}

// Envelope returns a copy of the task envelope.
func (c HandlerContext) Envelope() TaskEnvelope {
	return c.envelope.Clone()
}

// TaskID returns the current task ID.
func (c HandlerContext) TaskID() TaskID {
	return c.envelope.ID
}

// TaskName returns the current task name.
func (c HandlerContext) TaskName() TaskName {
	return c.envelope.Name
}

// Queue returns the queue from which the task was read.
func (c HandlerContext) Queue() QueueName {
	return c.envelope.Queue
}

// Attempt returns the current task attempt count.
func (c HandlerContext) Attempt() int {
	return c.envelope.Attempt
}

// Metadata returns a copy of task metadata.
func (c HandlerContext) Metadata() map[string]string {
	return c.envelope.Metadata.Values()
}
