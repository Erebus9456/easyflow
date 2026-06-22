package ui

import (
	"github.com/Erebus9456/easyflow/internal/config"
	"github.com/Erebus9456/easyflow/internal/git"
	"github.com/Erebus9456/easyflow/internal/github"
	"github.com/Erebus9456/easyflow/internal/workflow"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type AppModel struct {
	Engine      *workflow.Engine
	RepoCtx     *git.RepoContext
	MenuItems   []MainMenuItem
	Cursor      int
	Issues      []github.Issue
	IssueCursor int

	// Unified Configuration Structure
	Config config.Config

	TextInput textinput.Model
	Spinner   spinner.Model

	Loading      bool
	ErrorMessage string
	SuccessMsg   string
}

// Init implements [tea.Model].
func (m AppModel) Init() tea.Cmd {
	return tea.Batch(
		m.Spinner.Tick,  // Starts the loading dot animation loop
		textinput.Blink, // Starts the text input field cursor blinking
	)
}

func InitialModel(repo *git.RepoContext) AppModel {
	ti := textinput.New()
	ti.Placeholder = "Enter value..."
	ti.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = StyleSelectedOption

	return AppModel{
		Engine:    workflow.NewEngine(),
		RepoCtx:   repo,
		MenuItems: GetMainMenuOptions(),
		Config:    config.NewDefaultConfig(), // 👈 Fixes compiler error! Uses the new unified config initialization
		TextInput: ti,
		Spinner:   s,
	}
}
