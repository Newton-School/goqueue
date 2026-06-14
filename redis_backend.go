package goqueue

import "github.com/Newton-School/goqueue/redisbackend"

// NewRedisBackend creates a Redis backend from the app configuration.
func (a *App) NewRedisBackend(opts ...redisbackend.BackendOption) (*redisbackend.Backend, error) {
	return redisbackend.New(
		redisbackend.NewOptions(
			a.config.RedisURL,
			redisbackend.WithNamespace(a.config.Namespace),
		),
		opts...,
	)
}
