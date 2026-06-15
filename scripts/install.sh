#!/bin/bash
# Install claude-cli from GitHub Releases
# Usage: curl -sSL https://raw.githubusercontent.com/myflavor/claude-cli/main/scripts/install.sh | bash
set -e

REPO="myflavor/claude-cli"
BINARY_NAME="claude-cli"

# Detect OS and arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS-$ARCH" in
    linux-x86_64)  ASSET="claude-cli" ;;
    linux-aarch64) ASSET="claude-cli-linux-arm64" ;;
    darwin-x86_64) ASSET="claude-cli-darwin-amd64" ;;
    darwin-arm64)  ASSET="claude-cli-darwin-arm64" ;;
    *) echo "Unsupported platform: $OS-$ARCH"; exit 1 ;;
esac

DOWNLOAD_URL="https://github.com/$REPO/releases/latest/download/$ASSET"

# Find claude location
CLAUDE_PATH=$(command -v claude)
if [ -z "$CLAUDE_PATH" ]; then
    echo "Error: claude command not found in PATH"
    echo "Please install Claude Code first: https://code.claude.com"
    exit 1
fi

CLAUDE_DIR=$(cd "$(dirname "$CLAUDE_PATH")" && pwd)
TARGET="$CLAUDE_DIR/$BINARY_NAME"

# Check if already installed
if [ -f "$TARGET" ]; then
    echo "Existing claude-cli found at: $TARGET"
    echo "Updating to latest..."
fi

# Check write permission
if [ -f "$TARGET" ] && [ ! -w "$TARGET" ]; then
    echo "Error: no write permission to $TARGET"
    echo "Try running with sudo, or install claude to a user-writable location"
    exit 1
fi

# Download to temp file first, then move (atomic)
TMP_FILE=$(mktemp)
trap "rm -f $TMP_FILE" EXIT

echo "Downloading $ASSET from $REPO..."
if ! curl -sSL -f -o "$TMP_FILE" "$DOWNLOAD_URL"; then
    echo "Error: download failed"
    exit 1
fi

# Verify the download looks like a binary
if [ ! -s "$TMP_FILE" ]; then
    echo "Error: downloaded file is empty"
    exit 1
fi

chmod +x "$TMP_FILE"
mv -f "$TMP_FILE" "$TARGET"

echo "Installed to: $TARGET"
echo "Test it with: $BINARY_NAME start -P claude"
