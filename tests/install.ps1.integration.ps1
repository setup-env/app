param(
    [Parameter(Mandatory = $true)][string]$ReleaseDir,
    [string]$Version = "v0.1.0",
    [string]$UpgradeReleaseDir,
    [string]$UpgradeVersion
)

$ErrorActionPreference = "Stop"
$TestVersion = $Version
$TestUpgradeVersion = if ($UpgradeVersion) { $UpgradeVersion } else { $Version }
. (Join-Path $PSScriptRoot "..\install.ps1")

$root = Join-Path ([IO.Path]::GetTempPath()) ("setup-env-ps-integration-" + [Guid]::NewGuid().ToString("N"))
$installDir = Join-Path $root "install"
$fixtureDir = Join-Path $root "fixtures"
$upgradeFixtureDir = Join-Path $root "upgrade-fixtures"
$scratchDir = Join-Path $root "scratch"
New-Item -ItemType Directory -Path $root, $fixtureDir, $upgradeFixtureDir, $scratchDir | Out-Null
$originalFixture = $env:SETUP_ENV_RELEASE_DIR
$originalTestVersion = $env:SETUP_ENV_TEST_VERSION
$originalTemp = $env:TEMP
$originalTmp = $env:TMP
try {
    Get-ChildItem -LiteralPath $ReleaseDir -File | Copy-Item -Destination $fixtureDir -Force
    if ($UpgradeReleaseDir) {
        Get-ChildItem -LiteralPath $UpgradeReleaseDir -File | Copy-Item -Destination $upgradeFixtureDir -Force
    } else {
        Get-ChildItem -LiteralPath $ReleaseDir -File | Copy-Item -Destination $upgradeFixtureDir -Force
    }
    $env:SETUP_ENV_RELEASE_DIR = $fixtureDir
    $env:SETUP_ENV_TEST_VERSION = $TestVersion
    $env:TEMP = $scratchDir
    $env:TMP = $scratchDir

    Invoke-SetupEnvInstaller -RequestedVersion $TestVersion -RequestedInstallDir $installDir -AssumeYes -SkipPath
    $binary = Join-Path $installDir "setup-env.exe"
    $state = Join-Path $installDir ".setup-env-install.json"
    if (-not (Test-Path -LiteralPath $binary) -or -not (Test-Path -LiteralPath $state)) {
        throw "Installer did not create the binary and metadata."
    }
    $versionOutput = & $binary version | Out-String
    if ($LASTEXITCODE -ne 0 -or -not $versionOutput.Contains($TestVersion)) {
        throw "Installed binary metadata is incorrect: $versionOutput"
    }

    Set-Content -LiteralPath (Join-Path $installDir "unrelated.txt") -Value "preserve"
    $env:SETUP_ENV_RELEASE_DIR = $upgradeFixtureDir
    $env:SETUP_ENV_TEST_VERSION = $TestUpgradeVersion
    Invoke-SetupEnvInstaller -RequestedVersion $TestUpgradeVersion -RequestedInstallDir $installDir -IsUpgrade -AssumeYes -SkipPath
    $upgradedOutput = & $binary version | Out-String
    if ($LASTEXITCODE -ne 0 -or -not $upgradedOutput.Contains($TestUpgradeVersion)) {
        throw "Upgraded binary metadata is incorrect: $upgradedOutput"
    }
    Invoke-SetupEnvInstaller -RequestedInstallDir $installDir -IsUninstall -AssumeYes
    if ((Test-Path -LiteralPath $binary) -or (Test-Path -LiteralPath $state)) {
        throw "Uninstall left installer-owned files behind."
    }
    if (-not (Test-Path -LiteralPath (Join-Path $installDir "unrelated.txt"))) {
        throw "Uninstall removed an unrelated file."
    }

    $architecture = Get-SetupEnvArchitecture
    $env:SETUP_ENV_RELEASE_DIR = $fixtureDir
    $env:SETUP_ENV_TEST_VERSION = $TestVersion
    $asset = Get-SetupEnvAssetName -ReleaseVersion $TestVersion -Architecture $architecture
    Add-Content -LiteralPath (Join-Path $fixtureDir $asset) -Value "corrupt"
    $failed = $false
    try {
        Invoke-SetupEnvInstaller -RequestedVersion $TestVersion -RequestedInstallDir (Join-Path $root "corrupt-install") -AssumeYes -SkipPath
    } catch {
        $failed = $_.Exception.Message.Contains("Checksum verification failed")
    }
    if (-not $failed) {
        throw "Corrupt archive was not rejected by checksum verification."
    }
    if (Test-Path -LiteralPath (Join-Path $root "corrupt-install\setup-env.exe")) {
        throw "Corrupt archive installed an executable."
    }
    $leftovers = @(Get-ChildItem -LiteralPath $scratchDir -Directory -Filter "setup-env-install-*")
    if ($leftovers.Count -ne 0) {
        throw "Installer temporary directory cleanup is incomplete."
    }
} finally {
    $env:SETUP_ENV_RELEASE_DIR = $originalFixture
    $env:SETUP_ENV_TEST_VERSION = $originalTestVersion
    $env:TEMP = $originalTemp
    $env:TMP = $originalTmp
    if (Test-Path -LiteralPath $root) {
        Remove-Item -LiteralPath $root -Recurse -Force
    }
}

Write-Host "PowerShell installer integration tests passed."
