# Development

## Requirements

- Go 1.26; `go.mod` pins toolchain 1.26.5.
- Git for repository-aware context and normal contribution workflows.
- GitHub CLI only for authenticated GitHub operations; core local commands do
  not require it.

Go 1.26.5 is the current supported Go 1.26 patch selected for this milestone.
Keep the toolchain directive and CI in sync when updating it.

## Commands

```sh
go fmt ./...
go vet ./...
go test ./...
go test -race ./...
go build -o bin/setup-env ./cmd/setup-env
```

Race testing requires a supported C toolchain and runs on Linux CI. The regular
matrix validates Windows, macOS, and Linux.

Build metadata can be injected with linker flags:

```sh
go build -ldflags "-X github.com/setup-env/app/internal/version.Version=v0.1.0 -X github.com/setup-env/app/internal/version.Commit=<sha> -X github.com/setup-env/app/internal/version.Date=<date>" -o bin/setup-env ./cmd/setup-env
```

## Design rules

- Use `path/filepath` for local filesystem paths.
- Inject machine state in tests; do not depend on a contributor's username,
  home, installed tools, credentials, or organization.
- Never print credential command output.
- Keep module-specific behavior in its module repository.
- Do not describe proposed commands as implemented.
