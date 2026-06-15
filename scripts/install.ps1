# Install claude-cli from GitHub Releases
# Usage: irm https://raw.githubusercontent.com/myflavor/claude-cli/main/scripts/install.ps1 | iex
$ErrorActionPreference = "Stop"

$Repo = "myflavor/claude-cli"
$BinaryName = "claude-cli.exe"

# Detect arch
$Arch = $env:PROCESSOR_ARCHITECTURE
switch ($Arch) {
    "AMD64" { $Asset = "claude-cli.exe" }
    "ARM64" { $Asset = "claude-cli-arm64.exe" }
    default { Write-Error "Unsupported architecture: $Arch"; exit 1 }
}

$DownloadUrl = "https://github.com/$Repo/releases/latest/download/$Asset"

# Find claude location
$ClaudeCmd = Get-Command claude -ErrorAction SilentlyContinue
if (-not $ClaudeCmd) {
    Write-Error "Error: claude command not found in PATH"
    Write-Error "Please install Claude Code first: https://code.claude.com"
    exit 1
}

$ClaudeDir = Split-Path -Parent $ClaudeCmd.Source
$Target = Join-Path $ClaudeDir $BinaryName

# Check if already installed
if (Test-Path $Target) {
    Write-Host "Existing claude-cli found at: $Target"
    Write-Host "Updating to latest..."
}

# Download to temp file first
$TmpFile = [System.IO.Path]::GetTempFileName()
try {
    Write-Host "Downloading $Asset from $Repo..."
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $TmpFile -UseBasicParsing

    if ((Get-Item $TmpFile).Length -eq 0) {
        throw "Downloaded file is empty"
    }

    # Atomic move (overwrite)
    Move-Item -Force $TmpFile $Target

    Write-Host "Installed to: $Target"
    Write-Host "Test it with: claude-cli start -P claude"
}
catch {
    Remove-Item -Force $TmpFile -ErrorAction SilentlyContinue
    Write-Error "Error: $_"
    exit 1
}
