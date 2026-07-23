# Milestone 01 — Application Foundation

This document is the durable plan for the first application milestone and the
fallback tracker if GitHub milestone or issue creation is unavailable.

## GitHub tracking status

On 2026-07-24, milestone creation was unavailable through the connected GitHub
interface and the local `gh` authentication required renewal. A direct attempt
to create the parent issue through the connected GitHub app returned
`403 Resource not accessible by integration`. No GitHub milestone or issues
were created. The titles and descriptions below are ready to create after
repository issue permissions and `gh` authentication are corrected.

## Parent issue

### Complete Milestone 01 — Application Foundation

Establish a clean, tested, cross-platform Go foundation for the universal Setup
Env application. Deliver the initial CLI, environment and directory detection,
safe diagnostics, configuration foundation, CI, documentation, and governance.
Do not implement module-specific behavior or the workflow engine.

## Child workstreams

### Establish the Go project and initial CLI

Create the Go module and `setup-env`, `version`, `info`, and `doctor` commands.
Keep dependencies minimal and design command boundaries for future expansion.

### Implement cross-platform platform and directory detection

Detect OS, distribution, architecture, user, home, working directory,
development root, organization/repository context, and sanitized Git remote
identity with injected tests.

### Add Git, GitHub, and permission diagnostics

Report Git and GitHub CLI availability, versions, authentication readiness, and
development-root write access without reading or exposing credentials.

### Define versioned application configuration

Use OS-native user configuration directories and support a schema version,
development-root override, clone protocol, known organizations, and application
settings without secrets.

### Add tests and cross-platform CI

Test path models, directory context, configuration, platform data, diagnostics,
and commands. Validate formatting, vetting, tests, builds on Windows, macOS, and
Linux, plus race tests where supported.

### Document architecture, modules, development, and roadmap

Explain application/module separation, the provisional module contract, trust
and verification direction, directory conventions, contribution and security
policies, and future milestones.

## Completion criteria

Milestone 01 is complete when local validation passes, CI is configured for all
target operating systems, the current commands work in human and requested JSON
modes, documentation distinguishes implemented from proposed capability, and
the default branch contains the committed foundation.
