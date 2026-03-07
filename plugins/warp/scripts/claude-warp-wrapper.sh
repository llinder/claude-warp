#!/usr/bin/env bash
set -euo pipefail

REPO="llinder/claude-warp"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BINARY="$SCRIPT_DIR/claude-warp"
VERSION_FILE="$SCRIPT_DIR/.version"

# Determine OS and architecture
get_platform() {
  local os arch
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  arch="$(uname -m)"

  case "$arch" in
    x86_64|amd64) arch="amd64" ;;
    arm64|aarch64) arch="arm64" ;;
    *) echo "Unsupported architecture: $arch" >&2; exit 1 ;;
  esac

  case "$os" in
    darwin|linux) ;;
    mingw*|msys*|cygwin*) os="windows" ;;
    *) echo "Unsupported OS: $os" >&2; exit 1 ;;
  esac

  echo "${os}-${arch}"
}

# Get the latest release tag from GitHub
get_latest_version() {
  curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' \
    | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/'
}

# Download the binary for the current platform
download_binary() {
  local version="$1"
  local platform
  platform="$(get_platform)"

  local suffix=""
  if [[ "$platform" == windows-* ]]; then
    suffix=".exe"
  fi

  local asset="claude-warp-${platform}${suffix}"
  local url="https://github.com/${REPO}/releases/download/${version}/${asset}"

  echo "Downloading claude-warp ${version} (${platform})..." >&2
  curl -fsSL -o "$BINARY" "$url"
  chmod +x "$BINARY"
  echo "$version" > "$VERSION_FILE"
}

# Download if binary is missing
if [ ! -x "$BINARY" ]; then
  version="$(get_latest_version)"
  download_binary "$version"
fi

exec "$BINARY" "$@"
