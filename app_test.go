package goqueue

import (
	"errors"
	"testing"
)

func TestNewReturnsAppWithValidatedConfig(t *testing.T) {
	app, err := New(
		WithRedisURL("rediss://redis.example.com:6380/0"),
		WithDefaultQueue("emails"),
		WithNamespace("payments"),
	)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	cfg := app.Config()
	if cfg.RedisURL != "rediss://redis.example.com:6380/0" {
		t.Fatalf("RedisURL = %q, want configured URL", cfg.RedisURL)
	}
	if cfg.DefaultQueue != "emails" {
		t.Fatalf("DefaultQueue = %q, want emails", cfg.DefaultQueue)
	}
	if cfg.Namespace != "payments" {
		t.Fatalf("Namespace = %q, want payments", cfg.Namespace)
	}
}

func TestNewRejectsMissingRedisURL(t *testing.T) {
	_, err := New()
	if !errors.Is(err, ErrMissingRedisURL) {
		t.Fatalf("New error = %v, want ErrMissingRedisURL", err)
	}
}
