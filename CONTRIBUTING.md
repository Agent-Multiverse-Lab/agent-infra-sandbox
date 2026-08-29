# Contributing

Changes should improve either reproducibility, measurement quality, isolation
coverage, or runtime understanding.

## Before opening a pull request

1. Describe the research question or defect.
2. Add or update a numbered experiment when behavior is being evaluated.
3. Keep potentially destructive workloads bounded and explicitly enabled.
4. Run `make check`.
5. Separate measured facts from hypotheses in documentation.

Implementation changes should include tests. Benchmark changes should record
the environment closely enough for another machine to reproduce the run.
