package backend

import (
	"errors"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/task"
)

func TestEnqueueRequestValidateAcceptsTaskMessage(t *testing.T) {
	message := testTaskMessage(t)
	request := EnqueueRequest{Message: message}

	if err := request.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestEnqueueRequestValidateRejectsMissingID(t *testing.T) {
	message := testTaskMessage(t)
	message.ID = ""

	err := (EnqueueRequest{Message: message}).Validate()
	if !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func testTaskMessage(t *testing.T) task.TaskMessage {
	t.Helper()

	envelope, err := task.NewTaskEnvelope(task.TaskEnvelopeInput{
		Name:      "email.send",
		Queue:     "default",
		CreatedAt: time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("NewTaskEnvelope returned error: %v", err)
	}

	message, err := task.TaskEnvelopeToMessage(envelope, task.JSONPayloadCodec{})
	if err != nil {
		t.Fatalf("TaskEnvelopeToMessage returned error: %v", err)
	}

	return message
}
