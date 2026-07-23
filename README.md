# Setup Env

Setup Env is a universal, modular setup and scaffolding platform for Windows,
macOS, and Linux. This repository contains the native Go application, command
line interface, shared runtime contracts, platform detection, and application
roadmap.

The application is the universal entrypoint. Domain-specific behavior belongs
in independently versioned module repositories such as
[`setup-env/workstation`](https://github.com/setup-env/workstation),
[`setup-env/terraform`](https://github.com/setup-env/terraform), and
[`setup-env/helm`](https://github.com/setup-env/helm). Workstation setup is one
module; it is not the main product.

## Milestone 01

Milestone 01 establishes the application foundation:

- a small, dependency-free CLI;
- cross-platform platform and directory-context detection;
- safe Git and GitHub CLI readiness diagnostics;
- a versioned, secret-free configuration model;
- tested package boundaries and cross-platform CI.

The current commands are:

```text
setup-env
setup-env version
setup-env info [--json]
setup-env doctor [--json]
```

Running `setup-env` without arguments displays help. The `module`, `workflow`,
and `run` commands are planned and are not implemented.

## Directory convention

Setup Env uses a universal hierarchy rooted at the current user's home:

```text
~/dev/<organization>/<repository>
```

Paths are resolved dynamically with operating-system APIs. No username, drive
letter, or slash convention is hard-coded. For example:

```text
~/dev/setup-env/app
~/dev/setup-env/workstation
```

See [Directory conventions](docs/directory-conventions.md).

## Build and run

Go 1.26 is required. The module pins the Go 1.26.5 toolchain used by CI.

```sh
go build -o bin/setup-env ./cmd/setup-env
go run ./cmd/setup-env
go run ./cmd/setup-env info --json
go test ./...
```

On Windows, the build output is `bin/setup-env.exe`. The executable name is
`setup-env` on every platform.

## Configuration

No configuration file is required. If present, `config.json` is loaded from
the operating system's standard user configuration directory under
`setup-env/`. See [the example configuration](docs/config.example.json).
Configuration never contains credentials or access tokens.

## Planned capabilities

Future milestones will define and validate the module manifest and catalog,
download and verify module releases, execute workflows with dry-run and
permission controls, publish cross-platform installers, and reuse the Go engine
in a desktop application. Organization discovery will use authenticated
provider APIs such as GitHub's API through `gh`; local credentials alone cannot
reliably enumerate every organization a user may access.

See the [roadmap](docs/roadmap.md) and [architecture](docs/architecture.md).

## Modules and contributions

Propose module-contract changes in this repository. Propose or implement
domain behavior in that module's repository. A new module proposal should
describe its domain, workflows, supported operating systems, required
permissions, maintainers, and why it should be official, verified, or
community-maintained.

Read [CONTRIBUTING.md](CONTRIBUTING.md) and the provisional
[module model](docs/module-model.md) before contributing.

## License

[MIT](LICENSE)
