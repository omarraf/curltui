package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleUpdate processes messages and updates the model state
// TODO: Implement keyboard navigation between fields
//   - Tab/Shift+Tab: cycle through focusedField (method, URL, headers, body)
//   - Up/Down arrows: navigate within lists (headers)
//   - Enter: confirm input or execute request
//
// TODO: Implement field editing
//   - Typed characters: add to current field value
//   - Backspace/Delete: remove characters
//   - Support multi-line editing for body field
//   - Support adding/removing header entries
//   - Method selector: cycle through GET, POST, PUT, PATCH, DELETE, etc.
//
// TODO: Implement request execution
//   - Ctrl+Enter or F5: execute the current request
//   - Build curl command from request using curl.BuildCurlCommand
//   - Execute using http.Execute
//   - Update model.response with results
//   - Handle execution errors
//
// TODO: Implement curl import on paste
//   - Detect paste events (large text input)
//   - Parse using curl.ParseCurlCommand
//   - Update model.request with parsed values
//   - Handle parsing errors gracefully
func handleUpdate(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		// TODO: Handle Tab, Shift+Tab for field navigation
		// TODO: Handle Up/Down for list navigation
		// TODO: Handle Enter for confirmation/execution
		// TODO: Handle Ctrl+Enter or F5 for request execution
		// TODO: Handle character input for editing
		// TODO: Handle Backspace/Delete for deletion
		}
	}

	return m, nil
}
