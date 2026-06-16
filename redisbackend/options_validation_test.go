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

func TestOptionsValidateRejectsZeroTTLs(t *testing.T) {
	tests := []struct {
		name string
		opt  Option
	}{
		{name: "message ttl", opt: WithMessageTTL(0)},
		{name: "state ttl", opt: WithStateTTL(0)},
		{name: "result ttl", opt: WithResultTTL(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := NewOptions("redis://localhost:6379/0", tt.opt)
			err := options.Validate()
			if !errors.Is(err, ErrInvalidRedisOptions) {
				t.Fatalf("Validate error = %v, want ErrInvalidRedisOptions", err)
			}
		})
	}
}
