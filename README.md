# Setup Env

Setup Env is a cross-platform Go application for understanding and preparing a
development environment. It currently provides a point-in-time system status,
environment diagnostics, and a versioned catalog of Setup Env modules on
Windows, Apple macOS, and Ubuntu Linux.

## Install

Release installers are planned but are not published yet. Build the current
application from source using the platform guide:

### Windows

[Install Setup Env on Windows](docs/install/windows.md)

### Apple macOS

[Install Setup Env on Apple macOS](docs/install/macos.md)

### Ubuntu Linux

[Install Setup Env on Ubuntu Linux](docs/install/ubuntu.md)

## System status

Running `setup-env` with no arguments displays one static system snapshot.
This is deliberately the default product experience; it does not refresh or
take over the terminal. Help remains available through `setup-env --help` and
`setup-env help`.

```sh
setup-env
setup-env status
setup-env status --json
```

The snapshot includes local date and time, host and operating-system identity,
uptime, CPU, memory, user-relevant filesystems, local network interfaces,
development-root context, Git and GitHub CLI readiness, diagnostics, and
structured collection warnings. Unsupported or inaccessible metrics are shown
as unavailable without suppressing the remaining snapshot.

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

[`catalog/modules.yaml`](catalog/modules.yaml) is the authoritative
machine-readable catalog embedded in the binary. It controls listing,
repository location, trust, and status. A module's `setup-env.yaml` controls
its capabilities, platforms, compatibility, workflows, and descriptive
metadata.

[`setup-env/awesome-setup-env`](https://github.com/setup-env/awesome-setup-env)
is a separate human-curated list. Inclusion there does not grant catalog trust
or installability, and the CLI never scrapes Markdown to discover or execute
modules.

## Directory convention

Setup Env resolves paths dynamically under the current user's home directory:

```text
~/dev/<organization>/<repository>
```

Examples are `~/dev/setup-env/app` and `~/dev/setup-env/workstation`. No
username, drive letter, or slash convention is hard-coded. See
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

The next iteration is a live terminal dashboard that reuses the snapshot
collectors. Cross-platform releases, module retrieval, and workflow execution
remain future work. See the [roadmap](docs/roadmap.md),
[architecture](docs/architecture.md), and
[Milestone 03 notes](docs/milestone-03.md).

Propose module-contract changes here and domain behavior in the relevant
module repository. Read [CONTRIBUTING.md](CONTRIBUTING.md), the
[module model](docs/module-model.md), and the
[module contribution process](docs/module-contributions.md).

## License

[MIT](LICENSE)
