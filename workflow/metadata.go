package workflow

const (
	WorkflowKindChain = "chain"
	WorkflowKindGroup = "group"
	WorkflowKindChord = "chord"
)

const (
	MetadataKindKey          = "goqueue.workflow.kind"
	MetadataChainIDKey       = "goqueue.workflow.chain_id"
	MetadataChainStepKey     = "goqueue.workflow.chain_step"
	MetadataGroupIDKey       = "goqueue.workflow.group_id"
	MetadataGroupIndexKey    = "goqueue.workflow.group_index"
	MetadataChordIDKey       = "goqueue.workflow.chord_id"
	MetadataChordCallbackKey = "goqueue.workflow.chord_callback"
)

// MergeMetadata copies base metadata and overlays workflow-reserved values.
func MergeMetadata(base map[string]string, reserved map[string]string) map[string]string {
	metadata := copyStringMap(base)
	if metadata == nil {
		metadata = map[string]string{}
	}
	for key, value := range reserved {
		metadata[key] = value
	}

	return metadata
}
