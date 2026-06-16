package redisbackend

import (
	"context"
	"errors"
	"testing"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
	"github.com/redis/go-redis/v9"
)

func TestGetTaskMessageRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	_, err := b.GetTaskMessage(context.Background(), task.TaskID("123e4567-e89b-42d3-a456-556642440111"))
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("GetTaskMessage error = %v, want ErrInvalidRedisOptions", err)
	}
}

func TestGetTaskMessageRejectsInvalidTaskID(t *testing.T) {
	b := &Backend{
		options: NewOptions("redis://localhost:6379/0"),
		keys:    newKeyBuilder("goqueue"),
		client:  redis.NewClient(&redis.Options{Addr: "localhost:6379"}),
	}

	_, err := b.GetTaskMessage(context.Background(), task.TaskID(""))
	if !errors.Is(err, task.ErrInvalidTaskID) {
		t.Fatalf("GetTaskMessage error = %v, want invalid task id", err)
	}
}

func TestReadDeadLetterRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	_, err := b.ReadDeadLetter(context.Background(), "default", "1-0")
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("ReadDeadLetter error = %v, want ErrInvalidRedisOptions", err)
	}
}

func TestReadDeadLetterRejectsInvalidQueue(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	_, err := b.ReadDeadLetter(context.Background(), "bad queue", "1-0")
	if err == nil {
		t.Fatal("ReadDeadLetter expected error")
	}
}

func TestReadDeadLetterReturnsBackendErrForMissingRecord(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	_, err := b.ReadDeadLetter(context.Background(), "default", "")
	if err == nil {
		t.Fatal("ReadDeadLetter expected error")
	}
}

func TestDeleteDeadLettersRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	deleted, err := b.DeleteDeadLetters(context.Background(), "default")
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("DeleteDeadLetters error = %v, want ErrInvalidRedisOptions", err)
	}
	if deleted != 0 {
		t.Fatalf("DeleteDeadLetters deleted = %d, want 0", deleted)
	}
}

func TestPurgeQueueRejectsInvalidQueue(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	_, err := b.PurgeQueue(context.Background(), backend.PurgeQueueRequest{Queue: "bad queue"})
	if err == nil {
		t.Fatal("PurgeQueue expected error")
	}
}
