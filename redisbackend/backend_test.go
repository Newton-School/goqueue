package redisbackend

import (
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestNewBuildsBackendWithValidatedOptions(t *testing.T) {
	backend, err := New(NewOptions("redis://localhost:6379/0"), WithClient(redis.NewClient(&redis.Options{Addr: "localhost:6379"})))
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	if backend == nil {
		t.Fatal("New returned nil backend")
	}
}

func TestNewRejectsInvalidOptions(t *testing.T) {
	_, err := New(NewOptions(""), WithClient(redis.NewClient(&redis.Options{Addr: "localhost:6379"})))
	if err == nil {
		t.Fatal("New should reject invalid options")
	}
}
