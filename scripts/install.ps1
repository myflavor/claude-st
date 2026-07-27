# Install claude-st and claude-cfg from GitHub Releases
# Usage: irm https://raw.githubusercontent.com/myflavor/claude-cli/main/scripts/install.ps1 | iex
#
# Downloads a single per-platform tarball (claude-cli-windows-<arch>.tar.gz)
# containing both binaries and extracts them next to claude. Uses tar.exe,
# which is built into Windows 10+.
$ErrorActionPreference = "Stop"

$Repo = "myflavor/claude-cli"

# Detect arch, map to the release tarball name
$Arch = $env:PROCESSOR_ARCHITECTURE
switch ($Arch) {
    "AMD64" { $Tgz = "claude-cli-windows-amd64.tar.gz" }
    "ARM64" { $Tgz = "claude-cli-windows-arm64.tar.gz" }
    default { Write-Error "Unsupported architecture: $Arch"; exit 1 }
}

# tar.exe ships with Windows 10 1803+; required to extract the tarball
if (-not (Get-Command tar -ErrorAction SilentlyContinue)) {
    Write-Error "Error: 'tar' not found. It is built into Windows 10 1803 and newer."
    exit 1
}

# Find claude location (both binaries install next to it)
$ClaudeCmd = Get-Command claude -ErrorAction SilentlyContinue
if (-not $ClaudeCmd) {
    Write-Error "Error: claude command not found in PATH"
    Write-Error "Please install Claude Code first: https://code.claude.com"
    exit 1
}

$ClaudeDir = Split-Path -Parent $ClaudeCmd.Source
$StTarget = Join-Path $ClaudeDir "claude-st.exe"
$CfgTarget = Join-Path $ClaudeDir "claude-cfg.exe"

$TmpTgz = [System.IO.Path]::GetTempFileName() + ".tar.gz"
$ExtractDir = Join-Path ([System.IO.Path]::GetTempPath()) ("claude-cli-" + [System.IO.Path]::GetRandomFileName())
New-Item -ItemType Directory -Path $ExtractDir -Force | Out-Null
try {
    $Url = "https://github.com/$Repo/releases/latest/download/$Tgz"
    Write-Host "Downloading $Tgz from $Repo..."
    Invoke-WebRequest -Uri $Url -OutFile $TmpTgz -UseBasicParsing
    if ((Get-Item $TmpTgz).Length -eq 0) {
        throw "Downloaded file is empty: $Tgz"
    }

    tar -xzf $TmpTgz -C $ExtractDir

    # Atomic install: stage in the target dir, then rename (same volume).
    $StStaged = Join-Path $ClaudeDir ".claude-st.exe.tmp"
    $CfgStaged = Join-Path $ClaudeDir ".claude-cfg.exe.tmp"
    Copy-Item -Force (Join-Path $ExtractDir "claude-st.exe") $StStaged
    Copy-Item -Force (Join-Path $ExtractDir "claude-cfg.exe") $CfgStaged
    Move-Item -Force $StStaged $StTarget
    Move-Item -Force $CfgStaged $CfgTarget

    # Set up ~/.claude-cli/ directory
    $HomeDir = $env:USERPROFILE
    $ProvidersDir = Join-Path $HomeDir ".claude-cli"

    if (-not (Test-Path $ProvidersDir)) {
        New-Item -ItemType Directory -Path $ProvidersDir -Force | Out-Null
        Write-Host ""
        Write-Host "Created providers directory: $ProvidersDir"
        Write-Host "Add a provider by placing a settings JSON file at: $ProvidersDir\<name>.json"
    }

    Write-Host ""
    Write-Host "Installed to: $StTarget"
    Write-Host "              $CfgTarget"
    Write-Host "Test it with: claude-st <provider-name>"
}
catch {
    Write-Error "Error: $_"
    exit 1
}
finally {
    Remove-Item -Force $TmpTgz -ErrorAction SilentlyContinue
    Remove-Item -Recurse -Force $ExtractDir -ErrorAction SilentlyContinue
    foreach ($s in @($StStaged, $CfgStaged)) {
        if ($s -and (Test-Path $s)) { Remove-Item -Force $s -ErrorAction SilentlyContinue }
    }
}
