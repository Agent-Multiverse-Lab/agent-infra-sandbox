# Research roadmap

## Phase 0: measurement harness

- Stabilize the runtime-neutral lifecycle contract.
- Produce versioned JSON output with bounded stdout and stderr.
- Record host and runtime metadata alongside benchmark results.
- Establish the unsafe host-process control group.

## Phase 1: OCI container baseline

- Compare cold creation with warm reuse.
- Verify CPU, memory, PID, and filesystem limits.
- Measure cleanup after timeout, crash, and orphaned child processes.
- Compare default, disabled, and allow-listed network modes.

## Phase 2: Linux primitives

- Build minimal mount, PID, user, UTS, IPC, and network namespace experiments.
- Apply cgroup v2 limits without a full container manager.
- Evaluate seccomp, capabilities, and `no_new_privileges` independently.

## Phase 3: alternative boundaries

- Compare WebAssembly capability models with native-code compatibility.
- Evaluate gVisor only after the OCI baseline is reproducible.
- Evaluate microVMs after startup and snapshot requirements are explicit.

## Phase 4: lifecycle strategies

- Warm pools and snapshot restore.
- Workspace persistence without runtime persistence.
- Lease, heartbeat, reclamation, and crash recovery.
- Multi-tenant fairness and noisy-neighbor measurements.
