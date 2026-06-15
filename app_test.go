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

func TestAppRegisterTaskStoresHandler(t *testing.T) {
	app, err := New(WithRedisURL("redis://localhost:6379/0"))
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	handler := TaskHandlerFunc(func(HandlerContext, TaskPayload) (TaskResult, error) {
		return SucceededResult(nil), nil
	})
	if err := app.RegisterTask("email.send", handler); err != nil {
		t.Fatalf("RegisterTask returned error: %v", err)
	}

	registered, err := app.LookupTask("email.send")
	if err != nil {
		t.Fatalf("LookupTask returned error: %v", err)
	}
	if registered == nil {
		t.Fatal("LookupTask returned nil handler")
	}
}

func TestAppTaskNamesReturnsRegisteredNames(t *testing.T) {
	app, err := New(WithRedisURL("redis://localhost:6379/0"))
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	handler := TaskHandlerFunc(func(HandlerContext, TaskPayload) (TaskResult, error) {
		return SucceededResult(nil), nil
	})
	if err := app.RegisterTask("email.send", handler); err != nil {
		t.Fatalf("RegisterTask returned error: %v", err)
	}

	names := app.TaskNames()
	if len(names) != 1 || names[0] != "email.send" {
		t.Fatalf("TaskNames = %#v, want registered task name", names)
	}
}

func TestAppNewProducerUsesConfigDefaults(t *testing.T) {
	app, err := New(
		WithRedisURL("redis://localhost:6379/0"),
		WithDefaultQueue("billing"),
		WithNamespace("inbox"),
	)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	producer, err := app.NewProducer()
	if err != nil {
		t.Fatalf("NewProducer returned error: %v", err)
	}
	if producer == nil {
		t.Fatal("NewProducer returned nil")
	}
}

func TestAppNewWorkerUsesConfigDefaults(t *testing.T) {
	app, err := New(WithRedisURL("redis://localhost:6379/0"), WithDefaultQueue("billing"))
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	worker, err := app.NewWorker()
	if err != nil {
		t.Fatalf("NewWorker returned error: %v", err)
	}
	if worker == nil {
		t.Fatal("NewWorker returned nil")
	}
}

func TestAppNewSchedulerUsesConfigDefaults(t *testing.T) {
	app, err := New(WithRedisURL("redis://localhost:6379/0"), WithDefaultQueue("billing"))
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	scheduler, err := app.NewScheduler()
	if err != nil {
		t.Fatalf("NewScheduler returned error: %v", err)
	}
	if scheduler == nil {
		t.Fatal("NewScheduler returned nil")
	}
}
