package redisbackend

import "testing"

func TestKeyBuilderWorkflowChainMeta(t *testing.T) {
	keys := newKeyBuilder("tenant")

	if got := keys.workflowChainMeta("chain-1"); got != "tenant:workflow:chain:chain-1:meta" {
		t.Fatalf("workflowChainMeta = %q", got)
	}
}

func TestKeyBuilderWorkflowChainSignatures(t *testing.T) {
	keys := newKeyBuilder("tenant")

	if got := keys.workflowChainSignatures("chain-1"); got != "tenant:workflow:chain:chain-1:signatures" {
		t.Fatalf("workflowChainSignatures = %q", got)
	}
}

func TestKeyBuilderWorkflowGroupMeta(t *testing.T) {
	keys := newKeyBuilder("tenant")

	if got := keys.workflowGroupMeta("group-1"); got != "tenant:workflow:group:group-1:meta" {
		t.Fatalf("workflowGroupMeta = %q", got)
	}
}

func TestKeyBuilderWorkflowGroupCompleted(t *testing.T) {
	keys := newKeyBuilder("tenant")

	if got := keys.workflowGroupCompleted("group-1"); got != "tenant:workflow:group:group-1:completed" {
		t.Fatalf("workflowGroupCompleted = %q", got)
	}
}

func TestKeyBuilderWorkflowGroupCallback(t *testing.T) {
	keys := newKeyBuilder("tenant")

	if got := keys.workflowGroupCallback("group-1"); got != "tenant:workflow:group:group-1:callback" {
		t.Fatalf("workflowGroupCallback = %q", got)
	}
}
