package redisbackend

import "github.com/redis/go-redis/v9"

// Backend stores and reads task messages from Redis.
type Backend struct {
	options Options
	client  redis.UniversalClient
	keys    keyBuilder
}

// New creates a Redis backend.
func New(options Options, opts ...BackendOption) (*Backend, error) {
	if err := options.Validate(); err != nil {
		return nil, err
	}

	backend := &Backend{
		options: options,
		keys:    newKeyBuilder(options.Namespace),
	}
	for _, opt := range opts {
		if opt != nil {
			opt(backend)
		}
	}

	if backend.client == nil {
		parsed, err := redis.ParseURL(options.RedisURL)
		if err != nil {
			return nil, err
		}
		backend.client = redis.NewClient(parsed)
	}

	return backend, nil
}

// BackendOption customizes a Redis backend.
type BackendOption func(*Backend)

// WithClient injects a Redis client for tests or advanced users.
func WithClient(client redis.UniversalClient) BackendOption {
	return func(backend *Backend) {
		backend.client = client
	}
}
