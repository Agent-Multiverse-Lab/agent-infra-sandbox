// Package process implements the deliberately unsafe host-process baseline.
package process

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/agent-multiverse-lab/agent-infra-sandbox/pkg/sandbox"
)

const defaultMaxOutput = 1 << 20

// Runtime executes commands directly on the host. It is a measurement control
// group and provides no security boundary.
type Runtime struct {
	mu        sync.RWMutex
	instances map[string]sandbox.Spec
}

// New constructs an empty process runtime.
func New() *Runtime {
	return &Runtime{instances: make(map[string]sandbox.Spec)}
}

func (r *Runtime) Name() string { return "process" }

func (r *Runtime) Create(_ context.Context, spec sandbox.Spec) (sandbox.Instance, error) {
	if err := sandbox.ValidateID(spec.ID); err != nil {
		return sandbox.Instance{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.instances[spec.ID] = cloneSpec(spec)

	return sandbox.Instance{ID: spec.ID, Runtime: r.Name()}, nil
}

func (r *Runtime) Exec(
	ctx context.Context,
	instance sandbox.Instance,
	command sandbox.Command,
) (sandbox.Result, error) {
	spec, err := r.lookup(instance)
	if err != nil {
		return sandbox.Result{}, err
	}
	if len(command.Args) == 0 {
		return sandbox.Result{}, errors.New("command args are required")
	}

	maxOutput := command.MaxOutput
	if maxOutput <= 0 {
		maxOutput = defaultMaxOutput
	}
	stdout := newCappedBuffer(maxOutput)
	stderr := newCappedBuffer(maxOutput)

	cmd := exec.CommandContext(ctx, command.Args[0], command.Args[1:]...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = mergeEnv(os.Environ(), spec.Env, command.Env)
	cmd.Dir = command.WorkingDir
	if cmd.Dir == "" {
		cmd.Dir = spec.WorkingDir
	}

	started := time.Now()
	runErr := cmd.Run()
	duration := time.Since(started)
	result := sandbox.Result{
		ExitCode:        0,
		Duration:        duration,
		DurationMillis:  duration.Milliseconds(),
		TimedOut:        errors.Is(ctx.Err(), context.DeadlineExceeded),
		Stdout:          stdout.String(),
		Stderr:          stderr.String(),
		StdoutTruncated: stdout.Truncated(),
		StderrTruncated: stderr.Truncated(),
	}

	if runErr == nil {
		return result, nil
	}

	var exitErr *exec.ExitError
	if errors.As(runErr, &exitErr) {
		result.ExitCode = exitErr.ExitCode()
		if result.TimedOut {
			result.ExitCode = -1
		}
		return result, nil
	}
	return sandbox.Result{}, fmt.Errorf("start host process: %w", runErr)
}

func (r *Runtime) Stats(
	_ context.Context,
	instance sandbox.Instance,
) (sandbox.Stats, error) {
	if _, err := r.lookup(instance); err != nil {
		return sandbox.Stats{}, err
	}
	return sandbox.Stats{}, sandbox.ErrStatsUnavailable
}

func (r *Runtime) Destroy(_ context.Context, instance sandbox.Instance) error {
	if instance.Runtime != r.Name() {
		return fmt.Errorf("instance runtime %q does not match %q", instance.Runtime, r.Name())
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.instances[instance.ID]; !ok {
		return sandbox.ErrNotFound
	}
	delete(r.instances, instance.ID)
	return nil
}

func (r *Runtime) lookup(instance sandbox.Instance) (sandbox.Spec, error) {
	if instance.Runtime != r.Name() {
		return sandbox.Spec{}, fmt.Errorf("instance runtime %q does not match %q", instance.Runtime, r.Name())
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	spec, ok := r.instances[instance.ID]
	if !ok {
		return sandbox.Spec{}, sandbox.ErrNotFound
	}
	return cloneSpec(spec), nil
}

func cloneSpec(spec sandbox.Spec) sandbox.Spec {
	cloned := spec
	cloned.Env = make(map[string]string, len(spec.Env))
	for key, value := range spec.Env {
		cloned.Env[key] = value
	}
	return cloned
}

type envSlice []string

func (values envSlice) apply(target map[string]string) {
	for _, item := range values {
		for index := 0; index < len(item); index++ {
			if item[index] == '=' {
				target[item[:index]] = item[index+1:]
				break
			}
		}
	}
}

type envMap map[string]string

func (values envMap) apply(target map[string]string) {
	for key, value := range values {
		target[key] = value
	}
}

// Go does not implicitly convert []string or map[string]string to an interface
// implemented by named aliases, so keep conversion at the call site concise.
func mergeEnv(base []string, overlays ...map[string]string) []string {
	values := make(map[string]string)
	envSlice(base).apply(values)
	for _, overlay := range overlays {
		envMap(overlay).apply(values)
	}

	result := make([]string, 0, len(values))
	for key, value := range values {
		result = append(result, key+"="+value)
	}
	return result
}

type cappedBuffer struct {
	data      []byte
	limit     int
	truncated bool
}

func newCappedBuffer(limit int) *cappedBuffer {
	return &cappedBuffer{data: make([]byte, 0, min(limit, 4096)), limit: limit}
}

func (b *cappedBuffer) Write(value []byte) (int, error) {
	written := len(value)
	remaining := b.limit - len(b.data)
	if remaining <= 0 {
		b.truncated = true
		return written, nil
	}
	if len(value) > remaining {
		value = value[:remaining]
		b.truncated = true
	}
	b.data = append(b.data, value...)
	return written, nil
}

func (b *cappedBuffer) String() string { return string(b.data) }

func (b *cappedBuffer) Truncated() bool { return b.truncated }

var _ io.Writer = (*cappedBuffer)(nil)
var _ sandbox.Runtime = (*Runtime)(nil)
