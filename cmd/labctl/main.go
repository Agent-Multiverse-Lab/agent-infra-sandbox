// Command labctl runs local sandbox experiments.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	processruntime "github.com/agent-multiverse-lab/agent-infra-sandbox/internal/runtime/process"
	"github.com/agent-multiverse-lab/agent-infra-sandbox/pkg/sandbox"
)

const reportSchemaVersion = 1

type runReport struct {
	SchemaVersion int            `json:"schema_version"`
	Runtime       string         `json:"runtime"`
	SandboxID     string         `json:"sandbox_id"`
	StartedAt     time.Time      `json:"started_at"`
	StartupMillis int64          `json:"startup_ms"`
	Result        sandbox.Result `json:"result"`
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "labctl:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		printUsage()
		return errors.New("a subcommand is required")
	}

	switch args[0] {
	case "runtimes":
		fmt.Println("process\tunsafe host-process control group")
		return nil
	case "run":
		return runCommand(args[1:])
	case "help", "-h", "--help":
		printUsage()
		return nil
	default:
		return fmt.Errorf("unknown subcommand %q", args[0])
	}
}

func runCommand(args []string) error {
	flags := flag.NewFlagSet("run", flag.ContinueOnError)
	runtimeName := flags.String("runtime", "process", "runtime backend")
	id := flags.String("id", "", "sandbox instance id")
	timeout := flags.Duration("timeout", 30*time.Second, "execution timeout")
	maxOutput := flags.Int("max-output-bytes", 1<<20, "maximum captured bytes per output stream")
	allowUnsafeProcess := flags.Bool(
		"allow-unsafe-process",
		false,
		"acknowledge that the process runtime has no isolation",
	)
	if err := flags.Parse(args); err != nil {
		return err
	}

	commandArgs := flags.Args()
	if len(commandArgs) > 0 && commandArgs[0] == "--" {
		commandArgs = commandArgs[1:]
	}
	if len(commandArgs) == 0 {
		return errors.New("a command is required after --")
	}
	if *timeout <= 0 {
		return errors.New("timeout must be positive")
	}
	if *maxOutput <= 0 {
		return errors.New("max-output-bytes must be positive")
	}

	backend, err := selectRuntime(*runtimeName, *allowUnsafeProcess)
	if err != nil {
		return err
	}
	if *id == "" {
		*id = fmt.Sprintf("run-%d", time.Now().UTC().UnixNano())
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	startedAt := time.Now().UTC()
	createStarted := time.Now()
	instance, err := backend.Create(ctx, sandbox.Spec{ID: *id})
	if err != nil {
		return fmt.Errorf("create instance: %w", err)
	}
	startupMillis := time.Since(createStarted).Milliseconds()
	defer func() {
		if destroyErr := backend.Destroy(context.Background(), instance); destroyErr != nil {
			fmt.Fprintln(os.Stderr, "labctl: destroy instance:", destroyErr)
		}
	}()

	result, err := backend.Exec(ctx, instance, sandbox.Command{
		Args:      commandArgs,
		MaxOutput: *maxOutput,
	})
	if err != nil {
		return fmt.Errorf("execute command: %w", err)
	}

	report := runReport{
		SchemaVersion: reportSchemaVersion,
		Runtime:       backend.Name(),
		SandboxID:     instance.ID,
		StartedAt:     startedAt,
		StartupMillis: startupMillis,
		Result:        result,
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		return fmt.Errorf("encode report: %w", err)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("workload exited with status %d", result.ExitCode)
	}
	return nil
}

func selectRuntime(name string, allowUnsafeProcess bool) (sandbox.Runtime, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "process":
		if !allowUnsafeProcess {
			return nil, errors.New("process runtime is not isolated; pass --allow-unsafe-process to use the control group")
		}
		return processruntime.New(), nil
	default:
		return nil, fmt.Errorf("runtime %q is not implemented", name)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  labctl runtimes")
	fmt.Fprintln(os.Stderr, "  labctl run [flags] -- <command> [args...]")
}
