package main

import (
	"os"
	"strings"
	"testing"
)

func TestResolveRedisURLUsesEnvIfPresent(t *testing.T) {
	const envKey = "GOQUEUE_REDIS_URL"
	t.Setenv(envKey, "redis://env-host:6379/0")

	if got := resolveRedisURL(); got != "redis://env-host:6379/0" {
		t.Fatalf("resolveRedisURL() = %q, want env-host URL", got)
	}
}

func TestResolveRedisURLFallsBackByDefault(t *testing.T) {
	t.Setenv("GOQUEUE_REDIS_URL", "")
	if got := resolveRedisURL(); got != "redis://localhost:6379/0" {
		t.Fatalf("resolveRedisURL() = %q, want localhost URL", got)
	}
}

func TestResolveNamespaceUsesEnvIfPresent(t *testing.T) {
	t.Setenv("GOQUEUE_NAMESPACE", "batch-worker")
	if got := resolveNamespace(defaultNamespace); got != "batch-worker" {
		t.Fatalf("resolveNamespace() = %q, want batch-worker", got)
	}
}

func TestResolveQueueUsesEnvIfPresent(t *testing.T) {
	t.Setenv("GOQUEUE_DEFAULT_QUEUE", "critical")
	if got := resolveQueue(defaultTaskQueue); got != "critical" {
		t.Fatalf("resolveQueue() = %q, want critical", got)
	}
}

func TestPrintJSONProducesMachineFormat(t *testing.T) {
	var captured struct {
		Status string `json:"status"`
	}
	captured.Status = "ok"
	output := captureStdout(t, func() {
		printJSON(captured)
	})

	if output == "" {
		t.Fatal("printJSON output is empty")
	}
	if expected := `"status": "ok"`; !strings.Contains(output, expected) {
		t.Fatalf("printJSON output = %q, want %q", output, expected)
	}
}

func TestSplitStreamIDsParsesCommaSeparatedValues(t *testing.T) {
	ids, err := splitStreamIDs("1-0,2-0, 3-0")
	if err != nil {
		t.Fatalf("splitStreamIDs returned error: %v", err)
	}

	got := strings.Join(ids, ",")
	if got != "1-0,2-0,3-0" {
		t.Fatalf("splitStreamIDs got=%q want=%q", got, "1-0,2-0,3-0")
	}
}

func TestSplitStreamIDsReturnsNilForEmptyValue(t *testing.T) {
	ids, err := splitStreamIDs("")
	if err != nil {
		t.Fatalf("splitStreamIDs returned error: %v", err)
	}
	if ids != nil {
		t.Fatalf("splitStreamIDs result = %v, want nil", ids)
	}
}

func TestPrintUsageIncludesControlCommand(t *testing.T) {
	output := captureStdout(t, func() {
		printUsage()
	})
	if !strings.Contains(output, "goqueue control") {
		t.Fatalf("printUsage output = %q, missing control command", output)
	}
}

func TestPrintControlUsageIncludesReplayCommand(t *testing.T) {
	output := captureStdout(t, func() {
		printControlUsage()
	})
	if !strings.Contains(output, "replay-dead-letter") {
		t.Fatalf("printControlUsage output = %q, missing replay-dead-letter", output)
	}
}

func TestCtxReturnsContext(t *testing.T) {
	ctx := ctx()
	if ctx == nil {
		t.Fatal("ctx() returned nil")
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() returned error: %v", err)
	}

	os.Stdout = writer
	fn()
	writer.Close()

	outBytes := make([]byte, 0, 256)
	buf := make([]byte, 512)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			outBytes = append(outBytes, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	os.Stdout = orig
	return string(outBytes)
}
