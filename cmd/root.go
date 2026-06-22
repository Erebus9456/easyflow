package cmd

import (
	"fmt"
	"os"

	"github.com/Erebus9456/easyflow/internal/git"
	"github.com/Erebus9456/easyflow/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "easyflow",
	Short: "EasyFlow is a terminal-first GitHub workflow automation engine.",
	Long:  `An interactive pipeline utility that replaces browser context-switching with an automated terminal dashboard.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Run dynamic repo auto-detection before starting the TUI
		repoCtx, err := git.GetCurrentRepo()
		if err != nil {
			fmt.Printf("❌ Environment Error: %v\n", err)
			os.Exit(1)
		}

		// 2. Initialize the Bubble Tea program with alternate screen enabled
		p := tea.NewProgram(ui.InitialModel(repoCtx), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("❌ Critical UI Error Encountered: %v\n", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
