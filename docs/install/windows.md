# Install Setup Env on Windows

## Support and prerequisites

Supported architectures are Windows amd64 and arm64. The bootstrap supports
Windows PowerShell 5.1 and PowerShell 7 and uses built-in HTTPS,
`Get-FileHash`, and `Expand-Archive`. It does not require Git, Go, `gh`, or
administrator rights.

Official assets exist only when the requested tag appears on
[GitHub Releases](https://github.com/setup-env/app/releases). Windows binaries
are initially unsigned and may trigger SmartScreen.

## Recommended reviewed installer

```powershell
$version = "v0.1.0"
Invoke-WebRequest https://raw.githubusercontent.com/setup-env/app/main/install.ps1 -OutFile install.ps1
Get-Content .\install.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File .\install.ps1 -Version $version
Remove-Item .\install.ps1
```

The process-scoped execution-policy argument does not change user or machine
policy. The default directory is
`%LOCALAPPDATA%\Programs\setup-env\bin`. The installer verifies before
extraction, tests `setup-env version` before and after replacement, and adds the
directory to current-user PATH without duplicates.

```powershell
.\install.ps1                         # latest published release
.\install.ps1 -Version v0.1.0         # explicit version
.\install.ps1 -Upgrade                # latest verified upgrade
.\install.ps1 -InstallDir C:\My\Bin   # custom user-writable path
.\install.ps1 -NoPath                 # do not modify user PATH
.\install.ps1 -Uninstall -Yes
.\install.ps1 -Uninstall -Purge -Yes  # also configuration/cache
```

## Manual archive installation

1. Download the matching `.zip` and `checksums.txt` from the same release.
2. Verify the archive:

   ```powershell
   $asset = "setup-env_0.1.0_windows_amd64.zip" # use arm64 when appropriate
   $line = (Select-String -Path .\checksums.txt -Pattern "  $([regex]::Escape($asset))$").Line
   $expected = ($line -split '\s+')[0].ToLowerInvariant()
   $actual = (Get-FileHash -Algorithm SHA256 -Path ".\$asset").Hash.ToLowerInvariant()
   if ($actual -ne $expected) { throw "checksum mismatch" }
   ```

3. Expand the verified archive and copy `setup-env.exe` to a user-controlled
   executable directory.
4. Add that directory to user PATH without replacing existing entries.
5. Open a new terminal and run `setup-env version`, then `setup-env`.

Never execute the binary when verification fails.

## Upgrade, rollback, and uninstall

Upgrade verifies and executes a temporary candidate, keeps the existing binary
as a backup, replaces it, and tests it again. A failed replacement restores the
backup where available.

Uninstall removes only `setup-env.exe`, its temporary backup, and
`.setup-env-install.json`. It removes PATH only when metadata says the installer
added it. Configuration/cache remain unless `-Purge` is explicit. Projects and
development repositories are never removed.

## Build from source

Git and Go 1.26.5 are required only for source builds:

```powershell
git clone https://github.com/setup-env/app.git
Set-Location app
go test ./...
go build -o bin\setup-env.exe .\cmd\setup-env
.\bin\setup-env.exe version
```

## Troubleshooting and restricted networks

- Use `-Version vX.Y.Z` to avoid the latest-release API request.
- Standard PowerShell proxy and TLS settings apply; no GitHub token is stored.
- A 404 usually means the tag or architecture asset is not published.
- API rate-limit failures suggest using an explicit version.
- The installer does not disable SmartScreen.
- If policy blocks process-scoped script execution, use manual installation
  instead of changing machine policy.
