# agent-infra-sandbox

Independent research on isolation runtimes and lifecycle infrastructure for AI
agents.

This repository is an experimental laboratory, not an extension of
`agent-stack` and not a production sandbox service. Its purpose is to compare
runtime designs with reproducible workloads, measurements, and threat models.

## Research questions

- How much latency and memory does each isolation boundary add?
- Which resource controls survive hostile or pathological workloads?
- How should workspaces, snapshots, and process trees be cleaned up?
- What are the security and operational trade-offs among containers, Linux
  namespaces, WebAssembly, gVisor, and microVMs?
- When does warm reuse improve performance without weakening isolation?

## Current baseline

The first executable backend is `process`. It intentionally provides **no
isolation** and exists only as a control group for measurements. Running it
requires an explicit acknowledgement flag.

```bash
go run ./cmd/labctl runtimes

go run ./cmd/labctl run \
  --runtime process \
  --allow-unsafe-process \
  --timeout 5s \
  -- python3 workloads/normal/hello.py
```

The command prints a JSON result containing startup time, execution time, exit
status, timeout state, and bounded stdout/stderr.

## Repository layout

| Path | Purpose |
| --- | --- |
| `cmd/labctl` | Local experiment CLI |
| `pkg/sandbox` | Runtime-neutral contracts |
| `internal/runtime` | Runtime implementations and prototypes |
| `internal/isolation` | cgroup, filesystem, and network experiments |
| `experiments` | Numbered, reproducible research experiments |
| `workloads` | Normal and bounded adversarial workloads |
| `docs` | Threat model, roadmap, and runtime comparisons |
| `results` | Local raw output policy and checked-in summaries |

## Runtime roadmap

1. `process`: unsafe control group.
2. `docker`: container baseline with explicit resource controls.
3. `namespace`: direct Linux namespace and cgroup v2 experiments.
4. `wasm`: capability-oriented WebAssembly experiments.
5. gVisor and microVM adapters only after the earlier baselines are measured.

## Research rules

- Every experiment starts with a hypothesis and records its environment.
- Security claims require a threat model and a test that can falsify them.
- Adversarial workloads must be bounded by default and require explicit opt-in.
- Raw measurements are machine-readable; conclusions live beside the
  experiment that produced them.

See [the research roadmap](docs/research-roadmap.md) and
[the threat model](docs/threat-model.md) before adding a new runtime.

## Development

Requires Go 1.24 or newer.

```bash
make check
```
