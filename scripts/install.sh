#!/bin/bash
# Install claude-st and claude-cfg from GitHub Releases
# Usage: curl -sSL https://raw.githubusercontent.com/myflavor/claude-cli/main/scripts/install.sh | bash
set -e

REPO="myflavor/claude-cli"

# Detect OS and arch, map to release asset names
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS-$ARCH" in
    linux-x86_64)  ST="claude-st";             CFG="claude-cfg" ;;
    linux-aarch64) ST="claude-st-linux-arm64"; CFG="claude-cfg-linux-arm64" ;;
    darwin-x86_64) ST="claude-st-darwin-amd64"; CFG="claude-cfg-darwin-amd64" ;;
    darwin-arm64)  ST="claude-st-darwin-arm64"; CFG="claude-cfg-darwin-arm64" ;;
    *) echo "Unsupported platform: $OS-$ARCH"; exit 1 ;;
esac

# Find claude location (both binaries install next to it)
CLAUDE_PATH=$(command -v claude)
if [ -z "$CLAUDE_PATH" ]; then
    echo "Error: claude command not found in PATH"
    echo "Please install Claude Code first: https://code.claude.com"
    exit 1
fi

CLAUDE_DIR=$(cd "$(dirname "$CLAUDE_PATH")" && pwd)
ST_TARGET="$CLAUDE_DIR/claude-st"
CFG_TARGET="$CLAUDE_DIR/claude-cfg"

# Check write permission on any existing binary (e.g. previously installed with sudo)
check_writable() {
    if [ -f "$1" ] && [ ! -w "$1" ]; then
        echo "Error: no write permission to $1"
        echo "Try running with sudo, or install claude to a user-writable location"
        exit 1
    fi
}
check_writable "$ST_TARGET"
check_writable "$CFG_TARGET"

TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

download() {
    local asset="$1"
    local out="$2"
    echo "Downloading $asset from $REPO..."
    if ! curl -sSL -f -o "$out" "https://github.com/$REPO/releases/latest/download/$asset"; then
        echo "Error: download failed for $asset"
        exit 1
    fi
    if [ ! -s "$out" ]; then
        echo "Error: downloaded file is empty: $asset"
        exit 1
    fi
    chmod +x "$out"
}

download "$ST"  "$TMP_DIR/claude-st"
download "$CFG" "$TMP_DIR/claude-cfg"

# Atomic move into place
mv -f "$TMP_DIR/claude-st"  "$ST_TARGET"
mv -f "$TMP_DIR/claude-cfg" "$CFG_TARGET"

# Set up ~/.claude-cli/ directory
HOME_DIR="${HOME:-$(eval echo ~$(id -un))}"
PROVIDERS_DIR="$HOME_DIR/.claude-cli"

if [ ! -d "$PROVIDERS_DIR" ]; then
    mkdir -p "$PROVIDERS_DIR"
    echo ""
    echo "Created providers directory: $PROVIDERS_DIR"
    echo "Add a provider by placing a settings JSON file at: $PROVIDERS_DIR/<name>.json"
fi

echo ""
echo "Installed to: $ST_TARGET"
echo "              $CFG_TARGET"
echo "Test it with: claude-st <provider-name>"
