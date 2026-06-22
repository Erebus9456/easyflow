#!/bin/bash

set -e

REPO_NAME="easyflow"

echo "🚀 Creating empty EasyFlow project structure..."

# 1. Create directory tree
mkdir -p cmd
mkdir -p internal/{ui,workflow,github,git,config}
mkdir -p utils scripts docs

# 2. Initialize Go Module
if [ ! -f go.mod ]; then
    echo "📦 Initializing Go module..."
    go mod init easyflow
fi

# 3. Create completely empty structural files
touch main.go
touch cmd/root.go
touch internal/ui/{model.go,update.go,view.go,menu.go,styles.go}
touch internal/workflow/{workflow.go,state.go}
touch internal/github/{issues.go,pr.go}
touch internal/git/{branch.go,commit.go,push.go}
touch internal/config/config.go
touch utils/{shell.go,errors.go}

echo "📁 Folders and empty files generated successfully."

# 4. Git & GitHub Setup
if [ ! -d .git ]; then
    echo "🗃️ Initializing local Git repository..."
    git init -b main
fi

# Stage and commit the empty skeleton
git add .
git commit -m "Initial commit: Skeleton workspace for EasyFlow"

# Verify GitHub CLI and create the remote repository
if command -v gh &> /dev/null; then
    echo "🐙 Connecting to GitHub..."
    if ! git remote get-url origin &> /dev/null; then
        echo "Creating GitHub repository '$REPO_NAME'..."
        # Creates a public repository and links it as origin
        gh repo create "$REPO_NAME" --public --source=. --remote=origin --push
        echo "✅ Repository successfully linked and pushed to GitHub!"
    else
        echo "ℹ️ Git remote 'origin' already exists."
    fi
else
    echo "⚠️ GitHub CLI ('gh') not found. Skipping automated GitHub remote creation."
    echo "Please link manually using: git remote add origin <your-repo-url>"
fi

echo "🎉 Workspace setup complete!"