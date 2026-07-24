# Setup Env

Setup Env is a cross-platform terminal application for inspecting and preparing
a development environment. It provides a live system dashboard, static and
machine-readable status, diagnostics, and a versioned local module catalog.

## Install

Official binaries are distributed only through
[GitHub Releases](https://github.com/setup-env/app/releases). Before running a
bootstrap command, confirm that the requested release exists and review the
downloaded installer. The prepared first version is `v0.1.0`; it is not
available until that tag and release appear on the releases page.

Installers use current-user directories, verify the archive with SHA-256 before
extraction, and do not require Git, Go, `gh`, administrator, or root access.

### Windows

PowerShell 5.1 or 7 on amd64 or arm64:

```powershell
$version = "v0.1.0"
Invoke-WebRequest https://raw.githubusercontent.com/setup-env/app/main/install.ps1 -OutFile install.ps1
Get-Content .\install.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File .\install.ps1 -Version $version
Remove-Item .\install.ps1
```

`-ExecutionPolicy Bypass` applies only to that reviewed installer process; the
script never changes system execution-policy settings. Detailed and manual
instructions: [Windows installation](docs/install/windows.md).

### Apple macOS

POSIX shell on Intel or Apple Silicon:

```sh
version=v0.1.0
curl -fL https://raw.githubusercontent.com/setup-env/app/main/install.sh -o install.sh
less install.sh
sh install.sh --version "$version"
rm install.sh
```

Detailed and manual instructions: [macOS installation](docs/install/macos.md).

### Ubuntu Linux

POSIX shell on amd64 or arm64:

```sh
version=v0.1.0
curl -fL https://raw.githubusercontent.com/setup-env/app/main/install.sh -o install.sh
less install.sh
sh install.sh --version "$version"
rm install.sh
```

Detailed and manual instructions: [Ubuntu installation](docs/install/ubuntu.md).

### Launch

Open a new terminal if the installer added or reported a PATH entry, then run:

```sh
setup-env version
setup-env
```

Useful non-interactive commands:

```sh
setup-env status
setup-env status --json
setup-env doctor
setup-env module list
```

### Upgrade and uninstall

Download and review the current installer again, then:

```powershell
# Windows
.\install.ps1 -Upgrade
.\install.ps1 -Uninstall -Yes
```

```sh
# macOS and Ubuntu
sh install.sh --upgrade
sh install.sh --uninstall --yes
```

Uninstall removes only installer-owned binary and metadata files. Configuration
and cache remain unless the explicit `-Purge`/`--purge` option is combined with
uninstall.

### Security and signing

Installers download only release archives and `checksums.txt` from the official
`setup-env/app` GitHub Release. A checksum mismatch stops installation before
the binary is executed. SHA-256 protects integrity relative to the checksum
published in the same release; it is not equivalent to independent code
signing.

Windows binaries are initially unsigned. macOS binaries are initially unsigned
and unnotarized, so SmartScreen or Gatekeeper may warn. The installers do not
disable those controls or remove quarantine attributes. See [Security](SECURITY.md).

Homebrew, WinGet, Chocolatey, Scoop, APT repositories, Snap, Flatpak, MSI, PKG,
and DMG distribution are not available.

## Dashboard and status

Launching `setup-env` in a suitable interactive terminal opens the live
dashboard. Redirected, piped, `TERM=dumb`, or failed automatic initialization
prints one static ANSI-free snapshot instead.

```text
q / Ctrl+C  quit
r           refresh all metrics now
p / Space   pause or resume
?           toggle help
```

CPU, memory, and network refresh every second; filesystems every five seconds;
development diagnostics every sixty seconds. Individual metric failures appear
as warnings without stopping the dashboard.

`setup-env status` is always static. `setup-env status --json` emits the stable
schema-versioned snapshot intended for automation.

## Module catalog

[`catalog/modules.yaml`](catalog/modules.yaml) is the authoritative embedded
catalog. Catalog discovery and local manifest validation are implemented.
Module downloading, caching, installation, updates, and execution are not.

## Build from source

Source development requires Git and Go 1.26; `go.mod` pins toolchain 1.26.5.

```sh
git clone https://github.com/setup-env/app.git
cd app
go test ./...
go build -o bin/setup-env ./cmd/setup-env
```

Use `bin/setup-env.exe` on Windows. Source-development details are in
[Development](docs/development.md).

## Release verification

Every release contains six native archives and `checksums.txt`. Each archive
contains only the native executable, project license, release README, and
third-party notices. Manual verification is documented in each platform guide.

Release preparation and rollback procedures are in
[Release operations](docs/releasing.md). Changes are curated in
[CHANGELOG.md](CHANGELOG.md).

## Configuration and directories

No configuration is required. If present, `config.json` uses the operating
system’s standard user configuration directory. Development paths are resolved
dynamically as:

```text
~/dev/<organization>/<repository>
```

No username, drive letter, or slash convention is hard-coded. See
[Directory conventions](docs/directory-conventions.md).

## Contributing and security

See [CONTRIBUTING.md](CONTRIBUTING.md) for development rules. Report
vulnerabilities using the private process in [SECURITY.md](SECURITY.md); do not
publish credentials or exploit details in public issues.

Setup Env is licensed under [Apache License 2.0](LICENSE). Dependency attribution
is summarized in [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).
