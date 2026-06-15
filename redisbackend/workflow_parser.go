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
