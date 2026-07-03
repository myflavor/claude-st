# Install claude-st and claude-cfg from GitHub Releases
# Usage: irm https://raw.githubusercontent.com/myflavor/claude-cli/main/scripts/install.ps1 | iex
$ErrorActionPreference = "Stop"

$Repo = "myflavor/claude-cli"

# Detect arch, map to release asset names
$Arch = $env:PROCESSOR_ARCHITECTURE
switch ($Arch) {
    "AMD64" { $St = "claude-st.exe";       $Cfg = "claude-cfg.exe" }
    "ARM64" { $St = "claude-st-arm64.exe"; $Cfg = "claude-cfg-arm64.exe" }
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

$TmpSt = [System.IO.Path]::GetTempFileName()
$TmpCfg = [System.IO.Path]::GetTempFileName()
try {
    function Download-Binary($Asset, $OutFile) {
        $Url = "https://github.com/$Repo/releases/latest/download/$Asset"
        Write-Host "Downloading $Asset from $Repo..."
        Invoke-WebRequest -Uri $Url -OutFile $OutFile -UseBasicParsing
        if ((Get-Item $OutFile).Length -eq 0) {
            throw "Downloaded file is empty: $Asset"
        }
    }

    Download-Binary $St $TmpSt
    Download-Binary $Cfg $TmpCfg

    # Atomic move (overwrite)
    Move-Item -Force $TmpSt $StTarget
    Move-Item -Force $TmpCfg $CfgTarget

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
    Remove-Item -Force $TmpSt -ErrorAction SilentlyContinue
    Remove-Item -Force $TmpCfg -ErrorAction SilentlyContinue
    Write-Error "Error: $_"
    exit 1
}
