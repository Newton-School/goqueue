package scheduler

import (
	"errors"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/task"
)

func TestPeriodicTaskValidateRequiresName(t *testing.T) {
	definition := validPeriodicTask()
	definition.Name = ""

	if err := definition.Validate(); !errors.Is(err, ErrInvalidPeriodicTask) {
		t.Fatalf("Validate error = %v, want ErrInvalidPeriodicTask", err)
	}
}

func TestPeriodicTaskValidateRequiresTaskName(t *testing.T) {
	definition := validPeriodicTask()
	definition.TaskName = ""

	if err := definition.Validate(); !errors.Is(err, task.ErrInvalidTaskName) {
		t.Fatalf("Validate error = %v, want ErrInvalidTaskName", err)
	}
}

func TestPeriodicTaskValidateRequiresValidQueueWhenSet(t *testing.T) {
	definition := validPeriodicTask()
	definition.Queue = "invalid queue"

	if err := definition.Validate(); !errors.Is(err, task.ErrInvalidQueueName) {
		t.Fatalf("Validate error = %v, want ErrInvalidQueueName", err)
	}
}

func TestPeriodicTaskValidateRequiresSchedule(t *testing.T) {
	definition := validPeriodicTask()
	definition.Schedule = IntervalSchedule{}

	if err := definition.Validate(); !errors.Is(err, ErrInvalidSchedule) {
		t.Fatalf("Validate error = %v, want ErrInvalidSchedule", err)
	}
}

func validPeriodicTask() PeriodicTask {
	return PeriodicTask{
		Name:        "welcome-email",
		TaskName:    "email.send",
		Queue:       "default",
		Args:        []any{"u_123"},
		Kwargs:      map[string]any{"template": "welcome"},
		Metadata:    map[string]string{"source": "scheduler"},
		Schedule:    Every(10 * time.Minute),
		Priority:    task.DefaultPriority,
		RetryPolicy: task.DefaultRetryPolicy(),
	}
}
