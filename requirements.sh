#!/bin/bash

echo "🔄 Checking global system requirements..."

# 1. Verify Go Installation
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go (https://go.dev/doc/install) and try again."
    exit 1
else
    echo "✅ Go is installed: $(go version)"
fi

# 2. Verify Git Installation
if ! command -v git &> /dev/null; then
    echo "❌ Git is not installed. Please install Git."
    exit 1
else
    echo "✅ Git is installed"
fi

# 3. Verify GitHub CLI Installation
if ! command -v gh &> /dev/null; then
    echo "⚠️ Warning: GitHub CLI ('gh') is missing. EasyFlow relies on it for API tasks."
    echo "👉 Install instructions: https://cli.github.com"
else
    echo "✅ GitHub CLI ('gh') is installed"
fi

# 4. Pull Go frameworks and libraries for your project
echo "📥 Installing Go project dependencies..."

# Bubble Tea TUI Framework
go get github.com/charmbracelet/bubbletea@latest
# Lip Gloss (UI Styling)
go get github.com/charmbracelet/lipgloss@latest
# Bubbles (TUI Components)
go get github.com/charmbracelet/bubbles@latest
# Cobra CLI framework
go get github.com/spf13/cobra@latest

# Tidy tracking matrices
go mod tidy

echo "🔥 All dependencies fetched. Workspace is ready for code!"