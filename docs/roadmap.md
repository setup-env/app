# Roadmap

## Completed

### Milestone 01 — Application Foundation

Established the Go module, initial CLI, environment and directory detection,
configuration, diagnostics, tests, CI, architecture, and governance.

### Milestone 02 — Module Manifest and Official Catalog

Defined manifest v1, embedded the authoritative catalog, validated trust,
status and compatibility, exposed local discovery commands, and documented the
contribution flow.

### Milestone 03 — System Snapshot and Status Output

Added the reusable partial-failure snapshot, deterministic static output,
schema-versioned JSON, system metrics, and installation landing documentation.

### Milestone 04 — Live Terminal Dashboard

Reused the snapshot collectors in a responsive live terminal application with
bounded CPU/memory histories, network throughput, keyboard controls, help,
resize handling, non-interactive fallback, and safe terminal restoration. See
[Milestone 04](milestone-04.md).

### Milestone 05 — Cross-Platform Releases and Bootstrap Installers

Added reproducible Windows, macOS, and Linux archives, release metadata,
SHA-256 checksums, user-scoped PowerShell and POSIX installers, safe upgrade and
uninstall behavior, offline tests, tag-only publication automation, and release
operations documentation. Public `v0.1.0` publication remains an explicit
approval step. See [Milestone 05](milestone-05.md).

## Next

### Milestone 06 — Module Download, Cache, and Verification

Resolve immutable module releases, download and cache artifacts, verify
checksums and release metadata, and implement safe update behavior.

### Milestone 07 — Workstation Reference Module

Integrate `setup-env/workstation` as a reference module without moving domain
logic into the application, using it to validate the module contract.

## Later work

Workflow execution, a native desktop application, and opt-in remote management
remain later initiatives. Remote management depends on mature local security,
identity, permission, and audit models.
