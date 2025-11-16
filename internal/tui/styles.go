package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// TODO: Implement consistent color scheme
//   - Primary color for focused elements
//   - Secondary color for active/selected items
//   - Muted colors for inactive elements
//   - Success/error colors for response status
//   - Syntax highlighting colors for code
//
// TODO: Define border styles
//   - Normal border for unfocused sections
//   - Highlighted border for focused sections
//   - Different border styles for different section types
//
// TODO: Define text styles
//   - Title/header styles
//   - Input field styles
//   - Code/monospace styles for curl preview and response
//   - Help text styles
//
// TODO: Define layout styles
//   - Padding and margins for sections
//   - Width calculations for split panes
//   - Height calculations for scrollable areas

var (
	// Basic color palette - TODO: expand and refine
	primaryColor   = lipgloss.Color("12")  // Blue
	secondaryColor = lipgloss.Color("10")  // Green
	mutedColor     = lipgloss.Color("8")   // Gray
	errorColor     = lipgloss.Color("9")   // Red
	successColor   = lipgloss.Color("10")  // Green

	// Title style
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	// Section border style - unfocused
	sectionStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(mutedColor).
			Padding(1).
			MarginBottom(1)

	// Section border style - focused
	focusedSectionStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Padding(1).
				MarginBottom(1)

	// Input field style
	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15"))

	// Code/monospace style for curl and response
	codeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Background(lipgloss.Color("0")).
			Padding(1)

	// Error message style
	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	// Success message style
	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)
)
