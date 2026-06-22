# EasyFlow 🚀

EasyFlow is a **terminal-first GitHub workflow automation tool** built in Go. It replaces repetitive Git + GitHub UI actions with a **guided interactive terminal application** powered by Bubble Tea.

![EasyFlow Demo](docs/demo-project.png)

## Quick Start

```bash
easyflow
```

Complete your entire development workflow inside one interactive UI.

## � Documentation

For comprehensive documentation, architecture details, and usage guides, visit the [docs folder](docs/).

### Quick Links

- [Architecture Overview](docs/architecture.md) - System design and component relationships
- [Module Documentation](docs/modules.md) - Detailed documentation for each package
- [Workflow Guide](docs/workflow.md) - Complete workflow automation guide
- [API Reference](docs/api.md) - Function and type references

## Core Workflow

```
Issue → Branch → Code → Commit → Push → PR → Merge → Close Issue
```

Everything is interactive, fast, and keyboard-driven.

## Tech Stack

- **Go** - Core language
- **Bubble Tea** - UI engine
- **GitHub CLI (`gh`)** - GitHub API operations
- **Git** - Local repository operations

## 🎯 Goal

Replace GitHub Web UI, manual git commands, and context switching with **a single terminal application** that runs your entire dev workflow.