#!/bin/sh
set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
. "$SCRIPT_DIR/../install.sh"

assert_equal() {
    actual=$1
    expected=$2
    message=$3
    if [ "$actual" != "$expected" ]; then
        printf '%s\n' "$message: got '$actual', want '$expected'" >&2
        exit 1
    fi
}

assert_fails() {
    message=$1
    shift
    if "$@" >/dev/null 2>&1; then
        printf '%s\n' "$message: command unexpectedly succeeded" >&2
        exit 1
    fi
}

assert_equal "$(normalize_version 0.1.0)" v0.1.0 "version normalization"
assert_equal "$(normalize_version v1.2.3-rc.1)" v1.2.3-rc.1 "prerelease normalization"
assert_fails "invalid version" normalize_version latest-ish
assert_fails "build metadata version" normalize_version v1.2.3+build.4
assert_equal "$(detect_os Darwin)" darwin "Darwin mapping"
assert_equal "$(detect_os Linux)" linux "Linux mapping"
assert_fails "unsupported OS" detect_os FreeBSD
assert_equal "$(detect_arch x86_64)" amd64 "x86_64 mapping"
assert_equal "$(detect_arch aarch64)" arm64 "aarch64 mapping"
assert_fails "unsupported architecture" detect_arch i386
assert_equal \
    "$(asset_name v0.1.0 linux arm64)" \
    setup-env_0.1.0_linux_arm64.tar.gz \
    "asset naming"

temporary=$(mktemp -d "${TMPDIR:-/tmp}/setup-env-shell-unit.XXXXXX")
trap 'rm -rf "$temporary"' EXIT HUP INT TERM
hash=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
printf '%s  %s\n' "$hash" setup-env_0.1.0_linux_amd64.tar.gz >"$temporary/checksums.txt"
assert_equal \
    "$(checksum_for "$temporary/checksums.txt" setup-env_0.1.0_linux_amd64.tar.gz)" \
    "$hash" \
    "checksum parsing"
assert_fails "missing checksum" checksum_for "$temporary/checksums.txt" missing.tar.gz
path_contains "/home/example/.local/bin" "/usr/bin:/home/example/.local/bin:/bin"
assert_fails "PATH false positive" path_contains "/home/example/.local/bin" "/usr/bin:/home/example/.local/bin-extra"

mkdir "$temporary/fixtures"
SETUP_ENV_RELEASE_DIR=$temporary/fixtures
export SETUP_ENV_RELEASE_DIR
assert_fails "unavailable local release" copy_release_file v9.9.9 missing.tar.gz "$temporary/out.tar.gz"

printf '%s\n' "POSIX installer unit tests passed."
