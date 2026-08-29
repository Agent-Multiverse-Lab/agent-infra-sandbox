# Experiments

Experiments are numbered and append-only. Each experiment directory should
contain:

- `README.md`: hypothesis, variables, procedure, and interpretation rules;
- `config/`: runtime configuration when applicable;
- `scripts/`: reproducible orchestration;
- `results/`: small checked-in summaries, not large raw traces;
- `environment.json`: host and runtime metadata captured by the runner.

Measured facts and interpretation must be written separately. Failed or
inconclusive experiments remain useful history and should not be deleted.
