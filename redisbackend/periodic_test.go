package redisbackend

import (
	"context"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
)

func TestUpsertPeriodicTaskRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	err := b.UpsertPeriodicTask(context.Background(), backend.UpsertPeriodicTaskRequest{Record: testPeriodicTaskRecord()})
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("UpsertPeriodicTask error = %v, want ErrInvalidRedisOptions", err)
	}
}

func TestDeletePeriodicTaskRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	err := b.DeletePeriodicTask(context.Background(), backend.DeletePeriodicTaskRequest{Name: "welcome-email"})
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("DeletePeriodicTask error = %v, want ErrInvalidRedisOptions", err)
	}
}

func TestListDuePeriodicTasksRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	_, err := b.ListDuePeriodicTasks(context.Background(), backend.ListDuePeriodicTasksRequest{
		Now:         time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC),
		Limit:       10,
		SchedulerID: "scheduler-1",
		LockTTL:     time.Minute,
	})
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("ListDuePeriodicTasks error = %v, want ErrInvalidRedisOptions", err)
	}
}

func TestMarkPeriodicTaskDispatchedRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	err := b.MarkPeriodicTaskDispatched(context.Background(), backend.MarkPeriodicTaskDispatchedRequest{
		Name:             "welcome-email",
		LockToken:        "token",
		DispatchedTaskID: "11111111-1111-4111-8111-111111111111",
		DispatchedAt:     time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC),
		NextDueAt:        time.Date(2026, time.June, 15, 10, 10, 0, 0, time.UTC),
	})
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("MarkPeriodicTaskDispatched error = %v, want ErrInvalidRedisOptions", err)
	}
}

func TestNewPeriodicLockTokenReturnsHexToken(t *testing.T) {
	first, err := newPeriodicLockToken()
	if err != nil {
		t.Fatalf("newPeriodicLockToken returned error: %v", err)
	}
	second, err := newPeriodicLockToken()
	if err != nil {
		t.Fatalf("newPeriodicLockToken returned error: %v", err)
	}

	if len(first) != 32 {
		t.Fatalf("token length = %d, want 32", len(first))
	}
	if _, err := hex.DecodeString(first); err != nil {
		t.Fatalf("token should be hex: %v", err)
	}
	if first == second {
		t.Fatal("lock tokens should be unique")
	}
}
