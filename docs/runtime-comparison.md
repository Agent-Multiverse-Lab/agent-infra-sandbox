# Runtime comparison framework

This table records hypotheses to test, not conclusions.

| Runtime | Expected startup | Native compatibility | Expected boundary | Primary question |
| --- | --- | --- | --- | --- |
| Host process | Lowest | Full | None | What is the measurement floor? |
| OCI container | Low | Full | Shared kernel | Which controls are reliable under hostile workloads? |
| Raw namespaces | Low | Full | Shared kernel | What overhead comes from the primitives themselves? |
| WebAssembly | Low | Limited | Capability runtime | Is reduced compatibility worth a narrower attack surface? |
| gVisor | Medium | High | Userspace kernel | What security gain justifies syscall overhead? |
| MicroVM | Highest | Full | Virtualized kernel | Can snapshots make the stronger boundary practical? |

Every comparison should report the same workload revision, host environment,
resource limits, warm/cold state, sample count, and percentile method.
