package workflow

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

func TestNewCanvasRequiresBackend(t *testing.T) {
	_, err := NewCanvas(nil)
	if !errors.Is(err, ErrNilBackend) {
		t.Fatalf("NewCanvas error = %v, want ErrNilBackend", err)
	}
}

func TestNewCanvasAppliesOptions(t *testing.T) {
	now := time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC)
	canvas, err := NewCanvas(
		&fakeBackend{},
		WithCanvasDefaultQueue("critical"),
		WithCanvasNow(func() time.Time { return now }),
	)
	if err != nil {
		t.Fatalf("NewCanvas returned error: %v", err)
	}

	if canvas.defaultQueue != "critical" {
		t.Fatalf("defaultQueue = %q, want critical", canvas.defaultQueue)
	}
	if got := canvas.now(); !got.Equal(now) {
		t.Fatalf("now = %v, want %v", got, now)
	}
}

func TestNewWorkflowIDGeneratesUUIDLikeID(t *testing.T) {
	id, err := newWorkflowID()
	if err != nil {
		t.Fatalf("newWorkflowID returned error: %v", err)
	}

	if len(id.String()) != 36 || !strings.Contains(id.String(), "-") {
		t.Fatalf("id = %q, want UUID-like id", id)
	}
}

func TestCanvasApplySignatureDispatchesTask(t *testing.T) {
	backend := &fakeBackend{}
	canvas, err := NewCanvas(backend, WithCanvasDefaultQueue("critical"))
	if err != nil {
		t.Fatalf("NewCanvas returned error: %v", err)
	}
	signature := validSignature()
	signature.Queue = ""

	result, err := canvas.ApplySignature(context.Background(), signature)
	if err != nil {
		t.Fatalf("ApplySignature returned error: %v", err)
	}

	if result == nil || result.ID() == "" {
		t.Fatal("ApplySignature should return async result with task id")
	}
	if len(backend.enqueueReadyRequests) != 1 {
		t.Fatalf("enqueue ready calls = %d, want 1", len(backend.enqueueReadyRequests))
	}
	message := backend.enqueueReadyRequests[0].Message
	if message.Name != "email.send" {
		t.Fatalf("message name = %q, want email.send", message.Name)
	}
	if message.Queue != "critical" {
		t.Fatalf("message queue = %q, want critical", message.Queue)
	}
	if message.Metadata["source"] != "workflow" {
		t.Fatalf("message metadata source = %q, want workflow", message.Metadata["source"])
	}
}

func TestCanvasApplyChainStoresChainBeforeDispatch(t *testing.T) {
	backend := &fakeBackend{}
	canvas, err := NewCanvas(backend, WithCanvasDefaultQueue("critical"))
	if err != nil {
		t.Fatalf("NewCanvas returned error: %v", err)
	}
	first := validSignature()
	first.Queue = ""
	second := validSignature()
	second.Name = "email.audit"
	second.Queue = ""

	result, err := canvas.ApplyChain(context.Background(), Chain{Signatures: []Signature{first, second}})
	if err != nil {
		t.Fatalf("ApplyChain returned error: %v", err)
	}

	if result.WorkflowID == "" || result.FirstTask == "" {
		t.Fatalf("result = %+v, want workflow and first task ids", result)
	}
	if len(backend.savedChains) != 1 {
		t.Fatalf("saved chains = %d, want 1", len(backend.savedChains))
	}
	if backend.savedChains[0].ID != result.WorkflowID.String() {
		t.Fatalf("chain id = %q, want %q", backend.savedChains[0].ID, result.WorkflowID)
	}
	if len(backend.enqueueReadyRequests) != 1 {
		t.Fatalf("enqueue calls = %d, want 1", len(backend.enqueueReadyRequests))
	}
	message := backend.enqueueReadyRequests[0].Message
	if message.ID != result.FirstTask.String() {
		t.Fatalf("message id = %q, want %q", message.ID, result.FirstTask)
	}
	if message.Metadata[MetadataKindKey] != WorkflowKindChain {
		t.Fatalf("workflow kind = %q, want chain", message.Metadata[MetadataKindKey])
	}
	if message.Metadata[MetadataChainIDKey] != result.WorkflowID.String() {
		t.Fatalf("chain id metadata = %q, want workflow id", message.Metadata[MetadataChainIDKey])
	}
	if message.Metadata[MetadataChainStepKey] != "0" {
		t.Fatalf("chain step metadata = %q, want 0", message.Metadata[MetadataChainStepKey])
	}
}

func TestCanvasApplyGroupStoresGroupAndDispatchesChildren(t *testing.T) {
	backend := &fakeBackend{}
	canvas, err := NewCanvas(backend, WithCanvasDefaultQueue("critical"))
	if err != nil {
		t.Fatalf("NewCanvas returned error: %v", err)
	}
	first := validSignature()
	first.Queue = ""
	second := validSignature()
	second.Name = "email.audit"
	second.Queue = ""

	result, err := canvas.ApplyGroup(context.Background(), Group{Signatures: []Signature{first, second}})
	if err != nil {
		t.Fatalf("ApplyGroup returned error: %v", err)
	}

	if result.GroupID == "" || len(result.TaskIDs) != 2 {
		t.Fatalf("result = %+v, want group id and two task ids", result)
	}
	if len(backend.savedGroups) != 1 {
		t.Fatalf("saved groups = %d, want 1", len(backend.savedGroups))
	}
	if backend.savedGroups[0].ID != result.GroupID.String() {
		t.Fatalf("group id = %q, want %q", backend.savedGroups[0].ID, result.GroupID)
	}
	if len(backend.enqueueReadyRequests) != 2 {
		t.Fatalf("enqueue calls = %d, want 2", len(backend.enqueueReadyRequests))
	}
	for index, request := range backend.enqueueReadyRequests {
		message := request.Message
		if message.ID != result.TaskIDs[index].String() {
			t.Fatalf("message %d id = %q, want %q", index, message.ID, result.TaskIDs[index])
		}
		if message.Metadata[MetadataKindKey] != WorkflowKindGroup {
			t.Fatalf("message %d kind = %q, want group", index, message.Metadata[MetadataKindKey])
		}
		if message.Metadata[MetadataGroupIDKey] != result.GroupID.String() {
			t.Fatalf("message %d group id = %q, want group id", index, message.Metadata[MetadataGroupIDKey])
		}
		if message.Metadata[MetadataGroupIndexKey] != workflowIndexMetadata(index) {
			t.Fatalf("message %d group index = %q, want %d", index, message.Metadata[MetadataGroupIndexKey], index)
		}
	}
}

func TestCanvasApplyChordStoresCallbackAndDispatchesHeader(t *testing.T) {
	backend := &fakeBackend{}
	canvas, err := NewCanvas(backend, WithCanvasDefaultQueue("critical"))
	if err != nil {
		t.Fatalf("NewCanvas returned error: %v", err)
	}
	header := validSignature()
	header.Queue = ""
	callback := validSignature()
	callback.Name = "email.complete"
	callback.Queue = ""

	result, err := canvas.ApplyChord(context.Background(), Chord{
		Header:   Group{Signatures: []Signature{header}},
		Callback: callback,
	})
	if err != nil {
		t.Fatalf("ApplyChord returned error: %v", err)
	}

	if result.GroupID == "" || len(result.TaskIDs) != 1 {
		t.Fatalf("result = %+v, want group id and one task id", result)
	}
	if len(backend.savedGroups) != 1 {
		t.Fatalf("saved groups = %d, want 1", len(backend.savedGroups))
	}
	if backend.savedGroups[0].Callback == nil {
		t.Fatal("saved group callback should be set")
	}
	if backend.savedGroups[0].Callback.Name != "email.complete" {
		t.Fatalf("callback name = %q, want email.complete", backend.savedGroups[0].Callback.Name)
	}
	message := backend.enqueueReadyRequests[0].Message
	if message.Metadata[MetadataKindKey] != WorkflowKindChord {
		t.Fatalf("workflow kind = %q, want chord", message.Metadata[MetadataKindKey])
	}
	if message.Metadata[MetadataChordIDKey] != result.GroupID.String() {
		t.Fatalf("chord id = %q, want group id", message.Metadata[MetadataChordIDKey])
	}
	if message.Metadata[MetadataGroupIDKey] != result.GroupID.String() {
		t.Fatalf("group id = %q, want group id", message.Metadata[MetadataGroupIDKey])
	}
}

type fakeBackend struct {
	mu                   sync.Mutex
	enqueueReadyRequests []backend.EnqueueRequest
	setStateRequests     []backend.TaskStateRecord
	savedChains          []backend.WorkflowChainRecord
	savedGroups          []backend.WorkflowGroupRecord
}

func (f *fakeBackend) EnqueueReady(_ context.Context, request backend.EnqueueRequest) (backend.EnqueueResponse, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.enqueueReadyRequests = append(f.enqueueReadyRequests, request)
	return backend.EnqueueResponse{TaskID: task.TaskID(request.Message.ID), StreamID: "1-0"}, nil
}
func (f *fakeBackend) EnqueueScheduled(context.Context, backend.EnqueueRequest) (backend.EnqueueResponse, error) {
	return backend.EnqueueResponse{}, nil
}
func (f *fakeBackend) MoveDueScheduled(context.Context, backend.MoveDueScheduledRequest) ([]backend.MovedScheduledMessage, error) {
	return nil, nil
}
func (f *fakeBackend) EnsureConsumerGroup(context.Context, backend.ConsumerGroupRequest) error {
	return nil
}
func (f *fakeBackend) ReadReady(context.Context, backend.ReadReadyRequest) ([]backend.ReadyMessage, error) {
	return nil, nil
}
func (f *fakeBackend) ClaimStaleReady(context.Context, backend.ClaimStaleReadyRequest) ([]backend.ReadyMessage, error) {
	return nil, nil
}
func (f *fakeBackend) Ack(context.Context, backend.AckRequest) error { return nil }
func (f *fakeBackend) EnqueueDeadLetter(context.Context, backend.DeadLetterRequest) (backend.DeadLetterRecord, error) {
	return backend.DeadLetterRecord{}, nil
}
func (f *fakeBackend) ReadDeadLetters(context.Context, backend.ReadDeadLettersRequest) ([]backend.DeadLetterRecord, error) {
	return nil, nil
}
func (f *fakeBackend) UpsertPeriodicTask(context.Context, backend.UpsertPeriodicTaskRequest) error {
	return nil
}
func (f *fakeBackend) DeletePeriodicTask(context.Context, backend.DeletePeriodicTaskRequest) error {
	return nil
}
func (f *fakeBackend) ListDuePeriodicTasks(context.Context, backend.ListDuePeriodicTasksRequest) ([]backend.DuePeriodicTask, error) {
	return nil, nil
}
func (f *fakeBackend) MarkPeriodicTaskDispatched(context.Context, backend.MarkPeriodicTaskDispatchedRequest) error {
	return nil
}
func (f *fakeBackend) SaveWorkflowChain(_ context.Context, record backend.WorkflowChainRecord) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.savedChains = append(f.savedChains, record)
	return nil
}
func (f *fakeBackend) AdvanceWorkflowChain(context.Context, backend.AdvanceWorkflowChainRequest) (backend.AdvanceWorkflowChainResponse, error) {
	return backend.AdvanceWorkflowChainResponse{}, nil
}
func (f *fakeBackend) SaveWorkflowGroup(_ context.Context, record backend.WorkflowGroupRecord) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.savedGroups = append(f.savedGroups, record)
	return nil
}
func (f *fakeBackend) RecordWorkflowTaskCompleted(context.Context, backend.RecordWorkflowTaskCompletedRequest) (backend.WorkflowGroupProgress, error) {
	return backend.WorkflowGroupProgress{}, nil
}
func (f *fakeBackend) SetTaskState(_ context.Context, record backend.TaskStateRecord) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.setStateRequests = append(f.setStateRequests, record)
	return nil
}
func (f *fakeBackend) GetTaskState(context.Context, task.TaskID) (backend.TaskStateRecord, error) {
	return backend.TaskStateRecord{}, nil
}
func (f *fakeBackend) SaveTaskResult(context.Context, backend.TaskResultRecord) error {
	return nil
}
func (f *fakeBackend) GetTaskResult(context.Context, task.TaskID) (backend.TaskResultRecord, error) {
	return backend.TaskResultRecord{}, nil
}
func (f *fakeBackend) ForgetTaskResult(context.Context, task.TaskID) error { return nil }
func (f *fakeBackend) QueueStats(context.Context, backend.QueueStatsRequest) (backend.QueueStats, error) {
	return backend.QueueStats{}, nil
}
func (f *fakeBackend) Ping(context.Context) error { return nil }
func (f *fakeBackend) Close() error               { return nil }
