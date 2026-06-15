package redisbackend

import "testing"

func TestParseAdvanceWorkflowChainResponseWithNextSignature(t *testing.T) {
	encoded, err := (workflowSignatureCodec{}).encode(testWorkflowSignatureRecord())
	if err != nil {
		t.Fatalf("encode returned error: %v", err)
	}

	response, err := parseAdvanceWorkflowChainResponse([]any{int64(1), int64(0), string(encoded)})
	if err != nil {
		t.Fatalf("parseAdvanceWorkflowChainResponse returned error: %v", err)
	}

	if !response.Advanced {
		t.Fatal("Advanced should be true")
	}
	if response.Completed {
		t.Fatal("Completed should be false")
	}
	if response.Next == nil {
		t.Fatal("Next should be set")
	}
	if response.Next.Name != "email.send" {
		t.Fatalf("Next name = %q, want email.send", response.Next.Name)
	}
}

func TestParseAdvanceWorkflowChainResponseWithCompletedWorkflow(t *testing.T) {
	response, err := parseAdvanceWorkflowChainResponse([]any{int64(1), int64(1), ""})
	if err != nil {
		t.Fatalf("parseAdvanceWorkflowChainResponse returned error: %v", err)
	}

	if !response.Advanced {
		t.Fatal("Advanced should be true")
	}
	if !response.Completed {
		t.Fatal("Completed should be true")
	}
	if response.Next != nil {
		t.Fatal("Next should be nil")
	}
}
