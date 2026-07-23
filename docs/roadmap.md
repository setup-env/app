# Roadmap

## Milestone 01 — Application Foundation (complete)

Establish the Go module, initial CLI, environment and directory detection,
configuration model, diagnostics, tests, CI, architecture, and governance.

## Milestone 02 — Module Manifest and Official Catalog (complete)

Define manifest v1, embed the authoritative official catalog, validate trust,
status and compatibility, expose local module discovery commands, and document
the contribution flow.

## Milestone 03 — Module Download, Cache, and Verification

Resolve immutable releases, download and cache artifacts, verify checksums and
release metadata, and implement safe update behavior.

## Milestone 04 — Workflow Execution Engine

Plan and execute typed workflows with platform gates, dependency resolution,
dry-run, permissions, secret redaction, cancellation, and execution history.

## Milestone 05 — Workstation Reference Module Integration

Integrate `setup-env/workstation` as a reference module without moving its
domain logic into the application. Use it to validate the module contract.

## Milestone 06 — Cross-Platform Releases and Installer

Produce verified Windows, macOS, and Linux artifacts, installation methods,
upgrade paths, release notes, and rollback guidance.

## Future desktop application

Build a native desktop experience that reuses the Go engine and typed services.

## Future remote management platform

Design opt-in fleet orchestration only after the local security, permission,
identity, audit, and workflow models have matured. Remote management is not in
the current application scope.
