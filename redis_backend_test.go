package goqueue

import (
	"testing"
)

func TestAppNewRedisBackendUsesAppConfig(t *testing.T) {
	app, err := New(WithRedisURL("redis://localhost:6379/0"), WithNamespace("payments"))
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	backend, err := app.NewRedisBackend()
	if err != nil {
		t.Fatalf("NewRedisBackend returned error: %v", err)
	}
	if backend == nil {
		t.Fatal("NewRedisBackend returned nil backend")
	}
}
