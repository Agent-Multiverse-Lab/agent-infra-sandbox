# Experiment 002: Docker baseline

## Status

Planned. Implementation begins only after Experiment 001 produces stable,
versioned output.

## Hypothesis

An OCI container adds measurable cold-start overhead but can enforce useful CPU,
memory, PID, filesystem, and network controls for local Agent workloads.

## Required comparisons

- Cold container versus pre-created container.
- Network disabled versus default bridge.
- Read-only root filesystem versus writable root filesystem.
- Enforced resource limits versus attempted resource exhaustion.
- Clean exit, timeout, runtime crash, and orphaned child process cleanup.
