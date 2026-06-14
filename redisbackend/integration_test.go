package redisbackend

import (
	"os"
	"testing"
)

func redisIntegrationOptions(t *testing.T) Options {
	t.Helper()

	if os.Getenv("GOQUEUE_RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("set GOQUEUE_RUN_INTEGRATION_TESTS=true to run Redis integration tests")
	}

	redisURL := os.Getenv("GOQUEUE_REDIS_URL")
	if redisURL == "" {
		t.Fatal("GOQUEUE_REDIS_URL is required for Redis integration tests")
	}

	return NewOptions(redisURL, WithNamespace("goqueue_test"))
}
