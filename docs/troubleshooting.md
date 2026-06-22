# Troubleshooting

This guide helps you diagnose and resolve common issues with EasyFlow.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Runtime Issues](#runtime-issues)
- [Git Issues](#git-issues)
- [GitHub Issues](#github-issues)
- [UI Issues](#ui-issues)
- [Performance Issues](#performance-issues)
- [Getting Help](#getting-help)

---

## Installation Issues

### Issue: "go: command not found"

**Symptoms**: Command not found when trying to build or run EasyFlow

**Causes**: Go is not installed or not in PATH

**Solutions**:

1. **Install Go**:
   ```bash
   # macOS
   brew install go
   
   # Linux (Ubuntu/Debian)
   sudo apt install golang
   
   # Linux (Fedora/RHEL)
   sudo dnf install golang
   ```

2. **Add Go to PATH**:
   ```bash
   # Add to ~/.zshrc or ~/.bash_profile
   export PATH=$PATH:/usr/local/go/bin
   export PATH=$PATH:$HOME/go/bin
   
   # Reload shell
   source ~/.zshrc
   ```

3. **Verify Installation**:
   ```bash
   go version
   ```

### Issue: "gh: command not found"

**Symptoms**: GitHub CLI command not found

**Causes**: GitHub CLI is not installed or not in PATH

**Solutions**:

1. **Install GitHub CLI**:
   ```bash
   # macOS
   brew install gh
   
   # Linux (Ubuntu/Debian)
   curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
   echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
   sudo apt update
   sudo apt install gh
   ```

2. **Verify Installation**:
   ```bash
   gh --version
   ```

### Issue: "git: command not found"

**Symptoms**: Git command not found

**Causes**: Git is not installed or not in PATH

**Solutions**:

1. **Install Git**:
   ```bash
   # macOS
   brew install git
   
   # Linux (Ubuntu/Debian)
   sudo apt install git
   
   # Linux (Fedora/RHEL)
   sudo dnf install git
   ```

2. **Verify Installation**:
   ```bash
   git --version
   ```

### Issue: Build fails with dependency errors

**Symptoms**: `go build` fails with dependency errors

**Causes**: Missing or outdated dependencies

**Solutions**:

1. **Download dependencies**:
   ```bash
   go mod download
   ```

2. **Tidy dependencies**:
   ```bash
   go mod tidy
   ```

3. **Verify dependencies**:
   ```bash
   go mod verify
   ```

4. **Clean and rebuild**:
   ```bash
   go clean -cache
   go build -o easyflow
   ```

---

## Runtime Issues

### Issue: "not a git repository"

**Symptoms**: Error when starting EasyFlow

**Causes**: Not running from within a Git repository

**Solutions**:

1. **Navigate to Git repository**:
   ```bash
   cd /path/to/your/git/repository
   easyflow
   ```

2. **Initialize Git repository** (if needed):
   ```bash
   git init
   git add .
   git commit -m "Initial commit"
   ```

3. **Verify repository status**:
   ```bash
   git status
   ```

### Issue: "git remote 'origin' is missing"

**Symptoms**: Error when trying to push or create PR

**Causes**: No remote configured for repository

**Solutions**:

1. **Add remote**:
   ```bash
   git remote add origin https://github.com/username/repo.git
   ```

2. **Verify remote**:
   ```bash
   git remote -v
   ```

3. **Update remote URL** (if incorrect):
   ```bash
   git remote set-url origin https://github.com/username/repo.git
   ```

### Issue: "gh cli is not authenticated"

**Symptoms**: Error when trying to fetch issues or create PR

**Causes**: GitHub CLI not authenticated

**Solutions**:

1. **Authenticate with GitHub CLI**:
   ```bash
   gh auth login
   ```

2. **Follow the prompts**:
   - Select "GitHub.com"
   - Select "HTTPS" or "SSH"
   - Choose "Login with a web browser"
   - Copy and paste the device code

3. **Verify authentication**:
   ```bash
   gh auth status
   ```

### Issue: Stuck in loading state

**Symptoms**: Spinner keeps spinning, no progress

**Causes**: Network issue, API rate limit, or command hanging

**Solutions**:

1. **Press ESC to reset**:
   - Returns to dashboard
   - Clears loading state

2. **Check network connection**:
   ```bash
   ping github.com
   ```

3. **Check GitHub CLI status**:
   ```bash
   gh auth status
   ```

4. **Enable debug mode**:
   ```bash
   export EASYFLOW_DEBUG=1
   easyflow
   ```

5. **Check for API rate limits**:
   ```bash
   gh api /rate_limit
   ```

---

## Git Issues

### Issue: Branch creation fails

**Symptoms**: Error when creating branch

**Causes**: Invalid branch name, branch already exists, or permission issue

**Solutions**:

1. **Check branch name validity**:
   - Avoid special characters: `*`, `?`, `~`, `^`, `:`, `\`
   - Use lowercase and hyphens

2. **Check if branch already exists**:
   ```bash
   git branch
   ```

3. **Delete existing branch** (if needed):
   ```bash
   git branch -D branch-name
   ```

4. **Check Git permissions**:
   ```bash
   git config --list | grep user
   ```

### Issue: Commit fails

**Symptoms**: Error when committing changes

**Causes**: No changes to commit, merge conflict, or pre-commit hook failure

**Solutions**:

1. **Check for changes**:
   ```bash
   git status
   ```

2. **Stage changes manually**:
   ```bash
   git add .
   ```

3. **Resolve merge conflicts**:
   ```bash
   git status
   # Edit conflicted files
   git add .
   git commit
   ```

4. **Skip pre-commit hooks** (if needed):
   ```bash
   git commit --no-verify -m "message"
   ```

### Issue: Push fails

**Symptoms**: Error when pushing to remote

**Causes**: Network issue, authentication failure, or merge conflict

**Solutions**:

1. **Check network connection**:
   ```bash
   ping github.com
   ```

2. **Check remote URL**:
   ```bash
   git remote -v
   ```

3. **Authenticate with Git**:
   ```bash
   git config --global credential.helper store
   git push
   # Enter credentials
   ```

4. **Pull before push**:
   ```bash
   git pull origin main
   git push origin branch-name
   ```

5. **Force push** (use with caution):
   ```bash
   git push --force-with-lease
   ```

---

## GitHub Issues

### Issue: Cannot fetch issues

**Symptoms**: Error when trying to list issues

**Causes**: Authentication failure, repository not found, or API rate limit

**Solutions**:

1. **Verify authentication**:
   ```bash
   gh auth status
   ```

2. **Re-authenticate**:
   ```bash
   gh auth login
   ```

3. **Check repository access**:
   ```bash
   gh repo view
   ```

4. **Check API rate limit**:
   ```bash
   gh api /rate_limit
   ```

5. **Test GitHub CLI directly**:
   ```bash
   gh issue list --limit 5
   ```

### Issue: Cannot create issue

**Symptoms**: Error when creating new issue

**Causes**: Invalid issue title, permission issue, or API error

**Solutions**:

1. **Check issue title**:
   - Ensure title is not empty
   - Avoid special characters

2. **Check repository permissions**:
   ```bash
   gh repo view
   ```

3. **Test GitHub CLI directly**:
   ```bash
   gh issue create --title "Test issue" --body "Test"
   ```

4. **Check for protected branches**:
   ```bash
   gh api repos/owner/repo/branches/main
   ```

### Issue: Cannot create PR

**Symptoms**: Error when creating pull request

**Causes**: Branch not pushed, permission issue, or branch protection rules

**Solutions**:

1. **Ensure branch is pushed**:
   ```bash
   git push -u origin branch-name
   ```

2. **Check branch exists on remote**:
   ```bash
   gh api repos/owner/repo/branches/branch-name
   ```

3. **Check branch protection rules**:
   ```bash
   gh api repos/owner/repo/branches/main/protection
   ```

4. **Test GitHub CLI directly**:
   ```bash
   gh pr create --title "Test PR" --body "Test"
   ```

### Issue: Cannot merge PR

**Symptoms**: Error when merging pull request

**Causes**: Merge conflict, permission issue, or branch protection rules

**Solutions**:

1. **Check PR status**:
   ```bash
   gh pr view
   ```

2. **Resolve merge conflicts**:
   ```bash
   git checkout main
   git pull origin main
   git checkout branch-name
   git rebase main
   # Resolve conflicts
   git push
   ```

3. **Check permissions**:
   ```bash
   gh repo view
   ```

4. **Bypass branch protection** (if allowed):
   ```bash
   gh pr merge --merge --bypass
   ```

---

## UI Issues

### Issue: UI not rendering correctly

**Symptoms**: Garbled text, incorrect colors, or layout issues

**Causes**: Terminal size, color support, or font issues

**Solutions**:

1. **Check terminal size**:
   - Minimum recommended: 80x24
   - Resize terminal window

2. **Check color support**:
   ```bash
   echo $TERM
   ```

3. **Test terminal compatibility**:
   ```bash
   # Test colors
   echo $'\e[31mRed\e[0m'
   ```

4. **Use compatible terminal**:
   - iTerm2 (macOS)
   - Terminal.app (macOS)
   - GNOME Terminal (Linux)
   - Windows Terminal (Windows)

5. **Adjust layout configuration**:
   ```go
   func DefaultLayout() LayoutConfig {
       return LayoutConfig{
           MenuSpacing: 2,
           ColumnWidth: 40,
       }
   }
   ```

### Issue: Keyboard shortcuts not working

**Symptoms**: Keys not responding or unexpected behavior

**Causes**: Terminal configuration or key binding conflicts

**Solutions**:

1. **Test keyboard input**:
   ```bash
   # Test key detection
   cat
   # Press keys, press Ctrl+D to exit
   ```

2. **Check terminal settings**:
   - Ensure no custom key bindings
   - Check for keyboard shortcuts in terminal app

3. **Use alternative keys**:
   - Use `j/k` instead of `↑/↓`
   - Use `q` instead of `Ctrl+C`

4. **Reset terminal**:
   ```bash
   reset
   ```

### Issue: Text input not accepting input

**Symptoms**: Cannot type in text input fields

**Causes**: Terminal focus issue or input mode conflict

**Solutions**:

1. **Click terminal window** to ensure focus
2. **Press ESC** to reset state
3. **Restart EasyFlow**
4. **Check terminal input mode**:
   ```bash
   stty -a
   ```

---

## Performance Issues

### Issue: Slow startup

**Symptoms**: EasyFlow takes long time to start

**Causes**: Large repository, slow network, or system resources

**Solutions**:

1. **Check repository size**:
   ```bash
   du -sh .git
   ```

2. **Check network speed**:
   ```bash
   ping github.com
   ```

3. **Check system resources**:
   ```bash
   top
   ```

4. **Use debug mode to identify bottleneck**:
   ```bash
   export EASYFLOW_DEBUG=1
   easyflow
   ```

### Issue: UI lag or stuttering

**Symptoms**: UI responds slowly or stutters

**Causes**: Large output, slow rendering, or system resources

**Solutions**:

1. **Reduce output size**:
   - Limit commit log display
   - Limit issue list size

2. **Check system resources**:
   ```bash
   top
   ```

3. **Close other applications**
4. **Adjust layout configuration**:
   ```go
   func DefaultLayout() LayoutConfig {
       return LayoutConfig{
           MenuSpacing: 1,
           ColumnWidth: 40,
       }
   }
   ```

### Issue: High memory usage

**Symptoms**: EasyFlow using excessive memory

**Causes**: Memory leak or large data structures

**Solutions**:

1. **Monitor memory usage**:
   ```bash
   ps aux | grep easyflow
   ```

2. **Restart EasyFlow periodically**
3. **Report issue if memory leak suspected**

---

## Debugging

### Enable Debug Mode

Enable debug output for detailed information:

```bash
export EASYFLOW_DEBUG=1
easyflow
```

Debug output includes:
- Detailed error messages
- Stack traces
- Command execution details
- State transitions

### Check Logs

If debug mode is enabled, check terminal output for:
- Error messages
- Command execution results
- State transition logs

### Test Components Individually

Test Git and GitHub CLI separately:

```bash
# Test Git
git status
git branch
git log -5

# Test GitHub CLI
gh auth status
gh issue list --limit 5
gh pr list --limit 5
```

### Verify Environment

Check environment variables:

```bash
echo $EASYFLOW_DEBUG
echo $EASYFLOW_CONFIG
```

---

## Getting Help

### Self-Diagnosis Checklist

Before seeking help, check:

- [ ] Go installed and in PATH
- [ ] Git installed and in PATH
- [ ] GitHub CLI installed and authenticated
- [ ] Running from Git repository
- [ ] Remote configured
- [ ] Network connection working
- [ ] Terminal size adequate (80x24 minimum)
- [ ] Debug mode enabled for detailed output

### Documentation Resources

- [Installation Guide](installation.md) - Setup and installation
- [Quick Start Guide](quickstart.md) - Getting started
- [Workflow Guide](workflow.md) - Workflow automation
- [Configuration](configuration.md) - Configuration options
- [API Reference](api.md) - API documentation

### Community Resources

- [GitHub Issues](https://github.com/Erebus9456/easyflow/issues) - Report bugs
- [GitHub Discussions](https://github.com/Erebus9456/easyflow/discussions) - Ask questions
- [GitHub Wiki](https://github.com/Erebus9456/easyflow/wiki) - Community documentation

### Reporting Issues

When reporting issues, include:

1. **EasyFlow version**: `easyflow --version`
2. **Go version**: `go version`
3. **Git version**: `git --version`
4. **GitHub CLI version**: `gh --version`
5. **Operating system**: `uname -a`
6. **Terminal**: Terminal name and version
7. **Error message**: Full error output
8. **Steps to reproduce**: Detailed reproduction steps
9. **Debug output**: Output with `EASYFLOW_DEBUG=1`

### Issue Template

```markdown
## Description
Brief description of the issue

## Steps to Reproduce
1. Step one
2. Step two
3. Step three

## Expected Behavior
What should happen

## Actual Behavior
What actually happens

## Environment
- EasyFlow version: 
- Go version: 
- Git version: 
- GitHub CLI version: 
- OS: 
- Terminal: 

## Debug Output
```
Paste debug output here
```
```

---

## Common Error Messages

### "git binary not found on your system PATH"

**Solution**: Install Git and add to PATH (see Installation Issues)

### "github cli (gh) binary not found on your system PATH"

**Solution**: Install GitHub CLI (see Installation Issues)

### "not a git repository (or any of the parent directories)"

**Solution**: Navigate to a Git repository or initialize one

### "git remote 'origin' is missing or not configured"

**Solution**: Add remote: `git remote add origin <url>`

### "gh cli is not authenticated. run 'gh auth login' first"

**Solution**: Authenticate: `gh auth login`

### "failed to fetch issues via gh cli"

**Solution**: Check authentication and repository access

### "failed to create branch"

**Solution**: Check branch name validity and permissions

### "failed to push changes to remote"

**Solution**: Check network, authentication, and remote configuration

### "failed to create pull request"

**Solution**: Ensure branch is pushed and check permissions

### "failed to merge pull request"

**Solution**: Resolve conflicts and check permissions

---

**Related Documentation**:
- [Installation Guide](installation.md)
- [Quick Start Guide](quickstart.md)
- [Configuration](configuration.md)
- [Development Guide](development.md)
