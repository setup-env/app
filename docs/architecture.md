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

## Application foundation

Command parsing remains a small standard-library dispatcher in `internal/cli`.
The live dashboard alone uses a terminal framework; static, JSON, diagnostic,
and module commands do not initialize terminal state.

Package responsibilities are:

- `internal/app`: orchestration between reusable services;
- `internal/catalog`: authoritative catalog model, sources, filtering, and validation;
- `internal/cli`: arguments, human output, JSON output, and actionable errors;
- `internal/compatibility`: narrow semantic-version minimum checks;
- `internal/config`: versioned, secret-free settings and OS-native location;
- `internal/directory`: structural development/organization/repository context;
- `internal/diagnostics`: safe tool, authentication-readiness, and write checks;
- `internal/dashboard`: live state, scheduling, histories, rates, terminal lifecycle, and responsive views;
- `internal/git`: Git repository and sanitized remote metadata inspection;
- `internal/manifest`: strict YAML parsing and manifest v1 validation;
- `internal/paths`: home and development-root resolution;
- `internal/platform`: OS, distribution, architecture, user, home, and shell;
- `internal/status`: deterministic human formatting and static status rendering;
- `internal/system`: snapshot contract, section collectors, warnings, and health;
- `internal/version`: build metadata.

The detection packages use injected functions or interfaces where machine state
would otherwise make tests brittle.

The manifest parser uses `github.com/goccy/go-yaml` v1.19.2. YAML is a core
public contract and is not supported by the Go standard library; this focused,
maintained dependency provides strict unknown-field decoding without adding
external module dependencies.

System metrics use `github.com/shirou/gopsutil/v4` v4.26.6. It is one focused,
actively maintained, cross-platform dependency for host, CPU, memory, and disk
information and avoids fragile operating-system command parsing. Network
interface collection uses Go's standard library. The dependency does not
require CGO for supported Windows, macOS, and Linux targets. Its transitive
modules are platform adapters selected by Go build constraints: Windows
WMI/OLE, Darwin `purego`, AIX `perfstat`, Plan 9 statistics, Unix system
configuration helpers, and `golang.org/x/sys`.

The live terminal uses `charm.land/bubbletea/v2` v2.0.8 for its cross-platform
event loop, alternate-screen lifecycle, keyboard input, resize events, and
defensive terminal restoration. `github.com/charmbracelet/x/term` v0.2.2
provides explicit stdin/stdout terminal detection. Both are pure Go and require
no CGO. Bubble Tea's transitive graph provides terminal input, ANSI width,
Unicode grapheme, cancellation, color-profile, terminfo, and Windows terminal
adapters. Setup Env does not add Bubbles or a broad application framework.

## System snapshot

`internal/system.Snapshot` is the presentation-independent, point-in-time data
contract for the CLI, future live dashboard, and eventual desktop application.
It includes schema version 1, an RFC 3339 timestamp when encoded as JSON,
explicit byte and percentage fields, development diagnostics, structured
warnings, and overall health.

Individual section collectors accept a context and update a snapshot. The
composite collector runs them independently and returns partial data when a
section fails. Optional numbers use pointers: JSON `null` and human
`unavailable` are different from a legitimate zero. The only total failure is
a cancelled collection or a snapshot for which every section failed.

Collection and rendering are separate. `internal/status` owns static human
formatting; the JSON encoder serializes the typed model directly. Neither
collector emits ANSI terminal sequences. CPU utilization uses a 500 ms sample
to balance responsiveness and usefulness.

User-relevant filesystem filtering excludes known pseudo filesystems and
internal mount trees (`/dev`, `/proc`, `/run`, `/snap`, `/sys`,
`/System/Volumes`, and common container storage), removes duplicates, and keeps
Windows drive roots. Network collection reads only local interface metadata and
unicast addresses; it makes no public-IP or other external request.

The dashboard consumes this same model; collection rules never move into its
view layer.

## Live dashboard

`internal/dashboard.Model` owns presentation state: the latest snapshot, clock,
bounded CPU and memory histories, derived network rates, terminal dimensions,
help/pause state, and refresh lifecycle. Its pure layout functions render ASCII
structure for wide, compact, help, and minimum-size states. Dashboard styling
does not require color, so `NO_COLOR` is naturally respected.

The application gathers one full initial snapshot, including static platform
metadata and development diagnostics. Subsequent refreshes are non-overlapping:

- CPU, memory, local interface data, and network counters: every second;
- filesystem capacity and inventory: every five seconds;
- development-root and Git/GitHub diagnostics: every sixty seconds or on
  forced refresh;
- clock and resize query: every second.

CPU collection retains the 500 ms utilization window. Stable host, OS, user,
CPU-model, and tool metadata are merged forward rather than rediscovered each
second. A forced `r` refresh includes all sections but still runs as one
cancellable operation.

Per-interface cumulative receive/transmit counters are an additive snapshot
schema-v1 extension. Dashboard rates are derived only when the same interface
has two available, nondecreasing counters and positive monotonic elapsed time.
Counter reset, interface appearance/disappearance, and unavailable data yield
`unavailable`, never a negative or fabricated rate. Rates remain dashboard
state rather than snapshot fields because they describe a relationship between
samples.

Bubble Tea declaratively enters the alternate screen and restores raw mode,
cursor, and screen state during handled exits. The dashboard passes a child
context through refresh commands and cancels it when the program returns.
Resize messages recompute layout without discarding histories; an explicit
size request each second also covers platforms without `SIGWINCH`.

Before dashboard initialization, both stdin and stdout must be terminals and
`TERM` must not be `dumb`. No-argument non-interactive execution routes to the
existing static renderer. Explicit `dashboard` mode instead returns an
actionable error. Static and JSON output remain ANSI-free.
If automatic dashboard initialization fails, the no-argument path falls back
to the static renderer unless the application context was canceled.

This separation also preserves future desktop reuse: the desktop application
can consume snapshots and rates without importing terminal layout or lifecycle
behavior.

## Manifest and catalog authority

The versioned module manifest identifies a module, compatibility, supported
platforms, security declarations, and workflow metadata. The app catalog is
authoritative for listing, repository location, trust, and status. A module
manifest is authoritative for capabilities, platform support, compatibility,
workflows, and module metadata. Catalog governance wins for trust even if a
third-party file makes a conflicting claim; strict parsing currently rejects a
manifest `trust` field entirely.

Milestone 02 embeds the repository catalog in the binary. Source interfaces
leave room for future explicit local and verified cached catalogs. Proposed
future precedence is explicit local, verified cache, then embedded. Only the
embedded source is active now.

## Future workflow engine

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

Trust is assigned only through review of the app catalog and cannot be
self-declared by a module. Trust must be visible and must not silently escalate.
Future downloads will use
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

## Awesome ecosystem list

`setup-env/awesome-setup-env` is human-readable curation and may include
official modules, community experiments, third-party projects, guides, and
related tools. Its Markdown is never an execution or trust source. Inclusion in
the Awesome list does not grant catalog inclusion, trust, compatibility, or
installability. Future automation may generate parts of that list from the app
catalog, but synchronization is outside Milestone 02.

## Desktop application

A future desktop application will call the same Go engine and typed services as
the CLI. Business rules and module execution will not be duplicated in a GUI.
Human CLI output is a presentation concern; structured data models remain
reusable by automation and desktop clients.
