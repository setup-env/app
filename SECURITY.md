# Security Policy

## Reporting a vulnerability

Do not disclose vulnerabilities or secrets in public issues. Use GitHub's
private vulnerability reporting for `setup-env/app` when available. If that
feature is unavailable, contact an organization owner privately and provide a
minimal reproduction, affected versions, and impact.

The maintainers will acknowledge a report, investigate it, coordinate a fix,
and publish remediation guidance when appropriate.

## Security boundaries

Milestone 01 detects command availability and authentication readiness only.
It does not read, print, export, or store Git credentials, SSH keys, credential
helper data, or provider tokens. Configuration files must not contain secrets.

Future module execution must add explicit permission declarations, least
privilege, dry-run support, redaction, checksum and release verification,
auditable execution records, and clear confirmation before destructive
actions.
