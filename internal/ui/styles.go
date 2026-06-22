package ui

import "github.com/charmbracelet/lipgloss"

// UI Color Palette Tokens
var (
	ColorPrimary   = lipgloss.Color("#8633FF") // Vibrant Purple
	ColorSecondary = lipgloss.Color("#00F5D4") // Bright Aqua
	ColorSuccess   = lipgloss.Color("#70E000") // Lime Green
	ColorError     = lipgloss.Color("#FF0054") // Deep Red
	ColorNeutral   = lipgloss.Color("#3A3A3A") // Slate Gray
	ColorTextMuted = lipgloss.Color("#757575") // Muted Label Gray
)

// UI Element Styles Definitions
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
