Setup Env’s first public release provides the live terminal dashboard, static
and JSON status commands, diagnostics, and the local module catalog.

Supported artifacts:

- Windows amd64 and arm64 (`.zip`)
- Apple macOS Intel and Apple Silicon (`.tar.gz`)
- Linux amd64 and arm64 (`.tar.gz`)

Each archive is covered by `checksums.txt`. Verify SHA-256 before extraction,
or use the repository bootstrap installer, which performs that verification
before executing or replacing a binary.

After installation:

```text
setup-env version
setup-env
```

Windows and macOS binaries are initially unsigned; macOS artifacts are not
notarized. SmartScreen or Gatekeeper may warn. SHA-256 verification protects
artifact integrity relative to the published checksum file but is not a
replacement for platform code signing.

Upgrade and uninstall commands, manual verification, known limitations, and
security details are documented in the repository README and platform guides.

This release does not download or execute Setup Env modules and does not include
Workstation setup behavior, package-manager distribution, a GUI, telemetry, or
remote management.
