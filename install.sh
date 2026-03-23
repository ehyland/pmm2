#!/bin/bash
set -e

REPO="ehyland/pmm2"
PMM2_DIR="$HOME/.pmm2"
INSTALL_DIR="$PMM2_DIR/bin"

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
RELEASE_JSON=$(curl -fsSL -H "Accept: application/vnd.github+json" -H "User-Agent: pmm2-installer" "https://api.github.com/repos/$REPO/releases/latest")

if ! command -v node >/dev/null 2>&1; then
    echo "Error: Node.js is required to extract tag_name from https://api.github.com/repos/$REPO/releases/latest"
    echo "Install Node.js and retry. https://nodejs.org/en/download"
    exit 1
fi

LATEST_RELEASE=$(RELEASE_JSON="$RELEASE_JSON" node -e 'try { const release = JSON.parse(process.env.RELEASE_JSON || "{}"); if (release.tag_name) process.stdout.write(release.tag_name); } catch {}' 2>/dev/null || true)

if [ -z "$LATEST_RELEASE" ]; then
    echo "Error: Could not extract tag_name from https://api.github.com/repos/$REPO/releases/latest"
    exit 1
fi

FILENAME="pmm2_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${LATEST_RELEASE}/${FILENAME}"

echo "Downloading $URL..."
TMP_DIR=$(mktemp -d)
curl -L "$URL" -o "${TMP_DIR}/${FILENAME}"

echo "Installing to $INSTALL_DIR..."
tar -xzf "${TMP_DIR}/${FILENAME}" -C "${TMP_DIR}"

mv "${TMP_DIR}/pmm2" "$INSTALL_DIR/pmm2"
chmod +x "$INSTALL_DIR/pmm2"

echo "Creating symlinks..."
"$INSTALL_DIR/pmm2" setup

echo "pmm2 successfully installed to $INSTALL_DIR"

echo ""
echo ""
echo "Restart your terminal or run: source ~/.bashrc"
echo ""


