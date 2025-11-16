package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette - inspired by One Dark Pro
var (
	// Base colors
	bgColor        = lipgloss.Color("#1e1e1e")
	fgColor        = lipgloss.Color("#abb2bf")
	mutedColor     = lipgloss.Color("#5c6370")
	borderColor    = lipgloss.Color("#3e4451")

	// Accent colors
	primaryColor   = lipgloss.Color("#61afef")  // Blue - for focus
	successColor   = lipgloss.Color("#98c379")  // Green - 2xx
	warningColor   = lipgloss.Color("#e5c07b")  // Yellow - 3xx
	errorColor     = lipgloss.Color("#e06c75")  // Red - 4xx/5xx
	orangeColor    = lipgloss.Color("#d19a66")  // Orange - emphasis
	purpleColor    = lipgloss.Color("#c678dd")  // Purple - keywords
	cyanColor      = lipgloss.Color("#56b6c2")  // Cyan - strings

	// Status colors
	status2xx = successColor
	status3xx = warningColor
	status4xx = orangeColor
	status5xx = errorColor
)

// Text styles
var (
	// Title style
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		Background(bgColor)

	// Help text style
	helpStyle = lipgloss.NewStyle().
		Foreground(mutedColor).
		Background(bgColor)

	// Dimmed text
	dimStyle = lipgloss.NewStyle().
		Foreground(mutedColor)

	// Highlighted text
	highlightStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true)
)

// Border styles
var (
	// Section border style - unfocused
	sectionStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	// Section border style - focused (bright blue)
	focusedSectionStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1)

	// Status bar style
	statusBarStyle = lipgloss.NewStyle().
		Foreground(fgColor).
		Background(borderColor).
		Padding(0, 1).
		Bold(true)
)

// Code/syntax styles
var (
	// Curl command syntax highlighting
	curlCmdStyle    = lipgloss.NewStyle().Foreground(purpleColor).Bold(true) // curl
	curlMethodStyle = lipgloss.NewStyle().Foreground(orangeColor).Bold(true) // GET, POST, etc
	curlFlagStyle   = lipgloss.NewStyle().Foreground(cyanColor)              // -X, -H, -d
	curlStringStyle = lipgloss.NewStyle().Foreground(successColor)           // "strings"
	curlURLStyle    = lipgloss.NewStyle().Foreground(primaryColor).Underline(true)

	// JSON syntax highlighting
	jsonKeyStyle    = lipgloss.NewStyle().Foreground(primaryColor)
	jsonStringStyle = lipgloss.NewStyle().Foreground(successColor)
	jsonNumberStyle = lipgloss.NewStyle().Foreground(orangeColor)
	jsonBoolStyle   = lipgloss.NewStyle().Foreground(purpleColor)
	jsonNullStyle   = lipgloss.NewStyle().Foreground(purpleColor)

	// Code block style
	codeBlockStyle = lipgloss.NewStyle().
		Foreground(fgColor).
		Padding(0)
)

// Status styles
var (
	// HTTP status code styles
	successStatusStyle = lipgloss.NewStyle().
		Foreground(status2xx).
		Bold(true)

	redirectStatusStyle = lipgloss.NewStyle().
		Foreground(status3xx).
		Bold(true)

	clientErrorStatusStyle = lipgloss.NewStyle().
		Foreground(status4xx).
		Bold(true)

	serverErrorStatusStyle = lipgloss.NewStyle().
		Foreground(status5xx).
		Bold(true)

	// Error message style
	errorMessageStyle = lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true)

	// Loading style
	loadingStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true)
)

// UI element styles
var (
	// Method selector
	methodActiveStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true)

	methodInactiveStyle = lipgloss.NewStyle().
		Foreground(mutedColor)

	// Cursor
	cursorStyle = lipgloss.NewStyle().
		Foreground(primaryColor)

	// Placeholder text
	placeholderStyle = lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true)

	// Tab styles for response viewer
	tabActiveStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Underline(true)

	tabInactiveStyle = lipgloss.NewStyle().
		Foreground(mutedColor)
)
