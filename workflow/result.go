package workflow

import "github.com/Newton-School/goqueue/task"

// ChainResult describes a dispatched chain workflow.
type ChainResult struct {
	WorkflowID task.TaskID
	FirstTask  task.TaskID
}

// GroupResult describes a dispatched group workflow.
type GroupResult struct {
	GroupID task.TaskID
	TaskIDs []task.TaskID
}

// ChordResult describes a dispatched chord workflow.
type ChordResult struct {
	GroupID task.TaskID
	TaskIDs []task.TaskID
}
