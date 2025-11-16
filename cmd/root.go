package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/username/curltui/internal/tui"
)

var importCurl string

var rootCmd = &cobra.Command{
	Use:   "curltui",
	Short: "A TUI wrapper around curl for building and executing HTTP requests",
	Long:  `curltui is a terminal UI application that helps you build and execute HTTP requests using curl.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize the TUI model
		model := tui.NewModel(importCurl)

		// Create and run the Bubble Tea program
		p := tea.NewProgram(model, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.Flags().StringVar(&importCurl, "import", "", "Import a curl command string to edit")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
