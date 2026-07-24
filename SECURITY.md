# Security Policy

## Reporting a vulnerability

Do not disclose installer, release, supply-chain, credential, or application
vulnerabilities in public issues. Use GitHub private vulnerability reporting for
`setup-env/app` when available. Otherwise contact an organization owner
privately with the affected version, impact, and minimal reproduction.

## Release and installer trust model

- Official installers and binaries are published only from `setup-env/app`.
- Bootstrap scripts download only from the GitHub Releases channel for this
  repository and use public unauthenticated HTTPS.
- Release archives are verified against the release’s SHA-256 `checksums.txt`
  before extraction or execution.
- Explicit versions use predictable asset URLs and do not require API discovery.
- Installer metadata contains version, asset, checksum, timestamp, location, and
  PATH ownership only. It never contains credentials or tokens.
- Installation is current-user scoped by default and requires no elevation.
- Installers do not disable PowerShell policy globally, Gatekeeper, SmartScreen,
  antivirus, TLS validation, or other host security controls.
- Uninstall removes only installer-owned files. Configuration/cache removal
  requires an explicit purge option; projects and repositories are never removed.
- The application does not contact GitHub for automatic update checks.

SHA-256 verification detects corruption or substitution relative to the checksum
file published in the same GitHub Release. It does not provide an independent
trust root and is not equivalent to Authenticode signing, Apple signing and
notarization, or external provenance verification.

Initial Windows binaries may be unsigned. Initial macOS binaries may be unsigned
and unnotarized. Users may see SmartScreen or Gatekeeper warnings. Do not bypass
those warnings without reviewing the repository, installer, release origin, and
checksum.

## Workflow security

Ordinary push and pull-request CI has `contents: read` only and cannot publish.
The separate tag workflow has only `contents: write`, checks out the exact tag,
validates semantic version and release uniqueness, runs formatting, vet, tests,
Linux race tests, catalog validation, installer tests, and artifact verification
before creating a release. Official actions are pinned to immutable commit SHAs.

If publication fails, the workflow removes a partial release while preserving
the tag for investigation. Published version tags and releases must never be
recreated with different contents.

Artifact attestations, Windows Authenticode, and Apple signing/notarization are
deferred until maintainable credential and verification processes exist.

## Application boundaries

Setup Env detects command availability and authentication readiness only. It
does not read, print, export, or store Git credentials, SSH keys, credential
helper data, or provider tokens. Module download and execution are not
implemented.
