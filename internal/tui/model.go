package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/username/curltui/internal/curl"
	"github.com/username/curltui/internal/http"
)

// Model represents the application state
type Model struct {
	request      curl.Request
	response     http.Response
	focusedField int
	err          error
	width        int
	height       int
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return handleUpdate(m, msg)
}

// View renders the UI
func (m Model) View() string {
	return renderView(m)
}
