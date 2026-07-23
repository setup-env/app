# Module Contributions

Module registration is a pull-request review process. Catalog inclusion is not
automatic installation approval and does not execute repository content.

## Registration process

1. Create a public module repository.
2. Add a valid `setup-env.yaml` using manifest schema version 1.
3. Document purpose, workflows, platforms, prerequisites, permissions, and
   security behavior.
4. Add deterministic tests where applicable.
5. Create a feature branch in or fork of `setup-env/app`.
6. Add one sorted entry to `catalog/modules.yaml`.
7. Run `setup-env module validate`, `setup-env module validate-catalog`, and
   the standard Go checks.
8. Submit a pull request and complete the module-submission template.
9. Receive maintainer approval for the entry and trust classification.
10. Merge to `main` after required validation succeeds.

## Required manifest information

The manifest must contain identity, repository, version source, publisher,
license, minimum app version, supported platforms, categories, security
declarations, and at least one workflow metadata declaration. It must contain
no tokens, credentials, private keys, or other secrets.

## Catalog entry

Add only discovery and governance metadata: ID, display name, short
description, repository, manifest location, trust, status, categories, optional
tags, and optional publisher/version policy. Do not duplicate platform or
workflow capability fields from the manifest.

Entries, categories, and tags must remain alphabetically sorted.

## Trust review

- **Official**: owned by `setup-env`, governed and reviewed by Setup Env
  maintainers. Only maintainers can assign this level.
- **Verified**: externally owned and reviewed against published technical,
  maintenance, identity, and security requirements. It is not maintained or
  endorsed as an official module.
- **Community**: discoverable but not deeply reviewed or endorsed.

A contributor may request community inclusion. Promotion from community to
verified requires a new PR with evidence of stable ownership, maintenance,
release hygiene, tests, documentation, vulnerability handling, and accurate
permission declarations. Contributors cannot self-assign official or verified
trust.

## Security and maintenance expectations

Repositories must avoid embedded secrets, clearly disclose privileges, network
access and secret inputs, use least privilege, document destructive behavior,
and respond to vulnerabilities. Maintainers should publish versioned releases,
keep compatibility accurate, review dependencies, and provide deprecation
notice before abandonment.

## Deprecation and replacement

Deprecation requires a catalog status update and, when possible, manifest
`deprecated`, `deprecation_notice`, and `replacement` fields. Replacements need
independent catalog review; naming one does not grant it trust.

If a repository is transferred, renamed, archived, or changes control, submit a
catalog PR immediately. Trust does not automatically survive transfer. The
entry may be marked unavailable while ownership and security are re-reviewed.

## Rejection reasons

A submission may be rejected for an invalid or misleading manifest, duplicated
scope, unclear ownership, unavailable source, unsafe secret handling,
undeclared privileges, missing maintenance expectations, trademark or license
problems, malicious or arbitrary execution behavior, false compatibility
claims, self-assigned trust, or failure to pass deterministic validation.

## Awesome list

`awesome-setup-env` is separate human curation. Inclusion there neither grants
catalog inclusion nor changes trust. The CLI does not scrape its README.
