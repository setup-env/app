# Milestone 04 — Live Terminal Dashboard

## Outcome

Milestone 04 turns the Milestone 03 snapshot into a responsive live terminal
application without changing the static status or JSON contract. An interactive
`setup-env` launch opens the dashboard; redirected, piped, `TERM=dumb`, or
otherwise non-interactive no-argument execution prints one ANSI-free snapshot.
`setup-env dashboard` explicitly requires an interactive terminal.
Automatic no-argument startup also falls back to a static snapshot when
terminal initialization fails; explicit dashboard mode reports the failure.

This milestone does not implement process management, package or module
installation, workflow execution, remote monitoring, desktop UI, background
services, telemetry, public-IP lookup, or release packaging.

## GitHub tracking

- [Milestone 04](https://github.com/setup-env/app/milestone/2)
- [Parent issue #13](https://github.com/setup-env/app/issues/13)
- [Architecture and refresh loop #15](https://github.com/setup-env/app/issues/15)
- [Terminal lifecycle #16](https://github.com/setup-env/app/issues/16)
- [CPU and memory panels #17](https://github.com/setup-env/app/issues/17)
- [Filesystem and development panels #18](https://github.com/setup-env/app/issues/18)
- [Network counters and rates #19](https://github.com/setup-env/app/issues/19)
- [Keyboard controls and help #20](https://github.com/setup-env/app/issues/20)
- [Resize handling #21](https://github.com/setup-env/app/issues/21)
- [Non-interactive fallback #22](https://github.com/setup-env/app/issues/22)
- [Tests and CI #23](https://github.com/setup-env/app/issues/23)
- [Documentation and examples #24](https://github.com/setup-env/app/issues/24)

The children are linked to #13 through GitHub's native sub-issue relationship
and through explicit parent links in their bodies.

## Architecture

`internal/dashboard` contains terminal-specific behavior:

- `source.go`: reusable initial/refresh source contract and cadence defaults;
- `model.go`: Bubble Tea state machine, clock, scheduling, controls, and
  non-overlap guard;
- `history.go`: fixed-capacity oldest-to-newest metric histories;
- `rate.go`: interface-aware monotonic network-rate calculation;
- `format.go`: byte-rate, ASCII usage bar, history graph, and truncation;
- `layout.go`: pure wide, compact, help, and minimum-size rendering;
- `terminal.go`: stdin/stdout terminal suitability detection;
- `run.go`: cancellable Bubble Tea program lifecycle.

`internal/app.LiveCollector` orchestrates existing `internal/system` section
collectors. `internal/system` remains presentation-independent, and
`internal/status` remains the deterministic ANSI-free static renderer.

## Refresh cadence

| Data | Cadence | Notes |
| --- | --- | --- |
| Clock and terminal-size query | 1 second | Size polling supplements resize events on Windows |
| CPU, memory, interfaces, counters | 1 second | CPU includes its existing 500 ms utilization sample |
| Filesystem inventory and capacity | 5 seconds | Uses Milestone 03 filtering and ordering |
| Development and diagnostics | 60 seconds | Avoids repeated Git/GitHub subprocess checks |
| Stable host, OS, user, CPU model | Initial snapshot | Merged forward between dynamic samples |

`r` forces all sections once. `p` or Space pauses metric collection while the
clock and input remain responsive. Only one refresh may be active.

## Histories and rate calculation

CPU and memory histories retain at most sixty samples by default, ordered from
oldest to newest. Rendering clips them to the available width.

The snapshot schema remains version 1 and adds nullable cumulative receive,
transmit, and packet counters to each network interface. The dashboard compares
the same case-insensitive interface name across consecutive samples:

```text
rate = (current counter - previous counter) / monotonic elapsed seconds
```

A missing counter, new interface, removed interface, reset/decreasing counter,
or nonpositive elapsed time produces `unavailable`. Rates never become negative
and require no external network access.

## Terminal lifecycle and accessibility

Bubble Tea v2.0.8 owns alternate-screen entry, input mode, cursor handling,
resize events, panic cleanup, and terminal restoration. Dashboard refreshes use
a child context cancelled when the program returns. `q`, Ctrl+C, context
cancellation, and handled errors all return through the framework cleanup path.

The layout uses ASCII borders, bars, and graphs and is useful without color.
No color is required, so `NO_COLOR` is inherently respected and red/green is
never the only distinction. At very small dimensions the dashboard shows a
minimum-size message instead of attempting an unsafe layout.

Controls:

```text
q / Ctrl+C  quit and restore the terminal
r           force a full refresh
p / Space   pause or resume metrics
?           toggle help
```

## Dependency decision

`charm.land/bubbletea/v2` v2.0.8 was selected because it is a current stable,
maintained, no-CGO terminal framework with Windows, macOS, and Linux support,
declarative alternate-screen state, keyboard input, resizing, and defensive
cleanup. `github.com/charmbracelet/x/term` v0.2.2 supplies focused
platform-independent terminal detection.

The tradeoff is Bubble Tea's focused terminal transitive graph: terminal input,
width/grapheme handling, cancellation, terminfo, color-profile, and Windows
adapters. No Bubbles component suite, desktop framework, or application
framework was added. Static and JSON commands do not initialize Bubble Tea.

## Validation policy

Fixture tests cover bounded histories, CPU/memory graphs, bars, rates and
counter resets, interface churn, zero elapsed time, layouts, truncation, help,
controls, pause/refresh, partial errors, cancellation, alternate-screen
declaration, and non-TTY routing. CI runs tests and native status/fallback
smokes on Windows, macOS, and Ubuntu, Linux race tests, and six no-CGO
cross-builds.

A pseudo-terminal smoke is intentionally omitted. Adding another PTY/emulator
dependency for one timing-sensitive integration test would not improve the
deterministic model, lifecycle, and cross-platform compilation coverage enough
to justify its platform risk. Manual dashboard verification instructions are
in [Development](development.md).

## Limitations

- Network counters depend on operating-system support and interface naming.
- A newly seen or reset interface needs another sample before a rate appears.
- Disk I/O throughput is deferred; filesystem capacity is shown.
- Windows resize events are supplemented by a one-second size query.
- Very small terminals show a resize message; compact layouts intentionally
  hide additional rows.
- Transient high CPU or memory does not independently mark health unhealthy.

## Next milestone

Milestone 05 will publish verified cross-platform releases and bootstrap
installation paths. This milestone does not create installers, packages, or
GitHub Releases.
