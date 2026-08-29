# Repository working rules

## Scope

This repository studies sandbox infrastructure independently. Do not introduce
application-specific agent orchestration, conversation IDs, LangGraph state, or
dependencies on `agent-stack`.

## Research discipline

- State the hypothesis before implementing an experiment.
- Record runtime, kernel, CPU, memory, and configuration with benchmark data.
- Treat the `process` runtime as an unsafe control group, never as a sandbox.
- Do not claim an isolation property unless a test attempts to violate it.
- Keep adversarial workloads bounded and opt-in so accidental execution cannot
  exhaust a developer machine.

## Engineering

- Keep the public contract in `pkg/sandbox` small and runtime-neutral.
- Put runtime-specific code under `internal/runtime/<name>`.
- Prefer the standard library until an experiment demonstrates the need for a
  dependency.
- Run `make check` before submitting changes.
- Number experiments monotonically and do not rewrite historical results.
