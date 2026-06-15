package scheduler

const (
	// PeriodicMetadataNameKey stores the periodic definition name on dispatched tasks.
	PeriodicMetadataNameKey = "goqueue.periodic.name"

	// PeriodicMetadataDueAtKey stores the due timestamp that caused dispatch.
	PeriodicMetadataDueAtKey = "goqueue.periodic.due_at"
)

func periodicDispatchMetadata(base map[string]string, periodicName string, dueAt string) map[string]string {
	metadata := copyStringMap(base)
	if metadata == nil {
		metadata = map[string]string{}
	}
	metadata[PeriodicMetadataNameKey] = periodicName
	metadata[PeriodicMetadataDueAtKey] = dueAt
	return metadata
}
