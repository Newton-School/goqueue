package redisbackend

import (
	"errors"
	"testing"
	"time"
)

func TestOptionsValidateAcceptsRedisURL(t *testing.T) {
	options := NewOptions("rediss://redis.example.com:6380/0")
	if err := options.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestOptionsValidateRejectsMissingRedisURL(t *testing.T) {
	options := NewOptions("")
	err := options.Validate()
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("Validate error = %v, want ErrInvalidRedisOptions", err)
	}
}

func TestOptionsValidateRejectsNegativeTTL(t *testing.T) {
	options := NewOptions("redis://localhost:6379/0", WithMessageTTL(-time.Second))
	err := options.Validate()
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("Validate error = %v, want ErrInvalidRedisOptions", err)
	}
}
