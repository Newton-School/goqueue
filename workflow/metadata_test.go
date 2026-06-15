package workflow

import "testing"

func TestWorkflowMetadataMergesReservedValues(t *testing.T) {
	metadata := MergeMetadata(map[string]string{"source": "user"}, map[string]string{
		MetadataKindKey:    WorkflowKindChain,
		MetadataChainIDKey: "chain-1",
	})

	if metadata["source"] != "user" {
		t.Fatalf("source = %q, want user", metadata["source"])
	}
	if metadata[MetadataKindKey] != WorkflowKindChain {
		t.Fatalf("kind = %q, want chain", metadata[MetadataKindKey])
	}
	if metadata[MetadataChainIDKey] != "chain-1" {
		t.Fatalf("chain id = %q, want chain-1", metadata[MetadataChainIDKey])
	}
}

func TestWorkflowMetadataReservedValuesOverrideBase(t *testing.T) {
	metadata := MergeMetadata(map[string]string{MetadataKindKey: "user"}, map[string]string{MetadataKindKey: WorkflowKindGroup})

	if metadata[MetadataKindKey] != WorkflowKindGroup {
		t.Fatalf("kind = %q, want group", metadata[MetadataKindKey])
	}
}
