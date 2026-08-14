#!/bin/sh
# install.sh — install lokan from GitHub releases. No Go, no build.
# Usage: curl -fsSL https://raw.githubusercontent.com/avressatelier/lokan/main/install.sh | sh
set -e

BIN="lokan"
REPO="avressatelier/lokan"
DEST="${LOKAN_INSTALL_DIR:-$HOME/.local/bin}"

# detect OS + arch and map to goreleaser asset names
os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"
case "$arch" in
  x86_64 | amd64) arch="amd64" ;;
  aarch64 | arm64) arch="arm64" ;;
  *) echo "unsupported architecture: $arch" >&2; exit 1 ;;
esac
case "$os" in
  darwin | linux) ;;
  *) echo "unsupported OS: $os" >&2; exit 1 ;;
esac

url="https://github.com/$REPO/releases/latest/download/${BIN}_${os}_${arch}.tar.gz"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

# download + extract only the binary
curl -fsSL "$url" | tar -xz -C "$tmp" "$BIN"

# install
mkdir -p "$DEST"
install -m 0755 "$tmp/$BIN" "$DEST/$BIN"

# warn if the destination isn't on PATH
case ":$PATH:" in
  *":$DEST:"*) ;;
  *)
    echo "warning: $DEST is not on your PATH" >&2
    echo "  add it with: echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.zshrc" >&2
    ;;
esac

echo "installed $BIN to $DEST/$BIN"
echo "try it: lokan init board.md"
