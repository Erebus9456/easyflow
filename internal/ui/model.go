package ui

import (
	"github.com/Erebus9456/easyflow/internal/git"
	"github.com/Erebus9456/easyflow/internal/github"
	"github.com/Erebus9456/easyflow/internal/workflow"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// AppModel holds the state metrics for the runtime interface loop
type AppModel struct {
	Engine      *workflow.Engine
	RepoCtx     *git.RepoContext
	MenuItems   []MainMenuItem
	Cursor      int
	Issues      []github.Issue
	IssueCursor int

	// Input Subcomponents
	TextInput textinput.Model
	Spinner   spinner.Model

	// Runtime tracking properties
	Loading      bool
	ErrorMessage string
	SuccessMsg   string
}

// View implements [tea.Model].

// InitialModel configures a clean default state instantiation
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
		TextInput: ti,
		Spinner:   s,
	}
}

// Init triggers initial loading tick sequences for subcomponents
func (m AppModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.Spinner.Tick)
}
