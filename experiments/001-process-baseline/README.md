# Experiment 001: host-process baseline

## Hypothesis

Direct host execution establishes the minimum lifecycle and execution overhead
for the lab harness. It provides no isolation and must not be compared as a
security boundary.

## Procedure

1. Run the normal hello workload at least 100 times.
2. Record startup and execution duration from the JSON output.
3. Repeat with bounded stdout and a forced timeout.
4. Record the host environment before interpreting results.

## Metrics

- Create/startup latency.
- Execution latency.
- Timeout detection latency.
- Output truncation correctness.
- Cleanup errors.
