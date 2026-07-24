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

Changes to system collection must preserve partial results, context
cancellation, deterministic ordering, and the distinction between unavailable
and zero. Add fixture-based tests rather than asserting the contributor's
hardware, username, paths, credentials, network, or CPU count. Human rendering
and JSON serialization must remain separate from metric collection.

Dashboard changes must keep terminal behavior inside `internal/dashboard`,
preserve non-interactive static fallback, bound all histories, avoid
overlapping collectors, and honor context cancellation. Layout, rates,
controls, and refresh scheduling should be tested with injected state rather
than manual terminal timing. Verify that static and JSON output remain
ANSI-free after every dashboard change.

Commits should be focused and written in the imperative mood. Pull requests
should explain the problem, approach, user impact, validation, permissions, and
security considerations.

Module proposals should include maintainers, trust level, target platforms,
workflows, dependencies, permission requirements, and test strategy. The
manifest v1 contract is defined in [docs/module-model.md](docs/module-model.md).

## Registering a module

Module registration uses a pull request to this repository. Contributors add
one valid catalog entry and provide a public repository with a valid
`setup-env.yaml`, documentation, and tests where applicable. Automated
validation must pass and a maintainer must approve the assigned trust level.

Contributors may request `community` trust. They cannot assign themselves
`verified` or `official`; those classifications require separate maintainer
governance review. Follow
[docs/module-contributions.md](docs/module-contributions.md).

The recommended Git flow is:

```text
feature branch → pull request → validation → squash merge to main
```

The active `Protect main` repository ruleset requires a pull request, all
required CI jobs, resolved conversations, squash merging, and linear history.
It blocks ordinary direct pushes, force pushes, and deletion. Because the
repository currently has one direct maintainer, the required approval count is
zero; requiring that maintainer to approve their own pull request would be
unsatisfiable. New commits still dismiss any recorded approvals.

Repository administrators have a pull-request-only bypass for emergency
maintenance. It preserves a PR and audit trail, does not permit direct pushes
to `main`, and must not be used for routine feature or dependency work.

Release and installer changes must preserve deterministic archive names and
contents, SHA-256 verification before execution, user-scoped defaults, rollback
on replacement failure, and scoped uninstall. Installer tests must use
temporary locations and offline fixtures; pull-request workflows must never
receive release write permission.

Do not publish tags or GitHub Releases from a contribution branch. Follow
[`docs/releasing.md`](docs/releasing.md) only from clean, synchronized `main`
after explicit release authorization.
