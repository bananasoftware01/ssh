#!/usr/bin/env sh
set -e

# Banana Software CLI runner & installer
# Usage:
#   Quick Run:  curl -fsSL https://raw.githubusercontent.com/bananasoftware01/ssh/main/run.sh | bash
#   Install:    curl -fsSL https://raw.githubusercontent.com/bananasoftware01/ssh/main/run.sh | bash -s -- --install

REPO="${BANANA_REPO:-bananasoftware01/ssh}"
TAG="${BANANA_TAG:-latest}"

# 1. Detect Operating System
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
    linux*)
        OS="linux"
        EXT=""
        ;;
    darwin*)
        OS="darwin"
        EXT=""
        ;;
    msys*|mingw*|cygwin*|windows*)
        OS="windows"
        EXT=".exe"
        ;;
    *)
        echo "❌ Unsupported operating system: $OS" >&2
        exit 1
        ;;
esac

# 2. Detect Architecture
ARCH="$(uname -m | tr '[:upper:]' '[:lower:]')"
case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    arm64|aarch64|armv8*)
        ARCH="arm64"
        ;;
    *)
        echo "❌ Unsupported architecture: $ARCH" >&2
        exit 1
        ;;
esac

BINARY_NAME="banana-${OS}-${ARCH}${EXT}"

if [ "$TAG" = "latest" ]; then
    DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}"
else
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${TAG}/${BINARY_NAME}"
fi

# 3. Check for curl / wget
fetch() {
    url="$1"
    dest="$2"
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$url" -o "$dest"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "$dest" "$url"
    else
        echo "❌ Neither curl nor wget was found. Please install one to proceed." >&2
        exit 1
    fi
}

# 4. Handle Installation Mode
INSTALL_MODE=0
for arg in "$@"; do
    if [ "$arg" = "--install" ] || [ "$arg" = "install" ] || [ "$arg" = "-i" ]; then
        INSTALL_MODE=1
    fi
done

if [ "$INSTALL_MODE" -eq 1 ]; then
    INSTALL_DIR="${BANANA_INSTALL_DIR:-$HOME/.local/bin}"
    if [ "$(id -u)" -eq 0 ] && [ -d "/usr/local/bin" ] && [ -z "$BANANA_INSTALL_DIR" ]; then
        INSTALL_DIR="/usr/local/bin"
    fi

    mkdir -p "$INSTALL_DIR"
    TARGET="${INSTALL_DIR}/banana${EXT}"

    echo "🍌 Downloading Banana Software CLI (${OS}/${ARCH})..."
    fetch "$DOWNLOAD_URL" "$TARGET"
    chmod +x "$TARGET"

    echo "✅ Successfully installed to $TARGET"

    case ":$PATH:" in
        *":$INSTALL_DIR:"*) ;;
        *)
            echo ""
            echo "⚠️  Note: $INSTALL_DIR is not in your \$PATH."
            echo "   Add it to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
            echo "   export PATH=\"\$PATH:$INSTALL_DIR\""
            echo ""
            ;;
    esac

    echo "Run 'banana' anytime to launch the terminal landing page!"
    exit 0
fi

# 5. One-Shot Quick Run
TMP_DIR="$(mktemp -d 2>/dev/null || mktemp -d -t 'banana')"
TMP_BIN="${TMP_DIR}/banana${EXT}"
trap 'rm -rf "$TMP_DIR"' EXIT INT TERM

echo "🍌 Loading Banana Software Terminal..." >&2
fetch "$DOWNLOAD_URL" "$TMP_BIN"
chmod +x "$TMP_BIN"

# If input is piped (e.g. `curl ... | bash`), redirect stdin/stdout from /dev/tty
if [ ! -t 0 ] && [ -c /dev/tty ]; then
    exec "$TMP_BIN" "$@" </dev/tty >/dev/tty
else
    exec "$TMP_BIN" "$@"
fi
