# Roadmap

## Completed

### Milestone 01 — Application Foundation

Established the Go module, initial CLI, environment and directory detection,
configuration, diagnostics, tests, CI, architecture, and governance.

### Milestone 02 — Module Manifest and Official Catalog

Defined manifest v1, embedded the authoritative catalog, validated trust,
status and compatibility, exposed local discovery commands, and documented the
contribution flow.

## Current

### Milestone 03 — System Snapshot and Status Output

Provide a reusable, partial-failure-tolerant system snapshot with deterministic
human output and a versioned JSON contract. Add `setup-env status` and make a
static snapshot the no-argument experience. See [Milestone 03](milestone-03.md).

## Next

### Milestone 04 — Live Terminal Dashboard

Reuse the snapshot collectors in a lightweight, `htop`-inspired terminal
display with live refresh, trends, keyboard controls, resizing, and clean
terminal restoration. This milestone will not add process management.

### Milestone 05 — Cross-Platform Releases and Bootstrap Installers

Publish verified Windows, macOS, and Linux artifacts with supported bootstrap
installation, upgrade, release-note, and rollback paths.

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
