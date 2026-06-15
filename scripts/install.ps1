# Install claude-cli from GitHub Releases
# Usage: irm https://raw.githubusercontent.com/myflavor/claude-cli/main/scripts/install.ps1 | iex
$ErrorActionPreference = "Stop"

$Repo = "myflavor/claude-cli"
$BinaryName = "claude-cli.exe"

# Detect arch (Windows is always amd64 or arm64)
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

Write-Host "Downloading $Asset from $Repo..."
Invoke-WebRequest -Uri $DownloadUrl -OutFile $Target -UseBasicParsing

Write-Host "Installed to: $Target"
Write-Host "Test it with: claude-cli start -P claude"
