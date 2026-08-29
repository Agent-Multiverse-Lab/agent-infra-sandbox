# Bounded adversarial workloads

This directory will contain opt-in workloads that test resource and isolation
claims. Every workload must have a hard internal bound in addition to runtime
limits. Do not commit unbounded fork bombs, disk fillers, or memory allocators.

Initial workload categories:

- bounded process fan-out;
- bounded memory pressure;
- bounded output flooding;
- attempted host-path reads;
- attempted network and metadata access;
- child-process cleanup after timeout.
