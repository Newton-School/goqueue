package scheduler

import "testing"

func TestPeriodicDispatchMetadataMergesDefinitionMetadata(t *testing.T) {
	metadata := periodicDispatchMetadata(
		map[string]string{"source": "scheduler"},
		"welcome-email",
		"2026-06-15T10:00:00Z",
	)

	if metadata["source"] != "scheduler" {
		t.Fatalf("source = %q, want scheduler", metadata["source"])
	}
	if metadata[PeriodicMetadataNameKey] != "welcome-email" {
		t.Fatalf("periodic name = %q, want welcome-email", metadata[PeriodicMetadataNameKey])
	}
	if metadata[PeriodicMetadataDueAtKey] != "2026-06-15T10:00:00Z" {
		t.Fatalf("periodic due at = %q, want timestamp", metadata[PeriodicMetadataDueAtKey])
	}
}
