# Threat model

## Protected assets

- Host filesystem, processes, kernel, credentials, and network identity.
- Other sandbox instances and their workspaces.
- Control-plane availability and measurement integrity.

## Untrusted inputs

- Model-generated source code and shell commands.
- Uploaded files, archives, dependencies, and build scripts.
- Workloads that intentionally consume excessive resources or avoid cleanup.

## Initial attacker capabilities

Assume the workload can execute arbitrary native code inside the selected
runtime and can coordinate child processes. It may attempt to:

- exhaust CPU, memory, PIDs, file descriptors, disk, or output buffers;
- read host paths, environment variables, credentials, or neighboring data;
- scan local networks or reach cloud metadata services;
- leave background processes or persistent filesystem changes;
- exploit kernel, runtime, mount, or control-plane configuration mistakes.

## Out of scope for the first baseline

- Defending against a compromised host kernel.
- Confidential computing and hardware side-channel resistance.
- Production multi-region scheduling.

## Required interpretation

The `process` runtime has no isolation and is only a performance control group.
A runtime is not considered isolated merely because it starts a container. Each
claimed boundary must be tied to a negative test and a documented configuration.
