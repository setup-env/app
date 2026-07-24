# Install Setup Env on Ubuntu Linux

## Build from source

Requirements:

- Git;
- Go 1.26 (the repository pins toolchain 1.26.5);
- a user-writable directory on `PATH`.

Install Git using Ubuntu's normal package-management policy, and install the
required Go version from an official Go distribution. Then:

```sh
mkdir -p ~/dev/setup-env
cd ~/dev/setup-env
git clone https://github.com/setup-env/app.git
cd app
go build -o bin/setup-env ./cmd/setup-env
./bin/setup-env
```

Launching in a suitable interactive terminal opens the live dashboard. Press
`q` to exit. Redirection, pipelines, and non-TTY sessions print one static
snapshot; use `./bin/setup-env status --json` for automation. `TERM=dumb`
deliberately disables dashboard startup.

For a user-local installation:

```sh
mkdir -p ~/.local/bin
cp bin/setup-env ~/.local/bin/setup-env
```

Add `$HOME/.local/bin` to `PATH` in your shell profile if necessary. No
administrator access is required for this installation.

## Install a future release

Verified GitHub Release binaries and an Ubuntu bootstrap path are planned for
Milestone 05. APT repositories, `.deb` packages, Snap, and other installers are
not available. Future instructions will include checksum verification,
upgrades, and rollback.
