package redisbackend

import "errors"

var (
	// ErrInvalidRedisOptions is returned when Redis backend options are unsafe.
	ErrInvalidRedisOptions = errors.New("goqueue redis backend: invalid options")
)
