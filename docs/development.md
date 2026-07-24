# Development

## Requirements

- Go 1.26; `go.mod` pins toolchain 1.26.5.
- `github.com/goccy/go-yaml` v1.19.2, downloaded automatically by Go modules.
- `github.com/shirou/gopsutil/v4` v4.26.6 for cross-platform host, CPU, memory,
  and filesystem metrics. Platform-specific transitive modules are selected by
  build constraints and do not add a runtime service.
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
go run ./cmd/setup-env module validate-catalog
go run ./cmd/setup-env module validate examples/setup-env.yaml
go run ./cmd/setup-env status
go run ./cmd/setup-env status --json
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
- Keep catalog entries, categories, and tags deterministically sorted.
- Keep parsing, schema, semantics, and compatibility validation distinct.
- Never accept module trust from a manifest; trust belongs to the app catalog.
- Keep system collection independent from both human and JSON rendering.
- Encode unavailable numeric metrics as JSON `null`, never as a fabricated
  zero; use raw bytes and 0–100 percentages.
- Treat individual collector errors as structured warnings when a meaningful
  partial snapshot remains.
- Keep status output deterministic and free of ANSI escape sequences.
- Keep module-specific behavior in its module repository.
- Do not describe proposed commands as implemented.

## Status JSON contract

`setup-env status --json` emits schema version `1`. Timestamps use RFC 3339,
bytes and percentages are numeric, optional numeric values are `null` when
unavailable, and warnings contain a section, message, and severity. New
backward-compatible fields may be added within schema version 1; renaming,
removing, or changing the meaning or units of fields requires a new schema
version.

The status command should exit successfully when one metric or section cannot
be read but a meaningful snapshot exists. Invalid arguments, cancellation, or
failure of every section return an error.

## Supported status metrics

| Area | Windows | macOS | Ubuntu/Linux |
| --- | --- | --- | --- |
| Host, OS, architecture, user, uptime | Yes | Yes | Yes |
| Distribution name/version | OS platform data | OS platform data | `/etc/os-release` plus platform data |
| CPU model, counts, sampled utilization | Where exposed by OS APIs | Where exposed by OS APIs | Where exposed by OS APIs |
| Physical memory and utilization | Yes | Yes | Yes |
| User-relevant filesystem capacity | Drive roots | Relevant mounted volumes | Relevant mounted filesystems |
| Local interfaces, MAC, IPv4, IPv6 | Yes | Yes | Yes |
| Git, GitHub CLI, development root | Yes | Yes | Yes |

The underlying operating system, permissions, virtual machine, or container may
make an individual metric unavailable. Public IP, Wi-Fi secrets, process data,
disk I/O rates, and network throughput are intentionally not collected.
