# Official Module Catalog

`modules.yaml` is the authoritative machine-readable list consumed by the
Setup Env CLI. It determines whether a module is listed, where its repository
and manifest are located, its trust classification, and its listing status.

The module's `setup-env.yaml` is authoritative for capabilities, platforms,
compatibility, workflows, and module metadata. A catalog entry does not make a
module runnable. In particular, `planned` entries have no validated, usable
manifest and workflows yet.

Entries must be sorted by module ID. Categories and tags within an entry must
also be sorted. The catalog validator rejects duplicate IDs, duplicate
repositories, unsupported trust/status values, invalid paths, forbidden app or
Awesome entries, and self-assigned official trust outside `setup-env`.

The embedded catalog is the Milestone 02 default. Future milestones may add a
verified cached catalog or an explicit local override; proposed precedence is:

1. an explicit user-selected local catalog;
2. a verified cached catalog;
3. the embedded catalog.

Only the embedded source is active today.

[`setup-env/awesome-setup-env`](https://github.com/setup-env/awesome-setup-env)
is a human-curated discovery list. The CLI never scrapes or executes entries
from its Markdown.
