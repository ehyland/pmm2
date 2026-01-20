#!/bin/bash
set -e

REPO="ehyland/pmm2"
PMM_DIR="$HOME/.pmm2"
INSTALL_DIR="$PMM_DIR/bin"

# Create install directory
mkdir -p "$INSTALL_DIR"

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64) ARCH="x86_64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
    darwin) OS="Darwin" ;;
    linux) OS="Linux" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

echo "Detecting latest version..."
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo "Error: Could not find latest release for $REPO"
    exit 1
fi

FILENAME="pmm_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${LATEST_RELEASE}/${FILENAME}"

echo "Downloading $URL..."
TMP_DIR=$(mktemp -d)
curl -L "$URL" -o "${TMP_DIR}/${FILENAME}"

echo "Installing to $INSTALL_DIR..."
tar -xzf "${TMP_DIR}/${FILENAME}" -C "${TMP_DIR}"

mv "${TMP_DIR}/pmm" "$INSTALL_DIR/pmm"
chmod +x "$INSTALL_DIR/pmm"

echo "Creating symlinks..."
"$INSTALL_DIR/pmm" setup

echo "pmm successfully installed to $INSTALL_DIR"

echo ""
echo ""
echo "Please add the following to your .bashrc (or .zshrc):"
echo ""
echo "  export PATH=\"\$HOME/.pmm2/bin:\$PATH\""
echo ""
echo "After adding, restart your terminal or run: source ~/.bashrc"
echo ""


