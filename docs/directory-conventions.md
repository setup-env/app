# Directory Conventions

Setup Env resolves the current user's home directory at runtime and uses:

```text
<home>/dev/<organization>/<repository>
```

`path/filepath` handles local paths. The application never assumes a username,
Windows drive, or Unix home location.

## Context model

| Location | Type | Organization | Repository |
| --- | --- | --- | --- |
| `~/dev` | development root | none | none |
| `~/dev/setup-env` | organization | `setup-env` | none |
| `~/dev/setup-env/app` with Git metadata | repository | `setup-env` | `app` |
| Any unrelated or unverified nested directory | other | structural value when available | none |

Repository classification requires Git metadata. A nested directory is not
called a repository merely because it occupies the expected depth. When Git
has an `origin`, Setup Env parses non-sensitive host, owner, and repository
metadata to compare the local structure with its remote identity.

The current working directory may be anywhere, including inside a repository.
The development root defaults to `filepath.Join(home, "dev")` and can be
overridden by configuration.

## Platform examples

```text
Windows: C:\Users\<user>\dev\<organization>\<repository>
macOS:   /Users/<user>/dev/<organization>/<repository>
Linux:   /home/<user>/dev/<organization>/<repository>
```

These are illustrative shapes, not hard-coded paths.
