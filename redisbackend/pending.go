package redisbackend

import (
	"context"
	"fmt"

	"github.com/Newton-School/goqueue/backend"
)

// ClaimStaleReady claims pending ready messages that have been idle too long.
func (b *Backend) ClaimStaleReady(ctx context.Context, request backend.ClaimStaleReadyRequest) ([]backend.ReadyMessage, error) {
	if b.client == nil {
		return nil, fmt.Errorf("%w: redis client is nil", ErrInvalidRedisOptions)
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("%w: pending claim is not implemented", ErrInvalidRedisMessage)
}
