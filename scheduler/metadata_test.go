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

func TestPeriodicDispatchMetadataOverridesReservedKeys(t *testing.T) {
	metadata := periodicDispatchMetadata(
		map[string]string{
			PeriodicMetadataNameKey:  "user-value",
			PeriodicMetadataDueAtKey: "user-time",
		},
		"welcome-email",
		"2026-06-15T10:00:00Z",
	)

	if metadata[PeriodicMetadataNameKey] != "welcome-email" {
		t.Fatalf("periodic name = %q, want scheduler value", metadata[PeriodicMetadataNameKey])
	}
	if metadata[PeriodicMetadataDueAtKey] != "2026-06-15T10:00:00Z" {
		t.Fatalf("periodic due at = %q, want scheduler value", metadata[PeriodicMetadataDueAtKey])
	}
}
