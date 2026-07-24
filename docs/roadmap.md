# Roadmap

## Completed

### Milestone 01 — Application Foundation

Established the Go module, initial CLI, environment and directory detection,
application configuration, diagnostics, tests, CI, architecture, and
governance.

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
operations documentation. The first public release is
[`v0.0.1`](https://github.com/setup-env/app/releases/tag/v0.0.1). See
[Milestone 05](milestone-05.md).

## Planned

### Milestone 06 — Portable Configuration and Desired-State Planning

Define a portable, declarative, versioned desired-state model without executing
module workflows or mutating the machine. The milestone covers:

- configuration schema v1 and forward-compatible migration rules;
- a strict boundary between portable desired state and local observed state;
- desktop, headless, server, platform, and capability conditions;
- safe `config` initialization, validation, inspection, and normalization;
- import, export, backup, restore, and sharing workflows;
- declarative module selection and module configuration blocks;
- secret references that never embed secret values;
- desired-state diff and plan models with deterministic human and JSON output;
- cross-platform tests and documentation.

Portable configuration must be suitable for USB or cloud backup and for private
or public Git storage when its contents are safe to share. This milestone does
not install packages, download modules, execute workflows, or mutate the local
machine beyond explicitly requested configuration-file creation.

### Milestone 07 — Module Download, Cache, and Verification

Resolve catalog modules and immutable releases, download archives safely,
verify SHA-256 and trust metadata, extract defensively, and maintain a
verifiable local cache that supports offline reuse. It adds cache inspection
and removal commands, compatibility enforcement, deterministic fixtures, and
security documentation.

This milestone does not execute module workflows, install packages, run
arbitrary scripts, or apply the Workstation module.

### Milestone 08 — Apply Engine and Local State

Introduce the engine boundary that compares portable desired state with
observed local state, produces an approval-ready plan, applies module actions
idempotently, and records local results. It covers locking, concurrent-run
protection, partial failures, retry and compensation contracts, history,
auditing, drift detection, secret resolution boundaries, and `apply`, `diff`,
and `state` commands.

No module-specific Workstation behavior belongs in the application engine.

### Milestone 09 — Workstation Reference Module

Integrate `setup-env/workstation` as the first real module using the manifest,
download, planning, apply, and state contracts established by the application.
Detailed package and platform implementation belongs in the Workstation
repository; the app tracks only integration-level work.

Integration must verify Windows, macOS, Ubuntu desktop, and headless
Ubuntu/Linux profiles. Capability filtering must prevent GUI-only software,
including Visual Studio Code, Docker Desktop, browsers, and desktop terminals,
from being selected for headless nodes such as `mgmt`.

### Milestone 10 — Configuration Restore and Reproducibility Validation

Prove that a saved configuration can reproduce supported environments safely:
new Windows, macOS, and Ubuntu desktops; a new headless Ubuntu development
node; restoration after an OS reinstall; configuration sharing; repeated
idempotent application; unavailable-package handling; platform substitutions;
and external secret references or prompts.

## Sequencing note

The earlier roadmap placed module download at Milestone 06 and Workstation
integration at Milestone 07. Before either began, they were deliberately moved
to Milestones 07 and 09. Portable desired-state configuration now precedes
download, and the apply engine precedes the first real module. Milestone 10
then validates end-to-end restoration. This numbering is authoritative.

## Later work

A native desktop application and opt-in remote management remain later
initiatives. Remote management depends on mature local security, identity,
permission, and audit models.
