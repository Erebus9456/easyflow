# Configuration

This document covers configuration options and customization for EasyFlow.

## Table of Contents

- [Overview](#overview)
- [UI Configuration](#ui-configuration)
- [Environment Variables](#environment-variables)
- [Custom Configuration File](#custom-configuration-file)
- [Styling Customization](#styling-customization)
- [Advanced Configuration](#advanced-configuration)

---

## Overview

EasyFlow provides several ways to customize its behavior and appearance:

1. **UI Configuration** - Adjust layout and spacing
2. **Environment Variables** - Control runtime behavior
3. **Configuration File** - Persistent settings (future)
4. **Styling Customization** - Modify colors and themes
5. **Advanced Configuration** - Extend functionality

---

## UI Configuration

### Layout Configuration

The UI layout is controlled by the `LayoutConfig` struct in `internal/config/config.go`.

```go
type LayoutConfig struct {
    MenuSpacing int  // Number of newlines between menu items
    ColumnWidth int  // Width of UI columns
}
```

### Default Layout

```go
func DefaultLayout() LayoutConfig {
    return LayoutConfig{
        MenuSpacing: 4,  // Relaxed spacing between menu items
        ColumnWidth: 50, // Standard column width for split layout
    }
}
```

### Customizing Layout

To customize the layout, modify the `DefaultLayout()` function:

```go
func DefaultLayout() LayoutConfig {
    return LayoutConfig{
        MenuSpacing: 2,  // Tighter spacing
        ColumnWidth: 60, // Wider columns
    }
}
```

### Layout Options

#### Menu Spacing

Controls the vertical spacing between menu items:

| Value | Description | Use Case |
|-------|-------------|----------|
| `1` | Minimal spacing | Compact displays |
| `2` | Relaxed spacing | Standard terminals |
| `4` | Spacious spacing | Large displays |

#### Column Width

Controls the width of UI columns:

| Value | Description | Use Case |
|-------|-------------|----------|
| `40` | Narrow columns | Small terminals |
| `50` | Standard columns | Default setting |
| `60` | Wide columns | Large terminals |
| `80` | Extra wide | Very large displays |

### Example: Compact Layout

```go
func CompactLayout() LayoutConfig {
    return LayoutConfig{
        MenuSpacing: 1,  // Minimal spacing
        ColumnWidth: 40,  // Narrow columns
    }
}
```

### Example: Spacious Layout

```go
func SpaciousLayout() LayoutConfig {
    return LayoutConfig{
        MenuSpacing: 6,  // Extra spacing
        ColumnWidth: 70,  // Wide columns
    }
}
```

---

## Environment Variables

EasyFlow supports the following environment variables:

### EASYFLOW_DEBUG

Enable debug output for troubleshooting.

```bash
export EASYFLOW_DEBUG=1
easyflow
```

**Values**:
- `0` or unset: Debug disabled (default)
- `1`: Debug enabled

**Output**:
- Detailed error messages
- Stack traces
- Command execution details

### EASYFLOW_CONFIG

Specify a custom configuration file path.

```bash
export EASYFLOW_CONFIG=/path/to/config.yaml
easyflow
```

**Default**: `~/.easyflow/config.yaml`

### EASYFLOW_MENU_SPACING

Override menu spacing without recompiling.

```bash
export EASYFLOW_MENU_SPACING=2
easyflow
```

**Default**: `4`

### EASYFLOW_COLUMN_WIDTH

Override column width without recompiling.

```bash
export EASYFLOW_COLUMN_WIDTH=60
easyflow
```

**Default**: `50`

---

## Custom Configuration File

Configuration file support is planned for future releases. The configuration will use YAML format.

### Planned Configuration Structure

```yaml
# ~/.easyflow/config.yaml

# UI Settings
ui:
  menu_spacing: 4
  column_width: 50
  
# Git Settings
git:
  default_branch: main
  auto_stage: true
  
# GitHub Settings
github:
  default_pr_base: main
  auto_delete_branch: true
  
# Workflow Settings
workflow:
  pipeline_mode: false
  require_issue: true
  
# Theme Settings
theme:
  primary_color: "#8633FF"
  secondary_color: "#00F5D4"
  success_color: "#70E000"
  error_color: "#FF0054"
```

### Creating Configuration File

```bash
mkdir -p ~/.easyflow
cat > ~/.easyflow/config.yaml << EOF
ui:
  menu_spacing: 2
  column_width: 60
EOF
```

---

## Styling Customization

### Color Palette

Colors are defined in `internal/ui/styles.go`:

```go
var (
    ColorPrimary   = lipgloss.Color("#8633FF")  // Vibrant Purple
    ColorSecondary = lipgloss.Color("#00F5D4")  // Bright Aqua
    ColorSuccess   = lipgloss.Color("#70E000")  // Lime Green
    ColorError     = lipgloss.Color("#FF0054")  // Deep Red
    ColorNeutral   = lipgloss.Color("#3A3A3A")  // Slate Gray
    ColorTextMuted = lipgloss.Color("#757575")  // Muted Label Gray
)
```

### Customizing Colors

To customize the color palette, modify the color constants:

```go
var (
    ColorPrimary   = lipgloss.Color("#FF6B6B")  // Red
    ColorSecondary = lipgloss.Color("#4ECDC4")  // Teal
    ColorSuccess   = lipgloss.Color("#95E1D3")  // Mint
    ColorError     = lipgloss.Color("#FF4757")  // Bright Red
    ColorNeutral   = lipgloss.Color("#2F3542")  // Dark Gray
    ColorTextMuted = lipgloss.Color("#A4B0BE")  // Light Gray
)
```

### Style Definitions

Styles are defined in `internal/ui/styles.go`:

```go
var (
    StyleTitle = lipgloss.NewStyle().
            Bold(true).
            Background(ColorPrimary).
            Foreground(lipgloss.Color("#FFFFFF")).
            Padding(0, 1).
            MarginBottom(1)

    StyleHeader = lipgloss.NewStyle().
            Bold(true).
            Foreground(ColorSecondary).
            MarginBottom(1)

    StyleSelectedOption = lipgloss.NewStyle().
            Bold(true).
            Foreground(ColorPrimary)

    StyleUnselectedOption = lipgloss.NewStyle().
            Foreground(lipgloss.Color("#FFFFFF"))

    StyleSuccessBanner = lipgloss.NewStyle().
            Bold(true).
            Foreground(ColorSuccess).
            Border(lipgloss.NormalBorder()).
            BorderForeground(ColorSuccess).
            Padding(0, 2).
            MarginTop(1)

    StyleErrorBanner = lipgloss.NewStyle().
            Bold(true).
            Foreground(ColorError).
            Border(lipgloss.NormalBorder()).
            BorderForeground(ColorError).
            Padding(0, 2).
            MarginTop(1)

    StyleHelpText = lipgloss.NewStyle().
            Foreground(ColorTextMuted).
            MarginTop(1).
            Italic(true)
)
```

### Customizing Styles

Modify style definitions to change appearance:

```go
StyleTitle = lipgloss.NewStyle().
        Bold(true).
        Background(ColorPrimary).
        Foreground(lipgloss.Color("#FFFFFF")).
        Padding(1, 2).  // Increased padding
        MarginBottom(2) // Increased margin
```

### Theme Presets

#### Dark Theme (Default)

```go
ColorPrimary   = lipgloss.Color("#8633FF")
ColorSecondary = lipgloss.Color("#00F5D4")
ColorSuccess   = lipgloss.Color("#70E000")
ColorError     = lipgloss.Color("#FF0054")
ColorNeutral   = lipgloss.Color("#3A3A3A")
ColorTextMuted = lipgloss.Color("#757575")
```

#### Light Theme

```go
ColorPrimary   = lipgloss.Color("#6C5CE7")
ColorSecondary = lipgloss.Color("#00CEC9")
ColorSuccess   = lipgloss.Color("#00B894")
ColorError     = lipgloss.Color("#D63031")
ColorNeutral   = lipgloss.Color("#636E72")
ColorTextMuted = lipgloss.Color("#B2BEC3")
```

#### Monochrome Theme

```go
ColorPrimary   = lipgloss.Color("#FFFFFF")
ColorSecondary = lipgloss.Color("#CCCCCC")
ColorSuccess   = lipgloss.Color("#AAAAAA")
ColorError     = lipgloss.Color("#888888")
ColorNeutral   = lipgloss.Color("#666666")
ColorTextMuted = lipgloss.Color("#444444")
```

---

## Advanced Configuration

### Custom Menu Items

Add custom menu items in `internal/ui/menu.go`:

```go
func GetMainMenuOptions() []MainMenuItem {
    return []MainMenuItem{
        // ... existing items ...
        {
            Title:       "🔧 Custom Action",
            Description: "Your custom action description",
        },
    }
}
```

Then handle the selection in `internal/ui/update.go`:

```go
func (m AppModel) handleMenuSelection() (tea.Model, tea.Cmd) {
    // ... existing cases ...
    switch m.Cursor {
    // ... existing cases ...
    case 7: // Custom action index
        // Your custom logic here
        m.Engine.Advance(workflow.StateCustomAction)
    }
    return m, nil
}
```

### Custom Workflow States

Add new states in `internal/workflow/state.go`:

```go
const (
    // ... existing states ...
    StateCustomAction State = iota
)
```

Add validation in `internal/workflow/workflow.go`:

```go
func (e *Engine) Advance(next State) error {
    switch next {
    // ... existing cases ...
    case StateCustomAction:
        // Your validation logic
    }
    e.Ctx.CurrentStep = next
    return nil
}
```

Add UI handling in `internal/ui/update.go`:

```go
case workflow.StateCustomAction:
    if keyMsg, ok := msg.(tea.KeyMsg); ok {
        switch keyMsg.String() {
        case "enter":
            // Your custom action logic
        }
    }
```

Add view rendering in `internal/ui/view.go`:

```go
case workflow.StateCustomAction:
    leftSide.WriteString(StyleHeader.Render("🔧 Custom Action"))
    leftSide.WriteString("\n\n")
    leftSide.WriteString("Your custom action UI here")
```

### Custom Git Operations

Add new Git operations in `internal/git/`:

```go
// custom.go
package git

import (
    "fmt"
    "github.com/Erebus9456/easyflow/utils"
)

func CustomGitOperation() (string, error) {
    stdout, stderr, err := utils.ExecuteCommand("git", "custom", "command")
    if err != nil {
        return "", fmt.Errorf("failed to execute custom operation: %s %w", stderr, err)
    }
    return stdout, nil
}
```

### Custom GitHub Operations

Add new GitHub operations in `internal/github/`:

```go
// custom.go
package github

import (
    "fmt"
    "github.com/Erebus9456/easyflow/utils"
)

func CustomGitHubOperation() (string, error) {
    stdout, stderr, err := utils.ExecuteCommand("gh", "custom", "command")
    if err != nil {
        return "", fmt.Errorf("failed to execute custom operation: %s %w", stderr, err)
    }
    return stdout, nil
}
```

---

## Configuration Examples

### Example 1: Compact Terminal

For small terminal windows (80x24):

```go
func DefaultLayout() LayoutConfig {
    return LayoutConfig{
        MenuSpacing: 1,  // Minimal spacing
        ColumnWidth: 35, // Narrow columns
    }
}
```

### Example 2: Large Display

For large terminal windows (120x40+):

```go
func DefaultLayout() LayoutConfig {
    return LayoutConfig{
        MenuSpacing: 6,  // Extra spacing
        ColumnWidth: 70,  // Wide columns
    }
}
```

### Example 3: High Contrast

For better visibility:

```go
var (
    ColorPrimary   = lipgloss.Color("#FFFF00")  // Yellow
    ColorSecondary = lipgloss.Color("#00FFFF")  // Cyan
    ColorSuccess   = lipgloss.Color("#00FF00")  // Bright Green
    ColorError     = lipgloss.Color("#FF0000")  // Bright Red
    ColorNeutral   = lipgloss.Color("#FFFFFF")  // White
    ColorTextMuted = lipgloss.Color("#AAAAAA")  // Light Gray
)
```

### Example 4: Minimalist

For a clean, minimal look:

```go
var (
    ColorPrimary   = lipgloss.Color("#FFFFFF")
    ColorSecondary = lipgloss.Color("#FFFFFF")
    ColorSuccess   = lipgloss.Color("#FFFFFF")
    ColorError     = lipgloss.Color("#FFFFFF")
    ColorNeutral   = lipgloss.Color("#FFFFFF")
    ColorTextMuted = lipgloss.Color("#888888")
)

StyleTitle = lipgloss.NewStyle().
        Bold(true).
        Foreground(ColorPrimary).
        Padding(0, 0).
        MarginBottom(0)
```

---

## Configuration Best Practices

### 1. Test Configuration Changes

Always test configuration changes in a safe environment:

```bash
# Test with debug mode
export EASYFLOW_DEBUG=1
easyflow
```

### 2. Backup Original Configuration

Before making changes, backup the original:

```bash
cp internal/config/config.go internal/config/config.go.backup
```

### 3. Use Version Control

Commit configuration changes with descriptive messages:

```bash
git add internal/config/config.go
git commit -m "Customize UI layout for better readability"
```

### 4. Document Customizations

Add comments to explain custom configurations:

```go
// Customized for 80x24 terminal windows
func DefaultLayout() LayoutConfig {
    return LayoutConfig{
        MenuSpacing: 1,  // Reduced for compact display
        ColumnWidth: 35, // Narrower for small screens
    }
}
```

### 5. Share Configuration

Create different configuration profiles for different use cases:

```go
// CompactLayout for small terminals
func CompactLayout() LayoutConfig { ... }

// StandardLayout for normal use
func StandardLayout() LayoutConfig { ... }

// SpaciousLayout for large displays
func SpaciousLayout() LayoutConfig { ... }
```

---

## Troubleshooting Configuration

### Issue: Configuration Not Applied

**Solution**: Ensure you've rebuilt the binary after modifying code:

```bash
go build -o easyflow
```

### Issue: Colors Not Displaying Correctly

**Solution**: Check your terminal supports color:

```bash
# Test color support
echo $TERM
```

If colors don't work, your terminal may not support 256-color mode.

### Issue: Layout Issues on Small Terminals

**Solution**: Reduce column width and menu spacing:

```go
func DefaultLayout() LayoutConfig {
    return LayoutConfig{
        MenuSpacing: 1,
        ColumnWidth: 30,
    }
}
```

---

## Future Configuration Features

Planned features for future releases:

- **Configuration File Support** - YAML-based configuration
- **Theme Switching** - Built-in theme presets
- **Profile Management** - Multiple configuration profiles
- **Runtime Configuration** - Change settings without restart
- **Import/Export** - Share configurations with team

---

**Related Documentation**:
- [Installation Guide](installation.md) - Setup and installation
- [Architecture Overview](architecture.md) - System design
- [Module Documentation](modules.md) - Component details
