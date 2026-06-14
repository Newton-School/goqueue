package redisbackend

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Newton-School/goqueue/task"
)

const defaultNamespace = "goqueue"

// Options configures the Redis backend.
type Options struct {
	RedisURL   string
	Namespace  string
	MessageTTL time.Duration
	StateTTL   time.Duration
	ResultTTL  time.Duration
}

// Option customizes Options.
type Option func(*Options)

// NewOptions builds Redis backend options with safe defaults.
func NewOptions(redisURL string, opts ...Option) Options {
	options := Options{
		RedisURL:   redisURL,
		Namespace:  defaultNamespace,
		MessageTTL: 7 * 24 * time.Hour,
		StateTTL:   24 * time.Hour,
		ResultTTL:  24 * time.Hour,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}

	return options
}

// Validate verifies that options can safely initialize a Redis backend.
func (o Options) Validate() error {
	if strings.TrimSpace(o.RedisURL) == "" {
		return fmt.Errorf("%w: redis url is required", ErrInvalidRedisOptions)
	}

	parsed, err := url.Parse(o.RedisURL)
	if err != nil || parsed.Host == "" {
		return fmt.Errorf("%w: redis url must include host", ErrInvalidRedisOptions)
	}
	if parsed.Scheme != "redis" && parsed.Scheme != "rediss" {
		return fmt.Errorf("%w: redis url must use redis:// or rediss://", ErrInvalidRedisOptions)
	}

	if err := task.ValidateQueueName(o.Namespace); err != nil {
		return fmt.Errorf("%w: namespace: %v", ErrInvalidRedisOptions, err)
	}

	if o.MessageTTL < 0 {
		return fmt.Errorf("%w: message ttl cannot be negative", ErrInvalidRedisOptions)
	}
	if o.StateTTL < 0 {
		return fmt.Errorf("%w: state ttl cannot be negative", ErrInvalidRedisOptions)
	}
	if o.ResultTTL < 0 {
		return fmt.Errorf("%w: result ttl cannot be negative", ErrInvalidRedisOptions)
	}

	return nil
}

// WithNamespace configures the Redis key namespace.
func WithNamespace(namespace string) Option {
	return func(options *Options) {
		options.Namespace = namespace
	}
}

// WithMessageTTL configures task message retention.
func WithMessageTTL(ttl time.Duration) Option {
	return func(options *Options) {
		options.MessageTTL = ttl
	}
}

// WithStateTTL configures task state retention.
func WithStateTTL(ttl time.Duration) Option {
	return func(options *Options) {
		options.StateTTL = ttl
	}
}

// WithResultTTL configures task result retention.
func WithResultTTL(ttl time.Duration) Option {
	return func(options *Options) {
		options.ResultTTL = ttl
	}
}
