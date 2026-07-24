#!/bin/sh

INSTALLER_VERSION=1
REPOSITORY=setup-env/app

fail() {
    printf '%s\n' "setup-env installer: $*" >&2
    return 1
}

normalize_version() {
    value=$1
    if [ "$value" = "latest" ]; then
        printf '%s\n' latest
        return
    fi
    case "$value" in
        v*) normalized=$value ;;
        *) normalized=v$value ;;
    esac
    if ! printf '%s\n' "$normalized" | grep -Eq '^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-[0-9A-Za-z.-]+)?$'; then
        fail "version '$value' is invalid; use vMAJOR.MINOR.PATCH"
        return 1
    fi
    printf '%s\n' "$normalized"
}

detect_os() {
    value=${1:-$(uname -s)}
    case "$value" in
        Darwin|darwin) printf '%s\n' darwin ;;
        Linux|linux) printf '%s\n' linux ;;
        *) fail "unsupported operating system '$value'; Setup Env supports macOS and Linux" ;;
    esac
}

detect_arch() {
    value=${1:-$(uname -m)}
    case "$value" in
        x86_64|amd64) printf '%s\n' amd64 ;;
        arm64|aarch64) printf '%s\n' arm64 ;;
        *) fail "unsupported architecture '$value'; Setup Env supports amd64 and arm64" ;;
    esac
}

asset_name() {
    version=${1#v}
    os_name=$2
    architecture=$3
    printf 'setup-env_%s_%s_%s.tar.gz\n' "$version" "$os_name" "$architecture"
}

checksum_for() {
    checksums=$1
    asset=$2
    checksum=$(awk -v asset="$asset" '$2 == asset && length($1) == 64 && $1 ~ /^[0-9a-fA-F]+$/ { print tolower($1) }' "$checksums")
    count=$(printf '%s\n' "$checksum" | awk 'NF { count++ } END { print count+0 }')
    if [ "$count" -ne 1 ]; then
        fail "published checksums do not contain exactly one entry for '$asset'"
        return 1
    fi
    printf '%s\n' "$checksum"
}

sha256_file() {
    path=$1
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "$path" | awk '{ print tolower($1) }'
    elif command -v shasum >/dev/null 2>&1; then
        shasum -a 256 "$path" | awk '{ print tolower($1) }'
    else
        fail "SHA-256 verification requires sha256sum or shasum"
        return 1
    fi
}

http_get() {
    uri=$1
    destination=$2
    if command -v curl >/dev/null 2>&1; then
        curl --fail --location --silent --show-error \
            --user-agent "setup-env-installer/$INSTALLER_VERSION" \
            --output "$destination" "$uri"
    elif command -v wget >/dev/null 2>&1; then
        wget --quiet --user-agent="setup-env-installer/$INSTALLER_VERSION" \
            --output-document="$destination" "$uri"
    else
        fail "downloads require curl or wget"
        return 1
    fi
}

resolve_latest_version() {
    if [ -n "${SETUP_ENV_RELEASE_DIR:-}" ] && [ -n "${SETUP_ENV_TEST_VERSION:-}" ]; then
        normalize_version "$SETUP_ENV_TEST_VERSION"
        return
    fi
    api_base=${SETUP_ENV_API_BASE_URL:-https://api.github.com}
    temporary_response=$1
    uri=${api_base%/}/repos/$REPOSITORY/releases/latest
    if ! http_get "$uri" "$temporary_response"; then
        fail "unable to resolve the latest release; specify --version vMAJOR.MINOR.PATCH or check proxy, network, and GitHub rate limits"
        return 1
    fi
    tag=$(sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' "$temporary_response" | sed -n '1p')
    if [ -z "$tag" ]; then
        fail "latest release response did not contain tag_name"
        return 1
    fi
    normalize_version "$tag"
}

copy_release_file() {
    version=$1
    name=$2
    destination=$3
    if [ -n "${SETUP_ENV_RELEASE_DIR:-}" ]; then
        source_path=$SETUP_ENV_RELEASE_DIR/$name
        if [ ! -f "$source_path" ]; then
            fail "local release fixture '$source_path' does not exist"
            return 1
        fi
        cp "$source_path" "$destination"
        return
    fi
    download_base=${SETUP_ENV_DOWNLOAD_BASE_URL:-https://github.com/$REPOSITORY/releases/download}
    uri=${download_base%/}/$version/$name
    if ! http_get "$uri" "$destination"; then
        fail "unable to download '$name' from '$uri'; check the version, proxy, network access, or GitHub rate limits"
        return 1
    fi
}

path_contains() {
    wanted=${1%/}
    old_ifs=$IFS
    IFS=:
    for entry in ${2:-}; do
        if [ "${entry%/}" = "$wanted" ]; then
            IFS=$old_ifs
            return 0
        fi
    done
    IFS=$old_ifs
    return 1
}

confirm_action() {
    message=$1
    assume_yes=$2
    if [ "$assume_yes" -eq 1 ]; then
        return
    fi
    printf '%s [y/N] ' "$message"
    read answer
    case "$answer" in
        y|Y|yes|YES) ;;
        *) fail "operation canceled"; return 1 ;;
    esac
}

uninstall_setup_env() {
    install_dir=$1
    purge=$2
    assume_yes=$3
    binary=$install_dir/setup-env
    state=$install_dir/.setup-env-install
    backup=$binary.backup
    if [ ! -e "$binary" ] && [ ! -e "$state" ] && [ ! -e "$backup" ]; then
        printf 'Setup Env is not installed in %s.\n' "$install_dir"
        return
    fi
    confirm_action "Remove Setup Env installer-owned files from '$install_dir'?" "$assume_yes"
    rm -f "$binary" "$state" "$backup"
    for staged in "$install_dir"/.setup-env.new.*; do
        [ -e "$staged" ] && rm -f "$staged"
    done
    if [ "$purge" -eq 1 ]; then
        confirm_action "Also remove Setup Env configuration and cache data?" "$assume_yes"
        os_name=$(detect_os)
        if [ "$os_name" = darwin ]; then
            rm -rf "$HOME/Library/Application Support/setup-env" "$HOME/Library/Caches/setup-env"
        else
            rm -rf "${XDG_CONFIG_HOME:-$HOME/.config}/setup-env" "${XDG_CACHE_HOME:-$HOME/.cache}/setup-env"
        fi
    fi
    printf '%s\n' "Setup Env was removed. Existing projects and development repositories were not changed."
}

print_path_guidance() {
    install_dir=$1
    if path_contains "$install_dir" "${PATH:-}"; then
        return
    fi
    printf '\n%s\n' "'$install_dir' is not currently on PATH."
    printf '%s\n' "Bash: add  export PATH=\"$install_dir:\$PATH\"  to ~/.profile"
    printf '%s\n' "Zsh:  add  export PATH=\"$install_dir:\$PATH\"  to ~/.zprofile"
    printf '%s\n' "Fish: run  fish_add_path \"$install_dir\""
    printf '%s\n' "Restart the terminal after updating your shell configuration."
}

install_setup_env() {
    requested_version=$1
    install_dir=$2
    upgrade=$3
    temporary=$4
    version=$(normalize_version "$requested_version")
    if [ "$version" = latest ]; then
        version=$(resolve_latest_version "$temporary/latest.json")
    fi
    os_name=$(detect_os)
    architecture=$(detect_arch)
    asset=$(asset_name "$version" "$os_name" "$architecture")
    archive=$temporary/$asset
    checksums=$temporary/checksums.txt
    copy_release_file "$version" checksums.txt "$checksums"
    copy_release_file "$version" "$asset" "$archive"
    expected=$(checksum_for "$checksums" "$asset")
    actual=$(sha256_file "$archive")
    if [ "$actual" != "$expected" ]; then
        fail "checksum verification failed for '$asset'; expected $expected but received $actual; the archive will not be executed"
        return 1
    fi
    extract_dir=$temporary/extract
    mkdir -p "$extract_dir"
    tar -xzf "$archive" -C "$extract_dir"
    candidate=$extract_dir/setup-env
    if [ ! -f "$candidate" ]; then
        fail "verified archive does not contain setup-env"
        return 1
    fi
    chmod 755 "$candidate"
    candidate_output=$("$candidate" version 2>&1) || {
        fail "verified binary failed its version check: $candidate_output"
        return 1
    }
    case "$candidate_output" in
        *"$version"*) ;;
        *) fail "verified binary did not report expected version '$version'"; return 1 ;;
    esac

    mkdir -p "$install_dir"
    destination=$install_dir/setup-env
    staged=$install_dir/.setup-env.new.$$
    backup=$destination.backup
    old_version="not installed"
    if [ -f "$destination" ]; then
        old_version=$("$destination" version 2>/dev/null || printf '%s' unknown)
        cp "$destination" "$backup"
    fi
    cp "$candidate" "$staged"
    chmod 755 "$staged"
    if ! mv -f "$staged" "$destination"; then
        [ -f "$backup" ] && mv -f "$backup" "$destination"
        fail "replacement failed; the previous binary was restored where available"
        return 1
    fi
    installed_output=$("$destination" version 2>&1) || {
        [ -f "$backup" ] && mv -f "$backup" "$destination"
        fail "installed binary failed its version check; the previous binary was restored where available"
        return 1
    }
    case "$installed_output" in
        *"$version"*) ;;
        *)
            [ -f "$backup" ] && mv -f "$backup" "$destination"
            fail "installed binary reported unexpected metadata; the previous binary was restored where available"
            return 1
            ;;
    esac
    rm -f "$backup"
    {
        printf 'installed_version=%s\n' "$version"
        printf 'install_directory=%s\n' "$install_dir"
        printf 'installer_version=%s\n' "$INSTALLER_VERSION"
        printf 'installed_at=%s\n' "$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
        printf 'asset_name=%s\n' "$asset"
        printf 'checksum_sha256=%s\n' "$expected"
        printf 'managed_path=false\n'
    } >"$install_dir/.setup-env-install"
    action=installed
    if [ "$upgrade" -eq 1 ] || [ "$old_version" != "not installed" ]; then
        action=upgraded
    fi
    printf 'Setup Env %s was %s successfully.\n' "$version" "$action"
    printf 'Previous: %s\n' "$old_version"
    printf 'Installed: %s\n' "$destination"
    print_path_guidance "$install_dir"
}

setup_env_installer_main() {
    set -eu
    version=latest
    install_dir=${SETUP_ENV_INSTALL_DIR:-$HOME/.local/bin}
    upgrade=0
    uninstall=0
    purge=0
    assume_yes=0
    while [ "$#" -gt 0 ]; do
        case "$1" in
            --version)
                [ "$#" -ge 2 ] || { fail "--version requires a value"; return 1; }
                version=$2
                shift 2
                ;;
            --install-dir)
                [ "$#" -ge 2 ] || { fail "--install-dir requires a value"; return 1; }
                install_dir=$2
                shift 2
                ;;
            --upgrade) upgrade=1; shift ;;
            --uninstall) uninstall=1; shift ;;
            --purge) purge=1; shift ;;
            --yes) assume_yes=1; shift ;;
            --help|-h)
                printf '%s\n' "Usage: install.sh [--version vX.Y.Z] [--upgrade] [--uninstall [--purge]] [--install-dir PATH] [--yes]"
                return
                ;;
            *) fail "unknown option '$1'"; return 1 ;;
        esac
    done
    if [ "$upgrade" -eq 1 ] && [ "$uninstall" -eq 1 ]; then
        fail "--upgrade and --uninstall cannot be used together"
        return 1
    fi
    if [ "$purge" -eq 1 ] && [ "$uninstall" -ne 1 ]; then
        fail "--purge is valid only with --uninstall"
        return 1
    fi
    if [ "$uninstall" -eq 1 ]; then
        uninstall_setup_env "$install_dir" "$purge" "$assume_yes"
        return
    fi
    command -v tar >/dev/null 2>&1 || { fail "tar is required"; return 1; }
    temporary_base=${TMPDIR:-/tmp}
    temporary=$(mktemp -d "$temporary_base/setup-env-install.XXXXXX")
    trap 'rm -rf "$temporary"' EXIT HUP INT TERM
    install_setup_env "$version" "$install_dir" "$upgrade" "$temporary"
    rm -rf "$temporary"
    trap - EXIT HUP INT TERM
}

case "${0##*/}" in
    install.sh) setup_env_installer_main "$@" ;;
esac
