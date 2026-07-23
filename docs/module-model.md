# Module Manifest v1

`setup-env.yaml` is the versioned machine-readable contract for one module
repository. Milestone 02 validates manifest and workflow metadata but does not
download, install, update, or execute modules.

See the single maintained [reference manifest](../examples/setup-env.yaml).
It represents `workstation` as a fixture and is not the live manifest of the
module repository.

## Fields

Manifest v1 supports:

- `schema_version`: must be `1`;
- `id`: stable lowercase identifier using letters, numbers, and single hyphens;
- `name` and `description`;
- `repository.owner`, `repository.name`, and optional `issues_url`;
- `version.source`: `github-release`, `git-tag`, or `manifest`;
- `version.value`: required semantic version only for `manifest` source;
- `publisher` and `license`;
- optional `homepage` and `documentation` HTTP(S) URLs;
- `minimum_app_version` in semantic `MAJOR.MINOR.PATCH` form;
- supported `platforms.operating_systems` and `platforms.architectures`;
- identifier-safe `categories` and optional `tags`;
- non-authoritative security declarations for elevation, network access, and
  named secret inputs;
- workflow metadata with ID, name, description, and a relative YAML entrypoint
  under `workflows/`;
- optional `deprecated`, `replacement`, and `deprecation_notice`.

Unknown fields are rejected so mistakes and unsupported future contracts are
not silently accepted. Validation separates YAML parsing, schema-version
validation, semantic validation, and compatibility evaluation. Multiple
semantic failures are returned together where practical.

## Trust is not a manifest field

Modules cannot declare themselves `official`, `verified`, or `community`.
Trust is assigned by maintainers in `catalog/modules.yaml`. A manifest that
adds a `trust` field fails strict parsing. Security requirements describe what
a workflow may eventually need; they do not endorse its publisher or code.

## Repository layout

```text
<module>/
├── setup-env.yaml
├── README.md
├── workflows/
├── templates/
├── scripts/
└── tests/
```

Workflow files are declarations only in Milestone 02. The future workflow
engine will define their executable contract, inputs, actions, permissions,
idempotency, dry-run, cancellation, redaction, and audit behavior.

## Compatibility

The app compares its released semantic version with
`minimum_app_version`. Results are:

- `compatible`: app version meets the minimum;
- `incompatible`: app version is below the minimum;
- `unknown`: the app is a development build or version data is incomplete.

This is deliberately not a dependency solver and does not evaluate
module-to-module dependencies.

## Catalog relationship

The catalog locates and classifies modules; the manifest describes them. The
catalog wins for listing, repository, trust, and status. The manifest wins for
capabilities, platforms, compatibility, workflows, and module metadata.
