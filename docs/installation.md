# Installation & Setup

This guide covers installing and configuring EasyFlow on your system.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Verification](#verification)
- [Uninstallation](#uninstallation)

---

## Prerequisites

Before installing EasyFlow, ensure you have the following prerequisites installed:

### Required Software

| Software | Version | Purpose |
|----------|---------|---------|
| **Go** | 1.21 or higher | Build and run EasyFlow |
| **Git** | 2.0 or higher | Version control operations |
| **GitHub CLI (gh)** | 2.0 or higher | GitHub API operations |

### Checking Prerequisites

#### Check Go Installation

```bash
go version
```

Expected output: `go version go1.21.x darwin/amd64` (or similar)

#### Check Git Installation

```bash
git --version
```

Expected output: `git version 2.x.x`

#### Check GitHub CLI Installation

```bash
gh --version
```

Expected output: `gh version 2.x.x`

#### Check GitHub CLI Authentication

```bash
gh auth status
```

Expected output: Authentication status showing you're logged in

---

## Installation

### Method 1: Build from Source

#### Clone the Repository

```bash
git clone https://github.com/Erebus9456/easyflow.git
cd easyflow
```

#### Build the Binary

```bash
go build -o easyflow
```

This creates the `easyflow` binary in the current directory.

#### Install to System Path (Optional)

```bash
# macOS/Linux
sudo mv easyflow /usr/local/bin/

# Or add to your PATH
export PATH=$PATH:$(pwd)
```

#### Verify Installation

```bash
easyflow --help
```

### Method 2: Using Go Install

```bash
go install github.com/Erebus9456/easyflow@latest
```

This installs EasyFlow to `~/go/bin/`. Ensure this directory is in your PATH:

```bash
export PATH=$PATH:~/go/bin
```

### Method 3: Using Installation Script

```bash
curl -fsSL https://raw.githubusercontent.com/Erebus9456/easyflow/main/scripts/install.sh | bash
```

---

## Configuration

### GitHub CLI Authentication

If you haven't authenticated with GitHub CLI, run:

```bash
gh auth login
```

Follow the prompts:
1. Select "GitHub.com"
2. Select "HTTPS" or "SSH" for protocol
3. Choose "Login with a web browser"
4. Copy the device code
5. Paste in browser when prompted
6. Authorize EasyFlow to access your repositories

### Git Configuration

Ensure your Git identity is configured:

```bash
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

### Repository Setup

EasyFlow must be run from within a Git repository with a GitHub remote:

```bash
# Navigate to your repository
cd /path/to/your/repository

# Verify it's a Git repository
git status

# Verify remote is configured
git remote -v
```

Expected output should show an `origin` remote pointing to GitHub.

---

## Verification

### Test EasyFlow Installation

Run EasyFlow to verify it's working:

```bash
easyflow
```

You should see the EasyFlow dashboard interface.

### Test GitHub Integration

From the EasyFlow dashboard:
1. Select "🐛 Manage Issues Menu"
2. Select "List Repository Issues"
3. Verify your repository's issues are displayed

### Test Git Integration

From the EasyFlow dashboard:
1. Select "🌿 Manage Branches Menu"
2. Select "Select / Checkout Existing Branch"
3. Verify your repository's branches are displayed

---

## Platform-Specific Setup

### macOS

#### Install Prerequisites with Homebrew

```bash
# Install Go
brew install go

# Install Git
brew install git

# Install GitHub CLI
brew install gh
```

#### Add Go to PATH

Add to your `~/.zshrc` or `~/.bash_profile`:

```bash
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$HOME/go/bin
```

Reload your shell:

```bash
source ~/.zshrc
```

### Linux

#### Install Prerequisites

**Ubuntu/Debian:**

```bash
sudo apt update
sudo apt install golang git
```

**Install GitHub CLI:**

```bash
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
sudo apt update
sudo apt install gh
```

**Fedora/RHEL:**

```bash
sudo dnf install golang git
sudo dnf install gh
```

### Windows

#### Install Prerequisites

1. **Go**: Download from [golang.org/dl](https://golang.org/dl/)
2. **Git**: Download from [git-scm.com](https://git-scm.com/)
3. **GitHub CLI**: Download from [cli.github.com](https://cli.github.com/)

#### Build EasyFlow

Open PowerShell or Command Prompt:

```powershell
git clone https://github.com/Erebus9456/easyflow.git
cd easyflow
go build -o easyflow.exe
```

#### Add to PATH

Add the directory containing `easyflow.exe` to your system PATH.

---

## Environment Variables

EasyFlow supports the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `EASYFLOW_DEBUG` | Enable debug output | `0` (disabled) |
| `EASYFLOW_CONFIG` | Path to config file | `~/.easyflow/config.yaml` |

### Enable Debug Mode

```bash
export EASYFLOW_DEBUG=1
easyflow
```

---

## Common Installation Issues

### Issue: "go: command not found"

**Solution**: Install Go and add to PATH (see Platform-Specific Setup above)

### Issue: "gh: command not found"

**Solution**: Install GitHub CLI (see Platform-Specific Setup above)

### Issue: "git: command not found"

**Solution**: Install Git (see Platform-Specific Setup above)

### Issue: "not a git repository"

**Solution**: Navigate to a Git repository before running EasyFlow

### Issue: "git remote 'origin' is missing"

**Solution**: Add a GitHub remote to your repository:

```bash
git remote add origin https://github.com/username/repo.git
```

### Issue: "gh cli is not authenticated"

**Solution**: Authenticate with GitHub CLI:

```bash
gh auth login
```

---

## Uninstallation

### Remove Binary

```bash
# If installed to /usr/local/bin
sudo rm /usr/local/bin/easyflow

# If installed via go install
rm ~/go/bin/easyflow

# If in current directory
rm easyflow
```

### Remove Configuration

```bash
# Remove config directory
rm -rf ~/.easyflow
```

### Remove from PATH

Edit your shell configuration file (`~/.zshrc`, `~/.bash_profile`, etc.) and remove the EasyFlow PATH entries.

---

## Upgrading

### Upgrade from Source

```bash
cd easyflow
git pull origin main
go build -o easyflow
```

### Upgrade via Go Install

```bash
go install github.com/Erebus9456/easyflow@latest
```

---

## Next Steps

After installation, check out these guides:

- [Quick Start Guide](quickstart.md) - Get started with EasyFlow
- [Workflow Guide](workflow.md) - Learn about workflow automation
- [Configuration](configuration.md) - Customize EasyFlow settings

---

**Related Documentation**:
- [Quick Start Guide](quickstart.md)
- [Configuration](configuration.md)
- [Troubleshooting](troubleshooting.md)
