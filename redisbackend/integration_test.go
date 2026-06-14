package redisbackend

import (
	"context"
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

func cleanupIntegrationNamespace(ctx context.Context, t *testing.T, b *Backend) {
	t.Helper()

	keys, err := b.client.Keys(ctx, b.options.Namespace+":*").Result()
	if err != nil {
		t.Fatalf("list integration keys: %v", err)
	}
	if len(keys) == 0 {
		return
	}
	if err := b.client.Del(ctx, keys...).Err(); err != nil {
		t.Fatalf("delete integration keys: %v", err)
	}
}
