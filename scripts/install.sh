#!/bin/bash
# Install claude-st and claude-cfg from GitHub Releases
# Usage: curl -sSL https://raw.githubusercontent.com/myflavor/claude-cli/main/scripts/install.sh | bash
#
# Downloads a single per-platform tarball (claude-cli-<os>-<arch>.tar.gz)
# containing both binaries and extracts them next to claude. Uses only tar,
# which is built in on Linux and macOS.
set -e

REPO="myflavor/claude-cli"

# Detect OS and arch, map to the release tarball name
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)  ARCH=amd64 ;;
    aarch64|arm64) ARCH=arm64 ;;
    *) echo "Unsupported arch: $ARCH"; exit 1 ;;
esac

case "$OS-$ARCH" in
    linux-amd64)   TGZ="claude-cli-linux-amd64.tar.gz" ;;
    linux-arm64)   TGZ="claude-cli-linux-arm64.tar.gz" ;;
    darwin-amd64)  TGZ="claude-cli-darwin-amd64.tar.gz" ;;
    darwin-arm64)  TGZ="claude-cli-darwin-arm64.tar.gz" ;;
    *) echo "Unsupported platform: $OS-$ARCH"; exit 1 ;;
esac

# tar is required to extract the archive (built in on Linux/macOS)
if ! command -v tar >/dev/null 2>&1; then
    echo "Error: 'tar' is required but not found (very unusual on Linux/macOS)."
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

echo "Downloading $TGZ from $REPO..."
if ! curl -sSL -f -o "$TMP_DIR/$TGZ" "https://github.com/$REPO/releases/latest/download/$TGZ"; then
    echo "Error: download failed for $TGZ"
    exit 1
fi
if [ ! -s "$TMP_DIR/$TGZ" ]; then
    echo "Error: downloaded file is empty: $TGZ"
    exit 1
fi

tar -xzf "$TMP_DIR/$TGZ" -C "$TMP_DIR"
chmod +x "$TMP_DIR/claude-st" "$TMP_DIR/claude-cfg"

# Atomic install: stage in the target dir, then rename (same filesystem).
install_bin() {
    src="$1"; dst="$2"; tmp="$dst.new.$$"
    cp -f "$src" "$tmp" || { rm -f "$tmp"; echo "Error: failed to stage $dst"; exit 1; }
    mv -f "$tmp" "$dst" || { rm -f "$tmp"; echo "Error: failed to install $dst"; exit 1; }
}
install_bin "$TMP_DIR/claude-st"  "$ST_TARGET"
install_bin "$TMP_DIR/claude-cfg" "$CFG_TARGET"

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
