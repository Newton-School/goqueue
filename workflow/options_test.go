package workflow

import (
	"errors"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/task"
)

func TestDefaultCanvasConfigUsesProductionDefaults(t *testing.T) {
	config := defaultCanvasConfig()

	if config.defaultQueue != "default" {
		t.Fatalf("defaultQueue = %q, want default", config.defaultQueue)
	}
	if config.codec == nil {
		t.Fatal("codec should be set")
	}
	if config.now == nil {
		t.Fatal("now should be set")
	}
}

func TestDefaultCanvasClockAdvances(t *testing.T) {
	config := defaultCanvasConfig()
	first := config.now()
	time.Sleep(2 * time.Millisecond)
	second := config.now()

	if !second.After(first) {
		t.Fatalf("default canvas clock did not advance: first=%s second=%s", first.Format(time.RFC3339Nano), second.Format(time.RFC3339Nano))
	}
}

func TestCanvasOptionsApplyValues(t *testing.T) {
	config := defaultCanvasConfig()
	now := time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC)

	options := []CanvasOption{
		WithCanvasDefaultQueue("critical"),
		WithCanvasCodec(task.JSONPayloadCodec{}),
		WithCanvasNow(func() time.Time { return now }),
	}

	for _, option := range options {
		if err := option(&config); err != nil {
			t.Fatalf("option returned error: %v", err)
		}
	}

	if config.defaultQueue != "critical" {
		t.Fatalf("defaultQueue = %q, want critical", config.defaultQueue)
	}
	if got := config.now(); !got.Equal(now) {
		t.Fatalf("now = %v, want %v", got, now)
	}
}

func TestCanvasOptionsRejectInvalidValues(t *testing.T) {
	tests := []struct {
		name   string
		option CanvasOption
	}{
		{name: "invalid queue", option: WithCanvasDefaultQueue("invalid queue")},
		{name: "nil codec", option: WithCanvasCodec(nil)},
		{name: "nil now", option: WithCanvasNow(nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := defaultCanvasConfig()
			if err := tt.option(&config); !errors.Is(err, ErrInvalidWorkflow) {
				t.Fatalf("option error = %v, want ErrInvalidWorkflow", err)
			}
		})
	}
}
