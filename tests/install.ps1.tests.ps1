$ErrorActionPreference = "Stop"

. (Join-Path $PSScriptRoot "..\install.ps1")

function Assert-Equal {
    param($Actual, $Expected, [string]$Message)
    if ($Actual -ne $Expected) {
        throw "$Message. Actual: '$Actual'; expected: '$Expected'."
    }
}

function Assert-Throws {
    param([scriptblock]$Action, [string]$Message)
    try {
        & $Action
    } catch {
        return
    }
    throw "$Message. Expected an exception."
}

Assert-Equal (Normalize-SetupEnvVersion "0.1.0") "v0.1.0" "Version normalization failed"
Assert-Equal (Normalize-SetupEnvVersion "v1.2.3-rc.1") "v1.2.3-rc.1" "Prerelease normalization failed"
Assert-Throws { Normalize-SetupEnvVersion "latest-ish" } "Invalid version was accepted"
Assert-Throws { Normalize-SetupEnvVersion "v1.2.3+build.4" } "Build-metadata version was accepted"
Assert-Equal (Get-SetupEnvArchitecture "AMD64") "amd64" "AMD64 mapping failed"
Assert-Equal (Get-SetupEnvArchitecture "ARM64") "arm64" "ARM64 mapping failed"
Assert-Throws { Get-SetupEnvArchitecture "x86" } "Unsupported architecture was accepted"
Assert-Equal `
    (Get-SetupEnvAssetName -ReleaseVersion "v0.1.0" -Architecture "arm64") `
    "setup-env_0.1.0_windows_arm64.zip" `
    "Asset name is incorrect"

$pathValue = "C:\Tools;C:\Users\example\bin"
Assert-Equal `
    (Add-SetupEnvPathEntry -PathValue $pathValue -Entry "C:\Users\example\bin\") `
    $pathValue `
    "PATH entry was duplicated"
$addedPath = Add-SetupEnvPathEntry -PathValue $pathValue -Entry "C:\Setup Env\bin"
if (-not (Test-SetupEnvPathEntry -PathValue $addedPath -Entry "C:\Setup Env\bin")) {
    throw "PATH entry was not added"
}
Assert-Equal `
    (Remove-SetupEnvPathEntry -PathValue $addedPath -Entry "c:\setup env\bin\") `
    $pathValue `
    "PATH entry was not removed safely"

$temporary = Join-Path ([IO.Path]::GetTempPath()) ("setup-env-ps-unit-" + [Guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $temporary | Out-Null
try {
    $checksums = Join-Path $temporary "checksums.txt"
    $hash = "a" * 64
    Set-Content -LiteralPath $checksums -Value "$hash  setup-env_0.1.0_windows_amd64.zip"
    Assert-Equal `
        (Get-SetupEnvChecksum -ChecksumsPath $checksums -AssetName "setup-env_0.1.0_windows_amd64.zip") `
        $hash `
        "Checksum parsing failed"
    Assert-Throws {
        Get-SetupEnvChecksum -ChecksumsPath $checksums -AssetName "missing.zip"
    } "Missing checksum was accepted"

    $relative = Join-Path $temporary "..\custom-bin"
    Assert-Equal `
        (Resolve-SetupEnvInstallDir -Value $relative) `
        ([IO.Path]::GetFullPath($relative)) `
        "Custom install directory was not normalized"

    $oldFixture = $env:SETUP_ENV_RELEASE_DIR
    $env:SETUP_ENV_RELEASE_DIR = $temporary
    Assert-Throws {
        Copy-SetupEnvReleaseFile -ReleaseVersion "v9.9.9" -Name "missing.zip" -Destination (Join-Path $temporary "out.zip")
    } "Unavailable local release did not fail"
    $env:SETUP_ENV_RELEASE_DIR = $oldFixture
} finally {
    Remove-Item -LiteralPath $temporary -Recurse -Force
}

Write-Host "PowerShell installer unit tests passed."
