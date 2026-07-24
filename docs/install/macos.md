# Install Setup Env on Apple macOS

## Build from source

Requirements:

- Git (the Xcode Command Line Tools provide it);
- Go 1.26 (the repository pins toolchain 1.26.5);
- a user-writable directory on `PATH`.

```sh
mkdir -p ~/dev/setup-env
cd ~/dev/setup-env
git clone https://github.com/setup-env/app.git
cd app
go build -o bin/setup-env ./cmd/setup-env
./bin/setup-env
```

Launching in Terminal, iTerm2, or another suitable interactive terminal opens
the live dashboard. Press `q` to exit. Redirection and pipelines print one
static snapshot; use `./bin/setup-env status --json` for automation. Set
`TERM` to the terminal's correct value; `TERM=dumb` deliberately disables the
dashboard.

For a user-local installation:

```sh
mkdir -p ~/.local/bin
cp bin/setup-env ~/.local/bin/setup-env
```

Add `$HOME/.local/bin` to `PATH` in your shell profile if it is not already
present. This path does not require administrator access.

## Install a future release

Signed GitHub Release binaries and a macOS bootstrap path are planned for
Milestone 05. Homebrew, PKG, and DMG distribution are not available. This guide
will be updated with verification, upgrade, and rollback instructions when an
official release channel exists.
