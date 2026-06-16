package redisbackend

import (
	"testing"
	"time"
)

func TestNewOptionsAppliesDefaults(t *testing.T) {
	options := NewOptions("redis://localhost:6379/0")

	if options.RedisURL != "redis://localhost:6379/0" {
		t.Fatalf("RedisURL = %q, want configured URL", options.RedisURL)
	}
	if options.Namespace != defaultNamespace {
		t.Fatalf("Namespace = %q, want default namespace", options.Namespace)
	}
	if options.MessageTTL != 7*24*time.Hour {
		t.Fatalf("MessageTTL = %s, want 168h", options.MessageTTL)
	}
}

func TestNewOptionsAppliesOverrides(t *testing.T) {
	options := NewOptions(
		"redis://localhost:6379/0",
		WithNamespace("payments"),
		WithMessageTTL(time.Hour),
	)

	if options.Namespace != "payments" {
		t.Fatalf("Namespace = %q, want payments", options.Namespace)
	}
	if options.MessageTTL != time.Hour {
		t.Fatalf("MessageTTL = %s, want 1h", options.MessageTTL)
	}
}

func TestNewOptionsTrimsURLAndNamespace(t *testing.T) {
	options := NewOptions(" redis://localhost:6379/0 ", WithNamespace(" payments "))

	if options.RedisURL != "redis://localhost:6379/0" {
		t.Fatalf("RedisURL = %q, want trimmed URL", options.RedisURL)
	}
	if options.Namespace != "payments" {
		t.Fatalf("Namespace = %q, want payments", options.Namespace)
	}
}
