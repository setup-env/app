# Module Model

> Status: initial proposal. This specification is not implemented and will
> change as Milestones 02–05 validate it against real modules.

A module owns one setup domain and releases independently from the application.
The application owns discovery, validation, download, verification, execution,
and shared behavior. A module owns its workflows, scripts, templates, tests,
and domain documentation.

Expected repository layout:

```text
<module>/
├── setup-env.yaml
├── README.md
├── workflows/
├── templates/
├── scripts/
└── tests/
```

The future manifest is expected to declare:

- schema and module versions;
- identity, description, maintainers, and trust metadata;
- application compatibility;
- supported operating systems, distributions, and architectures;
- workflows, typed inputs, outputs, and dependencies;
- required tools, privileges, permissions, network access, and secrets;
- checksums or references to verified release artifacts.

Workflows should describe intent and use shared actions where possible. Scripts
remain appropriate for domain behavior that cannot be expressed safely in a
portable action. Every workflow must support validation before mutation, expose
destructive behavior, and define idempotency and rollback expectations.

Conceptual future commands include:

```text
setup-env module list
setup-env module info <module>
setup-env module install <module>
setup-env module update <module>
setup-env workflow list <module>
setup-env run <module> <workflow>
```

These commands are not implemented in Milestone 01.
