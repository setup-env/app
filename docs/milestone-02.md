# Milestone 02 — Module Manifest and Official Catalog

Milestone 02 defines the manifest v1 contract, authoritative embedded catalog,
trust and status governance, compatibility checks, local inspection CLI, test
matrix, and module contribution flow. It does not download, cache, install,
update, verify releases, or execute workflows.

## GitHub tracking status

On 2026-07-24, local `gh auth status` still reported an invalid token. The
connected GitHub app reports repository admin metadata but returned
`403 Resource not accessible by integration` when creating the parent issue.
The connector exposes neither milestone creation nor branch-protection/ruleset
configuration. Therefore no Milestone 02 milestone, issues, or protection
settings were created through those interfaces. The complete issue plan below
is the durable fallback.

## Parent issue

### Complete Milestone 02 — Module Manifest and Official Catalog

Define and implement the first stable module manifest, authoritative app
catalog, trust classification, compatibility validation, local catalog CLI,
contribution flow, and deterministic CI. Preserve Milestone 01 and exclude all
download, installation, caching, verification, and execution behavior.

## Child issues

### Define module manifest v1

Specify the deliberately small `setup-env.yaml` identity, repository, version,
publisher, platform, category, security, workflow-metadata, compatibility, and
deprecation fields. Document identifier and trust boundaries.

### Implement manifest parser and validation

Pin the focused YAML dependency, reject unknown fields, separate parsing,
schema, semantic and compatibility validation, aggregate actionable failures,
and add a reference manifest.

### Add official catalog model

Create and embed `catalog/modules.yaml`, model catalog sources for future
extension, populate the ten current official module repositories honestly, and
keep review ordering deterministic.

### Implement trust and status validation

Validate official, verified and community trust; active, planned, experimental,
deprecated and unavailable status; forbid self-promotion and app/Awesome module
entries; and reject duplicate IDs or repositories.

### Add module CLI commands

Implement local-only `module list`, `module info`, `module validate`, JSON
output, focused filters, actionable errors, and honest unavailable-manifest
messages.

### Add catalog and manifest CI validation

Test the embedded catalog, fixtures, CLI JSON, invalid exit behavior,
compatibility and trust rules without remote code or network-dependent tests.
Run explicit catalog and example validation in CI.

### Document module contribution workflow

Define the branch/PR/validation/approval process, trust reviews, security and
maintenance expectations, deprecation, transfers, promotion and rejection
criteria.

### Document awesome-setup-env relationship

Clarify that the app catalog is authoritative machine data while
`awesome-setup-env` is human curation, grants no trust, is never scraped for
execution, and is not synchronized during this milestone.

## Completion checklist

- [x] Manifest v1 model, strict parser, schema and semantic validation
- [x] Minimum app compatibility states
- [x] Embedded authoritative catalog with ten honest planned entries
- [x] Trust and status validation
- [x] Module list, info, validate and maintainer catalog validation commands
- [x] Reference manifest and deterministic tests
- [x] CI catalog and manifest checks
- [x] Contribution, architecture, Awesome relationship and roadmap updates
- [ ] GitHub milestone and linked issues, subject to GitHub write access
- [ ] Recommended branch protection, subject to repository rules access
