#!/bin/bash
set -e

echo "🚀 Starting EasyFlow Installer (macOS / Linux)..."
echo "----------------------------------------"

# Detect Operating System Platform
OS_TYPE="$(uname)"
if [ "$OS_TYPE" = "Darwin" ]; then
    echo "🍏 Platform detected: macOS"
    if ! command -v brew &> /dev/null; then
        echo "❌ Homebrew is required for automated macOS installation. Please install it from https://brew.sh/"
        exit 1
    fi
elif [ "$OS_TYPE" = "Linux" ]; then
    echo "🐧 Platform detected: Linux"
else
    echo "❌ Unsupported OS environment via this shell engine."
    exit 1
fi

# 1. Validate / Install Git & GitHub CLI (gh)
install_dep() {
    local cmd=$1
    if ! command -v "$cmd" &> /dev/null; then
        echo "📦 Installing dependency: $cmd..."
        if [ "$OS_TYPE" = "Darwin" ]; then
            brew install "$cmd"
        else
            sudo apt update && sudo apt install -y "$cmd"
        fi
    else
        echo "✅ $cmd is already installed."
    fi
}

install_dep "git"
install_dep "gh"

# 2. Check and evaluate the Go Toolchain environment
MIN_VERSION="1.21"
UPDATE_NEEDED=false

if ! command -v go &> /dev/null; then
    echo "⚠️ Go runtime is not installed."
    UPDATE_NEEDED=true
else
    CURRENT_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    echo "🔍 Found installed Go version: $CURRENT_VERSION"
    if [ "$(printf '%s\n' "$MIN_VERSION" "$CURRENT_VERSION" | sort -V | head -n1)" != "$MIN_VERSION" ]; then
        echo "⚠️ Go version is older than required v$MIN_VERSION."
        UPDATE_NEEDED=true
    fi
fi

# 3. Handle Go Installation / Upgrade Matrix
if [ "$UPDATE_NEEDED" = true ]; then
    if [ "$OS_TYPE" = "Darwin" ]; then
        echo "📥 Installing/Upgrading Go via Homebrew..."
        brew install go
    else
        echo "📥 Downloading and upgrading toolchain to Go 1.26.0 (Linux)..."
        wget -q --show-progress https://go.dev/dl/go1.26.0.linux-amd64.tar.gz
        sudo rm -rf /usr/local/go
        sudo tar -C /usr/local -xzf go1.26.0.linux-amd64.tar.gz
        rm go1.26.0.linux-amd64.tar.gz
        export PATH=/usr/local/go/bin:$PATH
        
        # Persist path configurations
        for rc in "$HOME/.bashrc" "$HOME/.zshrc"; do
            [ -f "$rc" ] && ! grep -q "/usr/local/go/bin" "$rc" && echo 'export PATH=$PATH:/usr/local/go/bin' >> "$rc"
        done
    fi
else
    echo "✅ Existing Go toolchain satisfies requirements."
fi

# 4. Clear proxy collision cache footprints & Install Tool
echo "🧼 Cleaning module cache footprint registries..."
go clean -modcache

echo "🏎️ Compiling and installing EasyFlow v1.0.2..."
GOPROXY=direct go install github.com/Erebus9456/easyflow@v1.0.2

# Ensure Go Bin is linked to target shell paths
for rc in "$HOME/.bashrc" "$HOME/.zshrc"; do
    [ -f "$rc" ] && ! grep -q '\$HOME/go/bin' "$rc" && echo 'export PATH=$PATH:$HOME/go/bin' >> "$rc"
done

# 5. Verify GitHub CLI (gh) Login Credentials Status
echo "----------------------------------------"
echo "🔐 Verifying GitHub CLI authentication context..."
if gh auth status &> /dev/null; then
    echo "✅ GitHub CLI is authenticated and connected properly!"
else
    echo "⚠️ GitHub CLI is installed but not authenticated."
    echo "🔄 Launching authentication pipeline tool context..."
    gh auth login
fi

echo "----------------------------------------"
echo "🎉 EasyFlow installation completed successfully!"
echo "💡 Please restart your terminal or run 'source ~/.bashrc' (or source ~/.zshrc) to activate commands."
echo "🚀 Run 'easyflow' to begin!"