# Architecture

## Responsibility

`setup-env/app` is the runtime, universal entrypoint, CLI, and future shared Go
engine for Setup Env. It owns platform and directory discovery, application
configuration, module/catalog contracts, download and verification policy,
workflow orchestration, shared actions, output, and execution history.

Each setup domain is an independently versioned repository. Workstation,
server, Terraform, Helm, cloud, and network behavior must not be embedded in
the application. This separation lets modules release at their own pace while
the application maintains a stable compatibility contract.

## Milestone 01 boundaries

The initial CLI uses the Go standard library. Four commands do not justify a
framework dependency, and command dispatch remains isolated in `internal/cli`
so a library can be introduced later if real nesting, completion, or lifecycle
needs make it valuable.

Package responsibilities are:

- `internal/app`: orchestration between reusable services;
- `internal/cli`: arguments, human output, JSON output, and actionable errors;
- `internal/config`: versioned, secret-free settings and OS-native location;
- `internal/directory`: structural development/organization/repository context;
- `internal/diagnostics`: safe tool, authentication-readiness, and write checks;
- `internal/git`: Git repository and sanitized remote metadata inspection;
- `internal/paths`: home and development-root resolution;
- `internal/platform`: OS, distribution, architecture, user, home, and shell;
- `internal/version`: build metadata.

The detection packages use injected functions or interfaces where machine state
would otherwise make tests brittle.

## Future module and workflow engine

A future versioned module manifest will identify a module, compatibility,
workflows, supported platforms, inputs, dependencies, actions, permissions,
artifacts, and verification metadata. A catalog will point to immutable module
releases rather than executing an arbitrary default branch.

The workflow engine will plan actions before execution, validate inputs and
platform support, resolve dependencies, require declared permissions, support
dry-run, redact secrets, make cancellation safe, and record an auditable result.
This contract will be implemented only after the proposal is tested against
reference modules.

## Trust and supply chain

The catalog will distinguish:

- **official** modules owned and released by the Setup Env organization;
- **verified** third-party modules whose identity and release process have been
  reviewed;
- **community** modules that meet the manifest contract without endorsement.

Trust must be visible and must not silently escalate. Future downloads will use
immutable versions, checksums, signed release metadata where practical,
compatibility checks, and cache verification. Scripts and templates remain
untrusted input until verified by policy.

Secrets must never appear in configuration, logs, plans, error messages, or
execution history. Provider integrations will rely on existing authenticated
tools or secure credential stores and pass only the minimum required data.

## Provider discovery

Git and SSH configuration can reveal useful local capabilities but cannot
enumerate all remote access. Future organization and repository discovery will
use authenticated provider APIs—initially GitHub APIs through `gh`—without
reading or storing token contents.

## Desktop application

A future desktop application will call the same Go engine and typed services as
the CLI. Business rules and module execution will not be duplicated in a GUI.
Human CLI output is a presentation concern; structured data models remain
reusable by automation and desktop clients.
