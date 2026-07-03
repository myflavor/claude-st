# Install claude-st and claude-cfg from GitHub Releases
# Usage: irm https://raw.githubusercontent.com/myflavor/claude-cli/main/scripts/install.ps1 | iex
#
# Downloads a single per-platform zip (claude-cli-windows-<arch>.zip) containing
# both binaries and extracts them next to claude.
$ErrorActionPreference = "Stop"

$Repo = "myflavor/claude-cli"

# Detect arch, map to the release zip name
$Arch = $env:PROCESSOR_ARCHITECTURE
switch ($Arch) {
    "AMD64" { $Zip = "claude-cli-windows-amd64.zip" }
    "ARM64" { $Zip = "claude-cli-windows-arm64.zip" }
    default { Write-Error "Unsupported architecture: $Arch"; exit 1 }
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

$TmpZip = [System.IO.Path]::GetTempFileName() + ".zip"
$ExtractDir = Join-Path ([System.IO.Path]::GetTempPath()) ("claude-cli-" + [System.IO.Path]::GetRandomFileName())
New-Item -ItemType Directory -Path $ExtractDir -Force | Out-Null
try {
    $Url = "https://github.com/$Repo/releases/latest/download/$Zip"
    Write-Host "Downloading $Zip from $Repo..."
    Invoke-WebRequest -Uri $Url -OutFile $TmpZip -UseBasicParsing
    if ((Get-Item $TmpZip).Length -eq 0) {
        throw "Downloaded file is empty: $Zip"
    }

    Expand-Archive -Path $TmpZip -DestinationPath $ExtractDir -Force

    # Atomic move (overwrite)
    Move-Item -Force (Join-Path $ExtractDir "claude-st.exe") $StTarget
    Move-Item -Force (Join-Path $ExtractDir "claude-cfg.exe") $CfgTarget

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
    Remove-Item -Force $TmpZip -ErrorAction SilentlyContinue
    Remove-Item -Recurse -Force $ExtractDir -ErrorAction SilentlyContinue
}
