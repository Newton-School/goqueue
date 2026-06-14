package goqueue

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	defaultQueueName = "default"
	defaultNamespace = "goqueue"
	maxNameLength    = 128
)

// Config contains application-level settings for producers and workers.
type Config struct {
	RedisURL     string
	DefaultQueue string
	Namespace    string
}

// Option customizes Config.
type Option func(*Config)

// NewConfig returns a Config with production-safe defaults and applied options.
func NewConfig(opts ...Option) Config {
	config := Config{
		DefaultQueue: defaultQueueName,
		Namespace:    defaultNamespace,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&config)
		}
	}

	config.RedisURL = strings.TrimSpace(config.RedisURL)
	config.DefaultQueue = strings.TrimSpace(config.DefaultQueue)
	config.Namespace = strings.TrimSpace(config.Namespace)

	return config
}

// WithRedisURL configures the Redis connection URL used by goqueue.
func WithRedisURL(redisURL string) Option {
	return func(config *Config) {
		config.RedisURL = redisURL
	}
}

// WithDefaultQueue configures the queue used when a task does not override it.
func WithDefaultQueue(queue string) Option {
	return func(config *Config) {
		config.DefaultQueue = queue
	}
}

// WithNamespace configures the Redis key namespace used by this application.
func WithNamespace(namespace string) Option {
	return func(config *Config) {
		config.Namespace = namespace
	}
}

// Validate checks that Config can safely initialize a goqueue application.
func (c Config) Validate() error {
	if c.RedisURL == "" {
		return fmt.Errorf("%w: configure RedisURL before creating a goqueue app", ErrMissingRedisURL)
	}

	parsed, err := url.Parse(c.RedisURL)
	if err != nil || parsed.Host == "" {
		return fmt.Errorf("%w: expected redis:// or rediss:// URL", ErrInvalidRedisURL)
	}

	if parsed.Scheme != "redis" && parsed.Scheme != "rediss" {
		return fmt.Errorf("%w: unsupported scheme %q", ErrInvalidRedisURL, parsed.Scheme)
	}

	if !validName(c.DefaultQueue) {
		return fmt.Errorf("%w: queue names must use 1-%d characters from [A-Za-z0-9._:-]", ErrInvalidQueueName, maxNameLength)
	}

	if !validName(c.Namespace) {
		return fmt.Errorf("%w: namespaces must use 1-%d characters from [A-Za-z0-9._:-]", ErrInvalidNamespace, maxNameLength)
	}

	return nil
}

// RedactedRedisURL returns the Redis URL with credentials hidden for logs.
func (c Config) RedactedRedisURL() string {
	if c.RedisURL == "" {
		return ""
	}

	parsed, err := url.Parse(c.RedisURL)
	if err != nil {
		return "<invalid redis url>"
	}

	return parsed.Redacted()
}

func validName(value string) bool {
	if value == "" || len(value) > maxNameLength {
		return false
	}

	for _, char := range value {
		if char >= 'a' && char <= 'z' {
			continue
		}
		if char >= 'A' && char <= 'Z' {
			continue
		}
		if char >= '0' && char <= '9' {
			continue
		}
		switch char {
		case '.', '_', ':', '-':
			continue
		default:
			return false
		}
	}

	return true
}
