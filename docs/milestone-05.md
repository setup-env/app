# Milestone 05 — Cross-Platform Releases and Bootstrap Installers

## Scope

Milestone 05 adds reproducible application archives, release metadata, SHA-256
checksums, user-scoped installers, upgrade/rollback/uninstall behavior, offline
tests, tag-only publication automation, and installation/security/release
documentation. Module retrieval/execution, package managers, signing,
notarization, self-update, GUI, telemetry, and remote behavior remain excluded.

## GitHub tracking

- [Parent #26](https://github.com/setup-env/app/issues/26)
- [Versioning and artifact standards #27](https://github.com/setup-env/app/issues/27)
- [Reproducible builds #28](https://github.com/setup-env/app/issues/28)
- [Checksums and verification #29](https://github.com/setup-env/app/issues/29)
- [PowerShell installer #30](https://github.com/setup-env/app/issues/30)
- [POSIX installer #31](https://github.com/setup-env/app/issues/31)
- [PATH, upgrade, and uninstall #32](https://github.com/setup-env/app/issues/32)
- [Installer security #33](https://github.com/setup-env/app/issues/33)
- [Release and installer tests #34](https://github.com/setup-env/app/issues/34)
- [GitHub Release workflow #35](https://github.com/setup-env/app/issues/35)
- [README installation landing #36](https://github.com/setup-env/app/issues/36)
- [First release candidate #37](https://github.com/setup-env/app/issues/37)
- [Release operations and rollback #38](https://github.com/setup-env/app/issues/38)

The GitHub App returned `403 Resource not accessible by integration` when
creating the parent issue. The authenticated GitHub CLI fallback created all
issues, assigned milestone 3, and linked #27–#38 as native sub-issues of #26.

## Release design

The repository-native `cmd/release` tool uses only the Go standard library and
the existing toolchain. It builds with `CGO_ENABLED=0`, `-trimpath`,
`-buildvcs=false`, and an empty build ID. Linker flags inject semantic version,
commit, commit timestamp, and dirty state without source mutation.

Archives use deterministic entry order, modes, owners, and commit timestamps.
Each contains:

- `setup-env` or `setup-env.exe`;
- `LICENSE`;
- release `README.md`;
- `THIRD_PARTY_NOTICES.md`.

Artifact matrix for `v0.1.0`:

| Operating system | Architecture | Asset |
| --- | --- | --- |
| Windows | amd64 | `setup-env_0.1.0_windows_amd64.zip` |
| Windows | arm64 | `setup-env_0.1.0_windows_arm64.zip` |
| macOS | amd64 | `setup-env_0.1.0_darwin_amd64.tar.gz` |
| macOS | arm64 | `setup-env_0.1.0_darwin_arm64.tar.gz` |
| Linux | amd64 | `setup-env_0.1.0_linux_amd64.tar.gz` |
| Linux | arm64 | `setup-env_0.1.0_linux_arm64.tar.gz` |

`checksums.txt` contains one sorted SHA-256 entry per archive. The builder
self-verifies checksum coverage, hashes, safe paths, and exact archive contents.

Dependency license files in the resolved Go module cache were reviewed before
packaging. Runtime dependencies use permissive MIT, BSD, Apache-2.0, and related
licenses; no new release-tool dependency or copyleft runtime dependency was
introduced. `THIRD_PARTY_NOTICES.md` summarizes the direct runtime licenses,
and `go.mod`/`go.sum` remain the version authority.

## Installer behavior

Both installers support latest or explicit versions, custom user-writable
directories, verified idempotent installation, upgrade, rollback, uninstall,
purge, offline fixtures, clear failures, and temporary cleanup.

Windows defaults to `%LOCALAPPDATA%\Programs\setup-env\bin`, updates only user
PATH without duplicates, and records whether that entry is installer-managed.
macOS/Linux default to `~/.local/bin` and print Bash, Zsh, and Fish instructions
instead of modifying an unknown profile.

Installers download only from official GitHub Releases, never require a token,
and never execute an archive before its checksum matches. Explicit versions use
stable URLs; latest uses the documented releases API and reports rate-limit,
proxy, and network failures.

## Testing strategy

- Go unit tests cover versions, targets, names, deterministic archives, minimal
  contents, checksum parsing/mismatch, and safe output cleanup.
- PowerShell 5.1-compatible tests cover version/architecture mapping, names,
  checksums, PATH deduplication/removal, custom paths, and missing releases.
- POSIX tests cover equivalent mappings, checksum behavior, PATH detection, and
  missing releases.
- Native CI integration builds a local fixture, installs, executes `version`,
  upgrades, uninstalls, preserves unrelated files, rejects corruption, and
  checks temporary cleanup on Windows, Ubuntu, and macOS.
- Pull-request snapshot CI builds all six archives, validates the matrix and
  metadata, and rebuilds Linux amd64 to compare reproducible hashes.

No test modifies persistent runner PATH or real user installation directories.
No installer test requires public GitHub access.

## Workflow security and publication

PR CI has `contents: read`. The separate `v*` workflow has only
`contents: write`, pins official actions by immutable SHA, validates the exact
tag and release uniqueness, and publishes only after every validation succeeds.
Artifact attestations are deferred; checksums remain mandatory.

The public `v0.1.0` tag and GitHub Release are intentionally pending an explicit
publication approval after implementation merges into synchronized `main`.

## Signing limitations and future hardening

Windows binaries are initially unsigned. macOS binaries are unsigned and
unnotarized. Checksums are not independent signing. Future hardening may add
artifact attestations, protected release tags, Authenticode, Apple signing and
notarization, and package-manager channels once maintainable trust processes
exist.
