package goqueue

import (
	"errors"
	"strings"
	"testing"
)

func TestNewConfigSetsProductionSafeDefaults(t *testing.T) {
	cfg := NewConfig(WithRedisURL("redis://localhost:6379/0"))

	if cfg.RedisURL != "redis://localhost:6379/0" {
		t.Fatalf("RedisURL = %q, want configured URL", cfg.RedisURL)
	}
	if cfg.DefaultQueue != "default" {
		t.Fatalf("DefaultQueue = %q, want default", cfg.DefaultQueue)
	}
	if cfg.Namespace != "goqueue" {
		t.Fatalf("Namespace = %q, want goqueue", cfg.Namespace)
	}
}

func TestConfigValidateRequiresRedisURL(t *testing.T) {
	cfg := NewConfig()

	err := cfg.Validate()
	if !errors.Is(err, ErrMissingRedisURL) {
		t.Fatalf("Validate error = %v, want ErrMissingRedisURL", err)
	}
}

func TestConfigValidateRejectsUnsupportedRedisScheme(t *testing.T) {
	cfg := NewConfig(WithRedisURL("http://localhost:6379/0"))

	err := cfg.Validate()
	if !errors.Is(err, ErrInvalidRedisURL) {
		t.Fatalf("Validate error = %v, want ErrInvalidRedisURL", err)
	}
}

func TestConfigValidateRejectsInvalidQueueName(t *testing.T) {
	cfg := NewConfig(
		WithRedisURL("redis://localhost:6379/0"),
		WithDefaultQueue("emails pending"),
	)

	err := cfg.Validate()
	if !errors.Is(err, ErrInvalidQueueName) {
		t.Fatalf("Validate error = %v, want ErrInvalidQueueName", err)
	}
}

func TestConfigValidateRejectsInvalidNamespace(t *testing.T) {
	cfg := NewConfig(
		WithRedisURL("redis://localhost:6379/0"),
		WithNamespace("../internal"),
	)

	err := cfg.Validate()
	if !errors.Is(err, ErrInvalidNamespace) {
		t.Fatalf("Validate error = %v, want ErrInvalidNamespace", err)
	}
}

func TestConfigRedactedRedisURLHidesCredentials(t *testing.T) {
	cfg := NewConfig(WithRedisURL("redis://worker:secret@redis.example.com:6379/0"))

	redacted := cfg.RedactedRedisURL()
	if strings.Contains(redacted, "secret") {
		t.Fatalf("RedactedRedisURL leaked password: %q", redacted)
	}
	if redacted != "redis://worker:xxxxx@redis.example.com:6379/0" {
		t.Fatalf("RedactedRedisURL = %q, want password redacted", redacted)
	}
}
