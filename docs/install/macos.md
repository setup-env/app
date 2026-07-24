# Install Setup Env on Apple macOS

## Support and prerequisites

Supported architectures are Intel (`amd64`) and Apple Silicon (`arm64`). The
installer needs a POSIX shell, `tar`, `curl` or `wget`, and `shasum -a 256` or
`sha256sum`. It does not require Git, Go, `gh`, Homebrew, or root.

Official assets exist only when the requested tag appears on
[GitHub Releases](https://github.com/setup-env/app/releases). Initial binaries
are unsigned and unnotarized, so Gatekeeper may warn.

## Recommended reviewed installer

```sh
version=v0.1.0
curl -fL https://raw.githubusercontent.com/setup-env/app/main/install.sh -o install.sh
less install.sh
sh install.sh --version "$version"
rm install.sh
```

The default directory is `~/.local/bin`. The script never edits profiles; if
needed, it prints Bash, Zsh, and Fish PATH instructions.

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
asset=setup-env_0.1.0_darwin_arm64.tar.gz # use amd64 on Intel
expected=$(awk -v asset="$asset" '$2 == asset {print $1}' checksums.txt)
actual=$(shasum -a 256 "$asset" | awk '{print $1}')
test "$actual" = "$expected"
tar -xzf "$asset"
mkdir -p "$HOME/.local/bin"
install -m 755 setup-env "$HOME/.local/bin/setup-env"
```

Download the archive and `checksums.txt` from the same official release. Add
`~/.local/bin` to the appropriate shell PATH, then run `setup-env version` and
`setup-env`. Never extract or execute after a failed checksum.

## Upgrade, rollback, and uninstall

Upgrade verifies and executes the candidate first, retains a backup until the
installed command passes, and restores it after failed validation. Uninstall
removes only `setup-env`, `.setup-env-install`, and the installer backup.
Configuration/cache require explicit `--purge`; projects remain untouched.

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
- The installer does not remove quarantine attributes or weaken Gatekeeper.
- Homebrew, PKG, DMG, and notarized distribution are not available.
