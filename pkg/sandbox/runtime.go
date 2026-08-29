// Package sandbox defines runtime-neutral contracts used by experiments.
package sandbox

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"
)

var (
	// ErrNotFound indicates that an instance is unknown to a runtime.
	ErrNotFound = errors.New("sandbox instance not found")
	// ErrStatsUnavailable indicates that a runtime cannot report live statistics.
	ErrStatsUnavailable = errors.New("sandbox statistics unavailable")
)

var validID = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,63}$`)

// Limits describes requested resource ceilings. A runtime must document which
// fields it actually enforces.
type Limits struct {
	MemoryBytes int64 `json:"memory_bytes,omitempty"`
	CPUQuota    int64 `json:"cpu_quota_micros,omitempty"`
	PIDs        int64 `json:"pids,omitempty"`
	DiskBytes   int64 `json:"disk_bytes,omitempty"`
}

// Spec describes a sandbox instance independently of a concrete runtime.
type Spec struct {
	ID         string            `json:"id"`
	WorkingDir string            `json:"working_dir,omitempty"`
	Env        map[string]string `json:"env,omitempty"`
	Limits     Limits            `json:"limits,omitempty"`
}

// Command is one execution request inside an instance.
type Command struct {
	Args       []string          `json:"args"`
	WorkingDir string            `json:"working_dir,omitempty"`
	Env        map[string]string `json:"env,omitempty"`
	MaxOutput  int               `json:"max_output_bytes,omitempty"`
}

// Instance is the runtime-neutral handle returned by Create.
type Instance struct {
	ID      string `json:"id"`
	Runtime string `json:"runtime"`
}

// Result captures one command execution without treating a non-zero program
// exit code as a runtime transport error.
type Result struct {
	ExitCode        int           `json:"exit_code"`
	Duration        time.Duration `json:"-"`
	DurationMillis  int64         `json:"duration_ms"`
	TimedOut        bool          `json:"timed_out"`
	Stdout          string        `json:"stdout"`
	Stderr          string        `json:"stderr"`
	StdoutTruncated bool          `json:"stdout_truncated"`
	StderrTruncated bool          `json:"stderr_truncated"`
}

// Stats is a deliberately small common measurement surface. Runtime-specific
// measurements belong in experiment output, not this interface.
type Stats struct {
	MemoryBytes int64 `json:"memory_bytes,omitempty"`
	PIDs        int64 `json:"pids,omitempty"`
}

// Runtime is the shared lifecycle contract used by the research harness.
type Runtime interface {
	Name() string
	Create(context.Context, Spec) (Instance, error)
	Exec(context.Context, Instance, Command) (Result, error)
	Stats(context.Context, Instance) (Stats, error)
	Destroy(context.Context, Instance) error
}

// ValidateID rejects identifiers that are awkward or unsafe in paths, labels,
// and runtime object names.
func ValidateID(id string) error {
	if !validID.MatchString(id) {
		return fmt.Errorf("invalid sandbox id %q", id)
	}
	return nil
}
