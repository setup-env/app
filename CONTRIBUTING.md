# Contributing

Thank you for helping build Setup Env.

Application runtime, CLI, platform, catalog, and workflow-specification changes
belong in this repository. Domain-specific setup logic belongs in its module
repository. Open an issue before a large architectural change so its contract
and milestone can be agreed first.

## Development workflow

1. Install the pinned Go toolchain.
2. Create a focused branch.
3. Add or update tests with behavior changes.
4. Run `go fmt ./...`, `go vet ./...`, `go test ./...`, and
   `go build ./cmd/setup-env`.
5. Update public documentation when behavior or contracts change.

Commits should be focused and written in the imperative mood. Pull requests
should explain the problem, approach, user impact, validation, permissions, and
security considerations.

Module proposals should include maintainers, trust level, target platforms,
workflows, dependencies, permission requirements, and test strategy. The
module specification remains provisional until the module-engine milestones
validate it.
