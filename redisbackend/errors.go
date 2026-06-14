package redisbackend

import "errors"

var (
	// ErrInvalidRedisOptions is returned when Redis backend options are unsafe.
	ErrInvalidRedisOptions = errors.New("goqueue redis backend: invalid options")

	// ErrInvalidRedisMessage is returned when a Redis task message cannot be encoded or decoded.
	ErrInvalidRedisMessage = errors.New("goqueue redis backend: invalid message")
)
