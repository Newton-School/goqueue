package redisbackend

import "time"

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
