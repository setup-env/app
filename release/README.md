# Setup Env release archive

This archive contains the native `setup-env` command for one operating system
and architecture.

Verify the archive against the release `checksums.txt` before extracting it.
After extraction, run:

```text
setup-env version
setup-env
```

The first command reports release metadata. The second opens the live dashboard
when attached to an interactive terminal and prints a static snapshot otherwise.

Installation, upgrade, uninstall, security, and platform-specific guidance:

https://github.com/setup-env/app#install

Setup Env releases may initially be unsigned and, on macOS, unnotarized.
Checksums detect artifact corruption or substitution relative to the published
checksum file, but are not a replacement for platform code signing.
