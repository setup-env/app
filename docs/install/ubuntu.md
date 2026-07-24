# Install Setup Env on Ubuntu Linux

## Support and prerequisites

Supported architectures are amd64 and arm64. The installer needs a POSIX shell,
`tar`, `sha256sum` (or `shasum`), and `curl` or `wget`. It does not require Git,
Go, `gh`, `sudo`, or a package manager.

Official assets exist only when the requested tag appears on
[GitHub Releases](https://github.com/setup-env/app/releases).

## Recommended reviewed installer

```sh
version=v0.1.0
curl -fL https://raw.githubusercontent.com/setup-env/app/main/install.sh -o install.sh
less install.sh
sh install.sh --version "$version"
rm install.sh
```

The default directory is `~/.local/bin`. The script prints Bash, Zsh, and Fish
PATH guidance rather than modifying a profile.

```sh
sh install.sh                              # latest release
sh install.sh --version v0.1.0             # explicit version
sh install.sh --upgrade                    # latest verified upgrade
sh install.sh --install-dir "$HOME/bin"    # custom path
sh install.sh --uninstall --yes
sh install.sh --uninstall --purge --yes    # also configuration/cache
```

## Manual archive installation

```sh
asset=setup-env_0.1.0_linux_amd64.tar.gz # use arm64 when appropriate
expected=$(awk -v asset="$asset" '$2 == asset {print $1}' checksums.txt)
actual=$(sha256sum "$asset" | awk '{print $1}')
test "$actual" = "$expected"
tar -xzf "$asset"
mkdir -p "$HOME/.local/bin"
install -m 755 setup-env "$HOME/.local/bin/setup-env"
```

Download the archive and `checksums.txt` from the same official release. Add
`~/.local/bin` to PATH, then run `setup-env version` and `setup-env`. Never
extract or execute after failed verification.

## Upgrade, rollback, and uninstall

Upgrade verifies and executes the candidate first, retains the working binary
until replacement succeeds, and restores it on failed validation. Uninstall
removes only installer-owned binary and metadata. Configuration/cache require
explicit `--purge`; projects remain untouched.

## Build from source

```sh
git clone https://github.com/setup-env/app.git
cd app
go test ./...
go build -o bin/setup-env ./cmd/setup-env
./bin/setup-env version
```

Git and Go 1.26.5 are required only for source builds.

## Troubleshooting and restricted networks

- Explicit versions avoid latest-release API rate limits.
- Standard `HTTPS_PROXY` and `NO_PROXY` behavior from curl/wget applies.
- A 404 usually means the tag or architecture asset is not published.
- Verify system time and TLS interception policy when HTTPS fails.
- APT repositories, `.deb`, Snap, and Flatpak are not available.
