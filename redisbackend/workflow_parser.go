package redisbackend

import (
	"fmt"

	"github.com/Newton-School/goqueue/backend"
)

func parseAdvanceWorkflowChainResponse(values []any) (backend.AdvanceWorkflowChainResponse, error) {
	if len(values) != 3 {
		return backend.AdvanceWorkflowChainResponse{}, fmt.Errorf("%w: workflow chain advance response must have 3 values", ErrInvalidRedisMessage)
	}

	advanced := redisInt(values[0]) == 1
	completed := redisInt(values[1]) == 1
	response := backend.AdvanceWorkflowChainResponse{
		Advanced:  advanced,
		Completed: completed,
	}

	rawNext, ok := values[2].(string)
	if !ok || rawNext == "" {
		return response, nil
	}

	next, err := (workflowSignatureCodec{}).decode([]byte(rawNext))
	if err != nil {
		return backend.AdvanceWorkflowChainResponse{}, err
	}
	response.Next = &next
	return response, nil
}

func redisInt(value any) int64 {
	switch typed := value.(type) {
	case int64:
		return typed
	case int:
		return int64(typed)
	default:
		return 0
	}
}

func parseWorkflowGroupProgress(groupID string, values []any) (backend.WorkflowGroupProgress, error) {
	if len(values) != 6 {
		return backend.WorkflowGroupProgress{}, fmt.Errorf("%w: workflow group progress response must have 6 values", ErrInvalidRedisMessage)
	}

	total := int(redisInt(values[0]))
	completed := int(redisInt(values[1]))
	failed := int(redisInt(values[2]))
	progress := backend.WorkflowGroupProgress{
		GroupID:   groupID,
		Total:     total,
		Completed: completed,
		Failed:    failed,
		Duplicate: redisInt(values[3]) == 1,
		Done:      completed+failed >= total && total > 0,
		Succeeded: redisInt(values[4]) == 1,
	}

	rawCallback, ok := values[5].(string)
	if !ok || rawCallback == "" {
		return progress, nil
	}

	callback, err := (workflowSignatureCodec{}).decode([]byte(rawCallback))
	if err != nil {
		return backend.WorkflowGroupProgress{}, err
	}
	progress.Callback = &callback
	return progress, nil
}
