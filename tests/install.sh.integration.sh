#!/bin/sh
set -eu

if [ "$#" -lt 1 ]; then
    printf '%s\n' "usage: install.sh.integration.sh RELEASE_DIR [VERSION [UPGRADE_RELEASE_DIR UPGRADE_VERSION]]" >&2
    exit 2
fi

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
RELEASE_DIR=$1
VERSION=${2:-v0.1.0}
UPGRADE_RELEASE_DIR=${3:-$RELEASE_DIR}
UPGRADE_VERSION=${4:-$VERSION}
root=$(mktemp -d "${TMPDIR:-/tmp}/setup-env-shell-integration.XXXXXX")
trap 'rm -rf "$root"' EXIT HUP INT TERM
fixture_dir=$root/fixtures
upgrade_fixture_dir=$root/upgrade-fixtures
install_dir=$root/install
scratch_dir=$root/tmp
mkdir -p "$fixture_dir" "$upgrade_fixture_dir" "$scratch_dir"
cp "$RELEASE_DIR"/* "$fixture_dir/"
cp "$UPGRADE_RELEASE_DIR"/* "$upgrade_fixture_dir/"
export SETUP_ENV_RELEASE_DIR=$fixture_dir
export SETUP_ENV_TEST_VERSION=$VERSION
export TMPDIR=$scratch_dir

sh "$SCRIPT_DIR/../install.sh" --version "$VERSION" --install-dir "$install_dir" --yes
test -x "$install_dir/setup-env"
test -f "$install_dir/.setup-env-install"
version_output=$("$install_dir/setup-env" version)
case "$version_output" in
    *"$VERSION"*) ;;
    *) printf '%s\n' "installed binary metadata is incorrect: $version_output" >&2; exit 1 ;;
esac

printf '%s\n' preserve >"$install_dir/unrelated.txt"
export SETUP_ENV_RELEASE_DIR=$upgrade_fixture_dir
export SETUP_ENV_TEST_VERSION=$UPGRADE_VERSION
sh "$SCRIPT_DIR/../install.sh" --version "$UPGRADE_VERSION" --upgrade --install-dir "$install_dir" --yes
upgraded_output=$("$install_dir/setup-env" version)
case "$upgraded_output" in
    *"$UPGRADE_VERSION"*) ;;
    *) printf '%s\n' "upgraded binary metadata is incorrect: $upgraded_output" >&2; exit 1 ;;
esac
sh "$SCRIPT_DIR/../install.sh" --uninstall --install-dir "$install_dir" --yes
test ! -e "$install_dir/setup-env"
test ! -e "$install_dir/.setup-env-install"
test -f "$install_dir/unrelated.txt"

os_name=$(uname -s)
case "$os_name" in
    Darwin) artifact_os=darwin ;;
    Linux) artifact_os=linux ;;
    *) printf '%s\n' "unsupported integration-test OS: $os_name" >&2; exit 1 ;;
esac
machine=$(uname -m)
case "$machine" in
    x86_64|amd64) artifact_arch=amd64 ;;
    arm64|aarch64) artifact_arch=arm64 ;;
    *) printf '%s\n' "unsupported integration-test architecture: $machine" >&2; exit 1 ;;
esac
plain_version=${VERSION#v}
export SETUP_ENV_RELEASE_DIR=$fixture_dir
export SETUP_ENV_TEST_VERSION=$VERSION
asset=$fixture_dir/setup-env_${plain_version}_${artifact_os}_${artifact_arch}.tar.gz
printf '%s\n' corrupt >>"$asset"
corrupt_install=$root/corrupt-install
if sh "$SCRIPT_DIR/../install.sh" --version "$VERSION" --install-dir "$corrupt_install" --yes >/dev/null 2>&1; then
    printf '%s\n' "corrupt archive was not rejected" >&2
    exit 1
fi
test ! -e "$corrupt_install/setup-env"

if find "$scratch_dir" -mindepth 1 -maxdepth 1 -name 'setup-env-install.*' | grep . >/dev/null 2>&1; then
    printf '%s\n' "installer temporary directory cleanup is incomplete" >&2
    exit 1
fi

printf '%s\n' "POSIX installer integration tests passed."
