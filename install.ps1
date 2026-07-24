[CmdletBinding()]
param(
    [string]$Version = "latest",
    [string]$InstallDir,
    [switch]$Upgrade,
    [switch]$Uninstall,
    [switch]$Purge,
    [switch]$Yes,
    [switch]$NoPath
)

$ErrorActionPreference = "Stop"
$script:InstallerVersion = "1"
$script:Repository = "setup-env/app"

function Normalize-SetupEnvVersion {
    param([Parameter(Mandatory = $true)][string]$Value)
    $normalized = $Value.Trim()
    if ($normalized -eq "latest") {
        return "latest"
    }
    if (-not $normalized.StartsWith("v")) {
        $normalized = "v$normalized"
    }
    if ($normalized -notmatch '^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-[0-9A-Za-z.-]+)?$') {
        throw "Version '$Value' is invalid. Use vMAJOR.MINOR.PATCH."
    }
    return $normalized
}

function Get-SetupEnvArchitecture {
    param([string]$Value = $env:PROCESSOR_ARCHITECTURE)
    switch ($Value.ToUpperInvariant()) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        default { throw "Unsupported Windows architecture '$Value'. Setup Env supports amd64 and arm64." }
    }
}

function Get-SetupEnvAssetName {
    param(
        [Parameter(Mandatory = $true)][string]$ReleaseVersion,
        [Parameter(Mandatory = $true)][string]$Architecture
    )
    $plainVersion = $ReleaseVersion.TrimStart("v")
    return "setup-env_${plainVersion}_windows_${Architecture}.zip"
}

function Get-SetupEnvDefaultInstallDir {
    if (-not $env:LOCALAPPDATA) {
        throw "LOCALAPPDATA is unavailable; pass -InstallDir explicitly."
    }
    return [IO.Path]::GetFullPath((Join-Path $env:LOCALAPPDATA "Programs\setup-env\bin"))
}

function Resolve-SetupEnvInstallDir {
    param([string]$Value)
    if ($Value) {
        return [IO.Path]::GetFullPath($Value)
    }
    return Get-SetupEnvDefaultInstallDir
}

function Test-SetupEnvPathEntry {
    param(
        [AllowEmptyString()][string]$PathValue,
        [Parameter(Mandatory = $true)][string]$Entry
    )
    $wanted = $Entry.TrimEnd('\', '/')
    foreach ($item in ($PathValue -split ';')) {
        if ($item.Trim().TrimEnd('\', '/').Equals($wanted, [StringComparison]::OrdinalIgnoreCase)) {
            return $true
        }
    }
    return $false
}

function Add-SetupEnvPathEntry {
    param(
        [AllowEmptyString()][string]$PathValue,
        [Parameter(Mandatory = $true)][string]$Entry
    )
    if (Test-SetupEnvPathEntry -PathValue $PathValue -Entry $Entry) {
        return $PathValue
    }
    if ([string]::IsNullOrWhiteSpace($PathValue)) {
        return $Entry
    }
    return "$($PathValue.TrimEnd(';'));$Entry"
}

function Remove-SetupEnvPathEntry {
    param(
        [AllowEmptyString()][string]$PathValue,
        [Parameter(Mandatory = $true)][string]$Entry
    )
    $wanted = $Entry.TrimEnd('\', '/')
    $kept = @()
    foreach ($item in ($PathValue -split ';')) {
        $trimmed = $item.Trim()
        if ($trimmed -and -not $trimmed.TrimEnd('\', '/').Equals($wanted, [StringComparison]::OrdinalIgnoreCase)) {
            $kept += $trimmed
        }
    }
    return ($kept -join ';')
}

function Get-SetupEnvChecksum {
    param(
        [Parameter(Mandatory = $true)][string]$ChecksumsPath,
        [Parameter(Mandatory = $true)][string]$AssetName
    )
    foreach ($line in Get-Content -LiteralPath $ChecksumsPath) {
        if ($line -match '^([0-9a-fA-F]{64})\s+(.+)$' -and $Matches[2] -eq $AssetName) {
            return $Matches[1].ToLowerInvariant()
        }
    }
    throw "The published checksums do not contain '$AssetName'."
}

function Resolve-SetupEnvLatestVersion {
    if ($env:SETUP_ENV_RELEASE_DIR -and $env:SETUP_ENV_TEST_VERSION) {
        return Normalize-SetupEnvVersion $env:SETUP_ENV_TEST_VERSION
    }
    $apiBase = if ($env:SETUP_ENV_API_BASE_URL) {
        $env:SETUP_ENV_API_BASE_URL.TrimEnd('/')
    } else {
        "https://api.github.com"
    }
    $uri = "$apiBase/repos/$script:Repository/releases/latest"
    try {
        $response = Invoke-RestMethod -Uri $uri -Headers @{
            "User-Agent" = "setup-env-installer/$script:InstallerVersion"
            "Accept" = "application/vnd.github+json"
        } -UseBasicParsing
    } catch {
        throw "Unable to resolve the latest Setup Env release from $uri. Specify -Version vMAJOR.MINOR.PATCH or retry later. $($_.Exception.Message)"
    }
    if (-not $response.tag_name) {
        throw "The latest release response did not contain a tag_name."
    }
    return Normalize-SetupEnvVersion ([string]$response.tag_name)
}

function Copy-SetupEnvReleaseFile {
    param(
        [Parameter(Mandatory = $true)][string]$ReleaseVersion,
        [Parameter(Mandatory = $true)][string]$Name,
        [Parameter(Mandatory = $true)][string]$Destination
    )
    if ($env:SETUP_ENV_RELEASE_DIR) {
        $source = Join-Path $env:SETUP_ENV_RELEASE_DIR $Name
        if (-not (Test-Path -LiteralPath $source -PathType Leaf)) {
            throw "Local release fixture '$source' does not exist."
        }
        Copy-Item -LiteralPath $source -Destination $Destination -Force
        return
    }
    $downloadBase = if ($env:SETUP_ENV_DOWNLOAD_BASE_URL) {
        $env:SETUP_ENV_DOWNLOAD_BASE_URL.TrimEnd('/')
    } else {
        "https://github.com/$script:Repository/releases/download"
    }
    $uri = "$downloadBase/$ReleaseVersion/$Name"
    try {
        Invoke-WebRequest -Uri $uri -OutFile $Destination -Headers @{
            "User-Agent" = "setup-env-installer/$script:InstallerVersion"
        } -UseBasicParsing
    } catch {
        throw "Unable to download '$Name' from $uri. Check the version, proxy, network access, or GitHub rate limits. $($_.Exception.Message)"
    }
}

function Confirm-SetupEnvAction {
    param(
        [Parameter(Mandatory = $true)][string]$Message,
        [switch]$AssumeYes
    )
    if ($AssumeYes) {
        return
    }
    $answer = Read-Host "$Message [y/N]"
    if ($answer -notmatch '^(y|yes)$') {
        throw "Operation canceled."
    }
}

function Invoke-SetupEnvUninstall {
    param(
        [Parameter(Mandatory = $true)][string]$TargetDir,
        [switch]$RemoveData,
        [switch]$AssumeYes
    )
    $binary = Join-Path $TargetDir "setup-env.exe"
    $statePath = Join-Path $TargetDir ".setup-env-install.json"
    $state = $null
    if (Test-Path -LiteralPath $statePath -PathType Leaf) {
        try {
            $state = Get-Content -Raw -LiteralPath $statePath | ConvertFrom-Json
        } catch {
            Write-Warning "Installer metadata is unreadable; owned files will be removed but PATH will not be changed."
        }
    }
    $owned = @($binary, $statePath, "$binary.backup")
    $owned += @(Get-ChildItem -LiteralPath $TargetDir -File -Filter ".setup-env.new.*.exe" -ErrorAction SilentlyContinue | Select-Object -ExpandProperty FullName)
    $present = @($owned | Where-Object { Test-Path -LiteralPath $_ })
    if ($present.Count -eq 0) {
        Write-Host "Setup Env is not installed in $TargetDir."
        return
    }
    Confirm-SetupEnvAction -Message "Remove Setup Env installer-owned files from '$TargetDir'?" -AssumeYes:$AssumeYes
    foreach ($path in $owned) {
        if (Test-Path -LiteralPath $path) {
            Remove-Item -LiteralPath $path -Force
        }
    }
    if ($state -and $state.managed_path -eq $true) {
        $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
        $updated = Remove-SetupEnvPathEntry -PathValue $userPath -Entry $TargetDir
        if ($updated -ne $userPath) {
            [Environment]::SetEnvironmentVariable("Path", $updated, "User")
        }
    }
    if ($RemoveData) {
        Confirm-SetupEnvAction -Message "Also remove Setup Env configuration and cache data?" -AssumeYes:$AssumeYes
        foreach ($path in @(
            (Join-Path $env:APPDATA "setup-env"),
            (Join-Path $env:LOCALAPPDATA "setup-env")
        )) {
            if ($path -and (Test-Path -LiteralPath $path)) {
                Remove-Item -LiteralPath $path -Recurse -Force
            }
        }
    }
    Write-Host "Setup Env was removed. Existing projects and development repositories were not changed."
}

function Invoke-SetupEnvInstall {
    param(
        [Parameter(Mandatory = $true)][string]$RequestedVersion,
        [Parameter(Mandatory = $true)][string]$TargetDir,
        [switch]$IsUpgrade,
        [switch]$SkipPath
    )
    if ($env:OS -ne "Windows_NT") {
        throw "install.ps1 supports Windows only."
    }
    $releaseVersion = Normalize-SetupEnvVersion $RequestedVersion
    if ($releaseVersion -eq "latest") {
        $releaseVersion = Resolve-SetupEnvLatestVersion
    }
    $architecture = Get-SetupEnvArchitecture
    $assetName = Get-SetupEnvAssetName -ReleaseVersion $releaseVersion -Architecture $architecture
    $temporary = Join-Path ([IO.Path]::GetTempPath()) ("setup-env-install-" + [Guid]::NewGuid().ToString("N"))
    New-Item -ItemType Directory -Path $temporary -Force | Out-Null
    try {
        $archivePath = Join-Path $temporary $assetName
        $checksumsPath = Join-Path $temporary "checksums.txt"
        Copy-SetupEnvReleaseFile -ReleaseVersion $releaseVersion -Name "checksums.txt" -Destination $checksumsPath
        Copy-SetupEnvReleaseFile -ReleaseVersion $releaseVersion -Name $assetName -Destination $archivePath
        $expected = Get-SetupEnvChecksum -ChecksumsPath $checksumsPath -AssetName $assetName
        $actual = (Get-FileHash -LiteralPath $archivePath -Algorithm SHA256).Hash.ToLowerInvariant()
        if ($actual -ne $expected) {
            throw "Checksum verification failed for '$assetName'. Expected $expected but received $actual. The archive will not be executed."
        }
        $extractDir = Join-Path $temporary "extract"
        Expand-Archive -LiteralPath $archivePath -DestinationPath $extractDir -Force
        $candidate = Join-Path $extractDir "setup-env.exe"
        if (-not (Test-Path -LiteralPath $candidate -PathType Leaf)) {
            throw "Verified archive does not contain setup-env.exe."
        }
        $candidateOutput = & $candidate version 2>&1 | Out-String
        if ($LASTEXITCODE -ne 0 -or -not $candidateOutput.Contains($releaseVersion)) {
            throw "The verified binary failed its version check. Output: $candidateOutput"
        }

        New-Item -ItemType Directory -Path $TargetDir -Force | Out-Null
        $destination = Join-Path $TargetDir "setup-env.exe"
        $staged = Join-Path $TargetDir ".setup-env.new.$PID.exe"
        $backup = "$destination.backup"
        $oldVersion = "not installed"
        if (Test-Path -LiteralPath $destination -PathType Leaf) {
            $oldOutput = & $destination version 2>&1 | Out-String
            if ($LASTEXITCODE -eq 0) {
                $oldVersion = $oldOutput.Trim()
            }
            Copy-Item -LiteralPath $destination -Destination $backup -Force
        }
        Copy-Item -LiteralPath $candidate -Destination $staged -Force
        try {
            Move-Item -LiteralPath $staged -Destination $destination -Force
            $installedOutput = & $destination version 2>&1 | Out-String
            if ($LASTEXITCODE -ne 0 -or -not $installedOutput.Contains($releaseVersion)) {
                throw "Installed binary failed its version check."
            }
        } catch {
            if (Test-Path -LiteralPath $backup -PathType Leaf) {
                Copy-Item -LiteralPath $backup -Destination $destination -Force
            }
            throw "Installation replacement failed and the previous binary was restored where available. $($_.Exception.Message)"
        } finally {
            if (Test-Path -LiteralPath $staged) {
                Remove-Item -LiteralPath $staged -Force
            }
        }
        if (Test-Path -LiteralPath $backup) {
            Remove-Item -LiteralPath $backup -Force
        }

        $statePath = Join-Path $TargetDir ".setup-env-install.json"
        $managedPath = $false
        if (Test-Path -LiteralPath $statePath -PathType Leaf) {
            try {
                $previousState = Get-Content -Raw -LiteralPath $statePath | ConvertFrom-Json
                $managedPath = $previousState.managed_path -eq $true
            } catch {
                $managedPath = $false
            }
        }
        if (-not $SkipPath) {
            $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
            $updated = Add-SetupEnvPathEntry -PathValue $userPath -Entry $TargetDir
            if ($updated -ne $userPath) {
                try {
                    [Environment]::SetEnvironmentVariable("Path", $updated, "User")
                    $managedPath = $true
                } catch {
                    Write-Warning "Unable to update current-user PATH. Add '$TargetDir' manually. $($_.Exception.Message)"
                }
            }
            if (-not (Test-SetupEnvPathEntry -PathValue $env:Path -Entry $TargetDir)) {
                $env:Path = Add-SetupEnvPathEntry -PathValue $env:Path -Entry $TargetDir
            }
        }
        $state = [ordered]@{
            installed_version = $releaseVersion
            install_directory = $TargetDir
            installer_version = $script:InstallerVersion
            installed_at = [DateTime]::UtcNow.ToString("o")
            asset_name = $assetName
            checksum_sha256 = $expected
            managed_path = $managedPath
        }
        $state | ConvertTo-Json | Set-Content -LiteralPath $statePath -Encoding UTF8
        $verb = if ($IsUpgrade -or $oldVersion -ne "not installed") { "upgraded" } else { "installed" }
        Write-Host "Setup Env $releaseVersion was $verb successfully."
        Write-Host "Previous: $oldVersion"
        Write-Host "Installed: $destination"
        if (-not $SkipPath) {
            Write-Host "Restart existing terminal sessions if 'setup-env' is not immediately available."
        }
    } finally {
        if (Test-Path -LiteralPath $temporary) {
            Remove-Item -LiteralPath $temporary -Recurse -Force
        }
    }
}

function Invoke-SetupEnvInstaller {
    param(
        [string]$RequestedVersion = "latest",
        [string]$RequestedInstallDir,
        [switch]$IsUpgrade,
        [switch]$IsUninstall,
        [switch]$RemoveData,
        [switch]$AssumeYes,
        [switch]$SkipPath
    )
    if ($IsUpgrade -and $IsUninstall) {
        throw "-Upgrade and -Uninstall cannot be used together."
    }
    if ($RemoveData -and -not $IsUninstall) {
        throw "-Purge is valid only with -Uninstall."
    }
    $targetDir = Resolve-SetupEnvInstallDir -Value $RequestedInstallDir
    if ($IsUninstall) {
        Invoke-SetupEnvUninstall -TargetDir $targetDir -RemoveData:$RemoveData -AssumeYes:$AssumeYes
        return
    }
    Invoke-SetupEnvInstall -RequestedVersion $RequestedVersion -TargetDir $targetDir -IsUpgrade:$IsUpgrade -SkipPath:$SkipPath
}

if ($MyInvocation.InvocationName -ne ".") {
    try {
        Invoke-SetupEnvInstaller `
            -RequestedVersion $Version `
            -RequestedInstallDir $InstallDir `
            -IsUpgrade:$Upgrade `
            -IsUninstall:$Uninstall `
            -RemoveData:$Purge `
            -AssumeYes:$Yes `
            -SkipPath:$NoPath
    } catch {
        Write-Error "Setup Env installer failed: $($_.Exception.Message)"
        exit 1
    }
}
