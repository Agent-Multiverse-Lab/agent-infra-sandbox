package process

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"

	"github.com/agent-multiverse-lab/agent-infra-sandbox/pkg/sandbox"
)

func TestLifecycleAndExec(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("baseline workload uses POSIX commands")
	}

	backend := New()
	instance, err := backend.Create(context.Background(), sandbox.Spec{ID: "test-1"})
	if err != nil {
		t.Fatal(err)
	}

	result, err := backend.Exec(context.Background(), instance, sandbox.Command{
		Args: []string{"sh", "-c", "printf hello"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.ExitCode != 0 || result.Stdout != "hello" {
		t.Fatalf("unexpected result: %+v", result)
	}

	if err := backend.Destroy(context.Background(), instance); err != nil {
		t.Fatal(err)
	}
	if _, err := backend.Stats(context.Background(), instance); !errors.Is(err, sandbox.ErrNotFound) {
		t.Fatalf("Stats after Destroy returned %v, want ErrNotFound", err)
	}
}

func TestOutputIsBounded(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("baseline workload uses POSIX commands")
	}

	backend := New()
	instance, err := backend.Create(context.Background(), sandbox.Spec{ID: "bounded"})
	if err != nil {
		t.Fatal(err)
	}

	result, err := backend.Exec(context.Background(), instance, sandbox.Command{
		Args:      []string{"sh", "-c", "printf 1234567890"},
		MaxOutput: 4,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Stdout != "1234" || !result.StdoutTruncated {
		t.Fatalf("unexpected bounded output result: %+v", result)
	}
}

func TestTimeout(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("baseline workload uses POSIX commands")
	}

	backend := New()
	instance, err := backend.Create(context.Background(), sandbox.Spec{ID: "timeout"})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	result, err := backend.Exec(ctx, instance, sandbox.Command{Args: []string{"sleep", "1"}})
	if err != nil {
		t.Fatal(err)
	}
	if !result.TimedOut || result.ExitCode != -1 {
		t.Fatalf("unexpected timeout result: %+v", result)
	}
}
