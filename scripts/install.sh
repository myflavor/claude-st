#!/bin/bash
# Install claude-st and claude-cfg from GitHub Releases
# Usage: curl -sSL https://raw.githubusercontent.com/myflavor/claude-cli/main/scripts/install.sh | bash
#
# Downloads a single per-platform zip (claude-cli-<os>-<arch>.zip) containing
# both binaries and extracts them next to claude.
set -e

REPO="myflavor/claude-cli"

# Detect OS and arch, map to the release zip name
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)  ARCH=amd64 ;;
    aarch64) ARCH=arm64 ;;
    *) echo "Unsupported arch: $ARCH"; exit 1 ;;
esac

case "$OS-$ARCH" in
    linux-amd64)   ZIP="claude-cli-linux-amd64.zip" ;;
    linux-arm64)   ZIP="claude-cli-linux-arm64.zip" ;;
    darwin-amd64)  ZIP="claude-cli-darwin-amd64.zip" ;;
    darwin-arm64)  ZIP="claude-cli-darwin-arm64.zip" ;;
    *) echo "Unsupported platform: $OS-$ARCH"; exit 1 ;;
esac

# unzip is required to extract the archive
if ! command -v unzip >/dev/null 2>&1; then
    echo "Error: 'unzip' is required but not installed."
    echo "Install it first, e.g.: sudo apt install unzip   /   brew install unzip"
    exit 1
fi

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

echo "Downloading $ZIP from $REPO..."
if ! curl -sSL -f -o "$TMP_DIR/$ZIP" "https://github.com/$REPO/releases/latest/download/$ZIP"; then
    echo "Error: download failed for $ZIP"
    exit 1
fi
if [ ! -s "$TMP_DIR/$ZIP" ]; then
    echo "Error: downloaded file is empty: $ZIP"
    exit 1
fi

unzip -o "$TMP_DIR/$ZIP" -d "$TMP_DIR" >/dev/null
chmod +x "$TMP_DIR/claude-st" "$TMP_DIR/claude-cfg"

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
