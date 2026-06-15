package redisbackend

import (
	"context"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
)

func TestPeriodicTaskRedisLifecycle(t *testing.T) {
	ctx := context.Background()
	options := redisIntegrationOptions(t)
	b, err := New(options)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	defer b.Close()
	cleanupIntegrationNamespace(ctx, t, b)

	record := testPeriodicTaskRecord()
	record.NextDueAt = time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC)

	if err := b.UpsertPeriodicTask(ctx, backend.UpsertPeriodicTaskRequest{Record: record}); err != nil {
		t.Fatalf("UpsertPeriodicTask returned error: %v", err)
	}

	due, err := b.ListDuePeriodicTasks(ctx, backend.ListDuePeriodicTasksRequest{
		Now:         record.NextDueAt,
		Limit:       10,
		SchedulerID: "scheduler-1",
		LockTTL:     time.Minute,
	})
	if err != nil {
		t.Fatalf("ListDuePeriodicTasks returned error: %v", err)
	}
	if len(due) != 1 {
		t.Fatalf("due count = %d, want 1", len(due))
	}
	if due[0].Record.Name != record.Name {
		t.Fatalf("due name = %q, want %q", due[0].Record.Name, record.Name)
	}

	lockedAgain, err := b.ListDuePeriodicTasks(ctx, backend.ListDuePeriodicTasksRequest{
		Now:         record.NextDueAt,
		Limit:       10,
		SchedulerID: "scheduler-2",
		LockTTL:     time.Minute,
	})
	if err != nil {
		t.Fatalf("second ListDuePeriodicTasks returned error: %v", err)
	}
	if len(lockedAgain) != 0 {
		t.Fatalf("second due count = %d, want 0", len(lockedAgain))
	}

	nextDue := record.NextDueAt.Add(10 * time.Minute)
	if err := b.MarkPeriodicTaskDispatched(ctx, backend.MarkPeriodicTaskDispatchedRequest{
		Name:             record.Name,
		LockToken:        due[0].LockToken,
		DispatchedTaskID: "11111111-1111-4111-8111-111111111111",
		DispatchedAt:     record.NextDueAt,
		NextDueAt:        nextDue,
	}); err != nil {
		t.Fatalf("MarkPeriodicTaskDispatched returned error: %v", err)
	}

	notDue, err := b.ListDuePeriodicTasks(ctx, backend.ListDuePeriodicTasksRequest{
		Now:         record.NextDueAt,
		Limit:       10,
		SchedulerID: "scheduler-1",
		LockTTL:     time.Minute,
	})
	if err != nil {
		t.Fatalf("ListDuePeriodicTasks after mark returned error: %v", err)
	}
	if len(notDue) != 0 {
		t.Fatalf("not due count = %d, want 0", len(notDue))
	}

	if err := b.DeletePeriodicTask(ctx, backend.DeletePeriodicTaskRequest{Name: record.Name}); err != nil {
		t.Fatalf("DeletePeriodicTask returned error: %v", err)
	}
}
