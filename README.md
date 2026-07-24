# Setup Env

Setup Env is a cross-platform terminal application for understanding and
preparing a development environment. It provides a live system dashboard,
static and machine-readable status, environment diagnostics, and a versioned
catalog of Setup Env modules on Windows, Apple macOS, and Ubuntu Linux.

## Install

Release installers are planned but are not published yet. Build the current
application from source using the platform guide:

### Windows

[Install Setup Env on Windows](docs/install/windows.md)

### Apple macOS

[Install Setup Env on Apple macOS](docs/install/macos.md)

### Ubuntu Linux

[Install Setup Env on Ubuntu Linux](docs/install/ubuntu.md)

## Dashboard and status

Launch `setup-env` in an interactive terminal to open the live dashboard:

```sh
setup-env
```

When input or output is redirected or piped, the same no-argument command
prints one static, ANSI-free snapshot instead. Automation never enters the
dashboard unexpectedly.

```sh
setup-env > status.txt
setup-env | less
setup-env status
setup-env status --json
setup-env dashboard
```

`setup-env dashboard` explicitly requires an interactive terminal and otherwise
returns an actionable error. `setup-env status` is always static, while
`setup-env status --json` emits schema-versioned machine data.
If automatic dashboard initialization fails, the no-argument path restores the
terminal and falls back to one static snapshot.

Dashboard controls:

```text
q / Ctrl+C  quit
r           refresh all metrics now
p / Space   pause or resume
?           toggle help
```

Representative wide-terminal layout:

```text
+ Setup Env -------------------------------------------------------------+
| example-host | Ubuntu 26.04 | amd64 | uptime 03:12:18 | 2026-07-24 ... |
+ CPU -------------------------------+ + Memory -------------------------+
| 17.3%  physical 4  logical 8       | | 3.2 GiB / 7.3 GiB  43.8%       |
| ...::--==++                         | | ...:::---===                    |
+ Filesystems -----------------------------------------------------------+
| /          28.4 GiB / 118.0 GiB  24.1% [####----------------]         |
+ Network ---------------------------------------------------------------+
| eth0  10.0.0.9  down 12.4 KiB/s  up 2.1 KiB/s                        |
+ Development and health -----------------------------------------------+
| root ~/dev | Git available | GitHub CLI available | health healthy    |
+------------------------------------------------------------------------+
 q quit | r refresh | p pause | ? help | live
```

The dashboard uses ASCII structure and remains useful without color. CPU,
memory, and network refresh every second; filesystems every five seconds; and
development diagnostics every sixty seconds. Individual metric failures
remain visible as warnings without stopping the interface.

Help remains available through `setup-env help` and `setup-env --help`.

Other commands are:

```text
setup-env version
setup-env info [--json]
setup-env doctor [--json]
setup-env module list [--json] [--trust <level>] [--status <status>] [--category <category>]
setup-env module info <module> [--json]
setup-env module validate <path> [--json]
setup-env module validate-catalog
```

Catalog discovery and local manifest validation are implemented. Module
downloading, caching, installation, updates, and workflow execution are not.
A listed `planned` module is not runnable.

## Module catalog

[`catalog/modules.yaml`](catalog/modules.yaml) is the authoritative embedded
catalog. It controls listing, repository location, trust, and status. A
module's `setup-env.yaml` controls capabilities, platforms, compatibility,
workflows, and metadata.

[`setup-env/awesome-setup-env`](https://github.com/setup-env/awesome-setup-env)
is a separate human-curated list. Inclusion there does not grant catalog trust
or installability.

## Directory convention

Setup Env resolves paths dynamically under the current user's home directory:

```text
~/dev/<organization>/<repository>
```

No username, drive letter, or slash convention is hard-coded. See
[Directory conventions](docs/directory-conventions.md).

## Build and run

Go 1.26 is required; `go.mod` pins Go 1.26.5.

```sh
go build -o bin/setup-env ./cmd/setup-env
go run ./cmd/setup-env status
go run ./cmd/setup-env status --json
go test ./...
```

On Windows, use `bin/setup-env.exe`. The command name is `setup-env` on every
platform.

## Configuration

No configuration file is required. If present, `config.json` is loaded from
the operating system's standard user configuration directory under
`setup-env/`. See [the example configuration](docs/config.example.json).
Configuration never contains credentials or access tokens.

## Roadmap and contributions

Cross-platform releases are next; module retrieval and workflow execution
remain future work. See the [roadmap](docs/roadmap.md),
[architecture](docs/architecture.md), and
[Milestone 04 notes](docs/milestone-04.md).

Read [CONTRIBUTING.md](CONTRIBUTING.md), the
[module model](docs/module-model.md), and the
[module contribution process](docs/module-contributions.md).

## License

[MIT](LICENSE)
