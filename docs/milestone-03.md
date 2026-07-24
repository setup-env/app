# Milestone 03 — System Snapshot and Status Output

## Purpose

Milestone 03 adds the reusable point-in-time snapshot that underpins the future
live dashboard and desktop application. It includes deterministic human output,
a schema-versioned JSON contract, partial-failure behavior, installation
landing documentation, tests, and cross-platform CI. It does not include live
refresh, a terminal UI, process management, module retrieval, or workflow
execution.

The no-argument command deliberately displays the same static snapshot as
`setup-env status`. `setup-env help` and `setup-env --help` retain explicit
help access. This is a migration from the earlier no-argument help behavior and
is covered by CLI tests.

## GitHub tracking status and access blocker

The public repository was inspected before implementation on 2026-07-24. It
had no existing milestones or issues, and the latest `main` CI run observed was
successful (run `30054347272`). GitHub CLI authentication was rechecked with
`gh auth status`; the configured credential was invalid. The connected GitHub
integration is installed only for the unrelated `smuts-tech` organization, not
`setup-env`; repository write attempts returned HTTP 403, and the integration
does not expose milestone creation. A normal `gh auth login` requires the
repository owner's interactive browser authorization; no token was requested,
printed, copied, or stored.

Consequently, no milestone or issue listed below was created through GitHub at
implementation time. Branch, pull-request, protection, and merge status must be
reported separately after final validation. These proposals are complete and
ready to create once an authenticated repository owner restores write access.

## Proposed Milestone 03

**Title:** Milestone 03 — System Snapshot and Status Output

**Description:** Build a reusable cross-platform, point-in-time system
snapshot for Windows, Apple macOS, and Ubuntu Linux. Expose deterministic human
output and schema-versioned JSON through `setup-env status`; make the same
static snapshot the deliberate no-argument experience; preserve help through
explicit flags; reuse environment diagnostics; tolerate partial collection
failures; validate all supported operating systems in CI; and publish honest
source-install documentation. Excludes live terminal refresh, process
management, release installers, module downloading/caching/installation, and
workflow execution.

**Completion criteria:** Snapshot and renderers are separate and reusable;
status human and JSON output work; unavailable values differ from zero;
warnings are structured; existing commands remain compatible; tests, vet,
formatting, Linux race testing, and six cross-builds pass; CI smokes and parses
status output on Windows, macOS, and Ubuntu; required documentation is current.

## Proposed parent issue

**Title:** Implement Milestone 03 — System Snapshot and Status Output

**Milestone:** Milestone 03 — System Snapshot and Status Output

**Body:** Implement the reusable point-in-time system snapshot and static
status experience described by this document. Coordinate the eight child
issues below. Keep collection independent from rendering, retain machine-ready
JSON, degrade gracefully on unavailable metrics, and preserve Milestones 01
and 02.

**Acceptance criteria:**

- [ ] `setup-env status` and `setup-env status --json` work cross-platform.
- [ ] No arguments show one static snapshot; explicit help remains available.
- [ ] Host, OS, user, time, uptime, CPU, memory, filesystem, network,
      development, diagnostics, warnings, and health are represented.
- [ ] JSON schema version 1 uses raw bytes, numeric percentages, RFC 3339 time,
      nullable unavailable values, and no ANSI sequences.
- [ ] Existing commands and catalog behavior remain compatible.
- [ ] Windows, macOS, Ubuntu, race, smoke, and cross-build validation pass.
- [ ] README, install guides, architecture, development, and roadmap are
      updated.

## Proposed child issues

### 1. Define cross-platform system snapshot model

**Body:** Add a presentation-independent schema-versioned snapshot with
explicit units, nullable optional numeric fields, structured warnings,
diagnostic health, section interfaces, context support, deterministic ordering,
and partial collection.

**Acceptance criteria:** JSON serialization is fixture-tested; legitimate zero
differs from unavailable; collection can run as a whole or by section; all
section failure and cancellation produce errors while partial failure does not
discard successful data.

### 2. Implement date, time, host, OS, user, and uptime collection

**Body:** Reuse the accepted platform detector and `/etc/os-release` handling,
add hostname, OS display/version/build, distribution, architecture, kernel,
current user, home, shell, uptime, local timezone name, and UTC offset without
external network access.

**Acceptance criteria:** JSON timestamp is RFC 3339; human time is local;
fixtures cover merged platform/host information; missing fields yield warnings
or unavailable values without suppressing other sections.

### 3. Implement CPU and memory collection

**Body:** Collect CPU model, physical/logical counts, total utilization using a
documented short sample, and physical-memory total/available/used/utilization
using a focused cross-platform library.

**Acceptance criteria:** Sampling stays between 250 and 1000 ms (target 500
ms); calculations and unavailable values are tested; no process monitoring is
added; high utilization alone does not mark the environment unhealthy.

### 4. Implement filesystem and disk-usage collection

**Body:** Collect mount or drive, filesystem type, and total/available/used
bytes and utilization. Filter pseudo, temporary internal, container, device,
and duplicate mounts while keeping relevant Windows drive roots and Unix
mounts.

**Acceptance criteria:** Filtering and calculations are fixture-tested;
ordering is deterministic; a failed mount does not remove successful mounts;
disk I/O throughput is deferred.

### 5. Implement network interface and address collection

**Body:** Collect local interface name, operational state, safe MAC address,
unicast IPv4/IPv6 addresses, and loopback classification using local APIs only.

**Acceptance criteria:** Addresses and interfaces are deterministic and
fixture-tested; no public-IP request, Wi-Fi secret, credential, or live
throughput is collected.

### 6. Add status CLI command and JSON output

**Body:** Add `status` and `status --json`, choose the static snapshot as the
no-argument experience, keep explicit help, and add reusable deterministic IEC
byte, percentage, duration, timestamp, offset, and unavailable formatting.

**Acceptance criteria:** Human output is readable and contains no terminal UI;
JSON schema version 1 is machine-oriented and ANSI-free; invalid arguments and
total collection failure are non-zero; partial snapshots remain successful.

### 7. Add cross-platform tests and CI validation

**Body:** Add injected fixture tests and CI status smoke tests for Windows,
macOS, and Ubuntu. Parse JSON structurally without asserting ephemeral runner
hardware. Retain Linux race coverage and cross-build six OS/architecture
targets locally.

**Acceptance criteria:** Format, diff check, vet, all tests, Linux race, six
cross-builds, command smoke tests, JSON parsing, no-ANSI checks, and catalog
validation pass without relying on a developer workstation or internet access.

### 8. Update README installation and usage documentation

**Body:** Turn the README into a Windows/macOS/Ubuntu installation landing
page; add honest build-from-source and future-release paths; document static
status, partial failures, JSON, dependencies, metrics, roadmap order, and
Milestone 04 boundaries.

**Acceptance criteria:** No unpublished installer is presented as available;
no personal home path appears; all required documents and platform-specific
source build/PATH guidance exist.

## Proposed Milestone 04

**Title:** Milestone 04 — Live Terminal Dashboard

**Description:** Reuse Milestone 03 snapshot and metric collectors in a
lightweight live terminal dashboard. Add refresh scheduling, trends,
interaction, resizing, clean shutdown, terminal restoration, and graceful
platform degradation. Do not duplicate collection logic or expand into process
management.

## Proposed Milestone 04 planning issue

**Title:** Implement live terminal dashboard

**Milestone:** Milestone 04 — Live Terminal Dashboard

**Body:** Design and implement an opt-in, lightweight `htop`-inspired live
terminal presentation backed by repeated Milestone 03 snapshot collection.
Display live refresh of date/time and health, CPU and memory trends, and disk
and network activity where supported. Include development-root and diagnostic
status, documented keyboard controls, terminal resize handling, context-driven
shutdown, restoration of terminal state on normal exit and errors, and graceful
degradation for unsupported or inaccessible metrics.

**Acceptance criteria:**

- [ ] Reuses the system snapshot and collectors instead of parsing CLI output.
- [ ] Provides a bounded/configurable refresh interval and CPU/memory trends.
- [ ] Adds disk/network activity only where counters are reliable.
- [ ] Displays date/time plus development and diagnostic health.
- [ ] Documents and tests keyboard controls and terminal resizing.
- [ ] Restores cursor, screen, and terminal modes on every shutdown path.
- [ ] Handles cancellation and unsupported metrics gracefully.
- [ ] Static `status` human and JSON behavior remains stable.

**Explicit exclusions:** process lists or killing, mouse support, desktop GUI,
background service, telemetry, public-IP discovery, remote monitoring, module
installation/execution, and release packaging.

## Snapshot and output notes

The snapshot JSON contract is version `1`. Bytes and percentages stay numeric;
optional values serialize as `null`; timestamps are RFC 3339; warnings expose a
section, message, and severity. Human output uses IEC byte units and the local
timezone. A 500 ms CPU sample is the only intentional collection delay.

Filesystem filtering excludes known pseudo types and internal mount prefixes,
deduplicates mount points, and retains relevant Windows drive roots. Network
collection is local-only. Collection warnings contribute to `warning` health;
failed diagnostics contribute to `unhealthy`; resource utilization does not
independently determine health.

See [Architecture](architecture.md) and [Development](development.md) for the
platform matrix, dependency impact, and contract evolution rules.
