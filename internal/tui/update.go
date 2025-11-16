package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/username/curltui/internal/curl"
	"github.com/username/curltui/internal/http"
)

// ExecuteMsg is sent when a request execution is complete
type ExecuteMsg struct {
	response http.Response
	err      error
}

// executeRequest runs the curl command in the background
func executeRequest(req curl.Request) tea.Cmd {
	return func() tea.Msg {
		curlCmd := curl.BuildCurlCommand(req)
		resp, err := http.Execute(curlCmd)
		return ExecuteMsg{response: resp, err: err}
	}
}

// handleUpdate processes messages and updates the model state
func handleUpdate(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case ExecuteMsg:
		m.loading = false
		m.response = msg.response
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		// Handle paste mode
		if m.pasteMode {
			return handlePasteMode(m, msg)
		}

		// Global shortcuts
		switch msg.String() {
		case "ctrl+c":
			if !m.editing {
				return m, tea.Quit
			}

		case "q":
			if !m.editing {
				return m, tea.Quit
			}

		case "ctrl+p":
			// Enter paste mode
			m.pasteMode = true
			m.pasteBuffer = ""
			return m, nil

		case "ctrl+e", "ctrl+enter":
			// Execute request
			if !m.loading {
				m.loading = true
				m.err = nil
				req := m.toRequest()
				return m, executeRequest(req)
			}
			return m, nil

		case "tab":
			// Cycle forward through fields
			m.editing = false
			m.focusedField = (m.focusedField + 1) % 4
			m.cursorPos = 0
			// Auto-start editing on URL and Body fields
			if m.focusedField == URLField || m.focusedField == BodyField {
				m.editing = true
				m.cursorPos = len(m.url)
				if m.focusedField == BodyField {
					m.cursorPos = len(m.body)
				}
			}
			return m, nil

		case "shift+tab":
			// Cycle backward through fields
			m.editing = false
			m.focusedField = (m.focusedField + 3) % 4 // +3 is same as -1 mod 4
			m.cursorPos = 0
			// Auto-start editing on URL and Body fields
			if m.focusedField == URLField || m.focusedField == BodyField {
				m.editing = true
				m.cursorPos = len(m.url)
				if m.focusedField == BodyField {
					m.cursorPos = len(m.body)
				}
			}
			return m, nil

		case "esc":
			// Exit editing mode
			m.editing = false
			return m, nil

		case "1", "2", "3":
			// Switch response tabs (only when not editing)
			if !m.editing && m.response.StatusCode != 0 {
				switch msg.String() {
				case "1":
					m.responseTab = 0
				case "2":
					m.responseTab = 1
				case "3":
					m.responseTab = 2
				}
			}
			return m, nil
		}

		// Field-specific handling
		switch m.focusedField {
		case MethodField:
			return handleMethodField(m, msg)
		case URLField:
			return handleURLField(m, msg)
		case HeadersField:
			return handleHeadersField(m, msg)
		case BodyField:
			return handleBodyField(m, msg)
		}
	}

	return m, nil
}

// handlePasteMode handles input in paste mode
func handlePasteMode(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "ctrl+c":
		// Cancel paste mode
		m.pasteMode = false
		m.pasteBuffer = ""
		return m, nil

	case "enter":
		// Parse the pasted curl command
		if m.pasteBuffer != "" {
			req, err := curl.ParseCurlCommand(m.pasteBuffer)
			if err != nil {
				m.err = err
			} else {
				m.method = req.Method
				m.url = req.URL
				m.body = req.Body
				m.headers = []HeaderPair{}
				for key, value := range req.Headers {
					m.headers = append(m.headers, HeaderPair{Key: key, Value: value})
				}
				m.err = nil
			}
		}
		m.pasteMode = false
		m.pasteBuffer = ""
		return m, nil

	case "backspace":
		if len(m.pasteBuffer) > 0 {
			m.pasteBuffer = m.pasteBuffer[:len(m.pasteBuffer)-1]
		}
		return m, nil

	default:
		// Add character to paste buffer
		if len(msg.String()) == 1 {
			m.pasteBuffer += msg.String()
		}
		return m, nil
	}
}

// handleMethodField handles input when method field is focused
func handleMethodField(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}

	switch msg.String() {
	case "left", "h":
		// Cycle to previous method
		for i, method := range methods {
			if method == m.method {
				if i > 0 {
					m.method = methods[i-1]
				} else {
					m.method = methods[len(methods)-1]
				}
				break
			}
		}
		return m, nil

	case "right", "l":
		// Cycle to next method
		for i, method := range methods {
			if method == m.method {
				if i < len(methods)-1 {
					m.method = methods[i+1]
				} else {
					m.method = methods[0]
				}
				break
			}
		}
		return m, nil
	}

	return m, nil
}

// handleURLField handles input when URL field is focused
func handleURLField(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Auto-enable editing when focused on URL field
	if !m.editing {
		m.editing = true
	}

	switch msg.String() {
	case "backspace":
		if len(m.url) > 0 && m.cursorPos > 0 {
			m.url = m.url[:m.cursorPos-1] + m.url[m.cursorPos:]
			m.cursorPos--
		}
		return m, nil

	case "delete":
		if m.cursorPos < len(m.url) {
			m.url = m.url[:m.cursorPos] + m.url[m.cursorPos+1:]
		}
		return m, nil

	case "left":
		if m.cursorPos > 0 {
			m.cursorPos--
		}
		return m, nil

	case "right":
		if m.cursorPos < len(m.url) {
			m.cursorPos++
		}
		return m, nil

	case "home", "ctrl+a":
		m.cursorPos = 0
		return m, nil

	case "end":
		m.cursorPos = len(m.url)
		return m, nil

	default:
		// Add character
		if len(msg.String()) == 1 {
			m.url = m.url[:m.cursorPos] + msg.String() + m.url[m.cursorPos:]
			m.cursorPos++
		}
		return m, nil
	}
}

// handleHeadersField handles input when headers field is focused
func handleHeadersField(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if !m.editing && m.headerIndex < len(m.headers) {
			m.headerIndex++
		}
		return m, nil

	case "k", "up":
		if !m.editing && m.headerIndex > 0 {
			m.headerIndex--
		}
		return m, nil

	case "a":
		if !m.editing {
			// Add new header
			m.headers = append(m.headers, HeaderPair{Key: "", Value: ""})
			m.headerIndex = len(m.headers) - 1
			m.editing = true
			m.headerEditField = 0
		}
		return m, nil

	case "e", "enter":
		if !m.editing && m.headerIndex < len(m.headers) {
			// Edit existing header
			m.editing = true
			m.headerEditField = 0
		}
		return m, nil

	case "d":
		if !m.editing && m.headerIndex < len(m.headers) {
			// Delete header
			m.headers = append(m.headers[:m.headerIndex], m.headers[m.headerIndex+1:]...)
			if m.headerIndex >= len(m.headers) && m.headerIndex > 0 {
				m.headerIndex--
			}
		}
		return m, nil

	case "tab":
		if m.editing {
			// Switch between key and value
			m.headerEditField = 1 - m.headerEditField
			return m, nil
		}

	case "backspace":
		if m.editing && m.headerIndex < len(m.headers) {
			if m.headerEditField == 0 && len(m.headers[m.headerIndex].Key) > 0 {
				m.headers[m.headerIndex].Key = m.headers[m.headerIndex].Key[:len(m.headers[m.headerIndex].Key)-1]
			} else if m.headerEditField == 1 && len(m.headers[m.headerIndex].Value) > 0 {
				m.headers[m.headerIndex].Value = m.headers[m.headerIndex].Value[:len(m.headers[m.headerIndex].Value)-1]
			}
		}
		return m, nil

	default:
		// Add character
		if m.editing && m.headerIndex < len(m.headers) && len(msg.String()) == 1 {
			if m.headerEditField == 0 {
				m.headers[m.headerIndex].Key += msg.String()
			} else {
				m.headers[m.headerIndex].Value += msg.String()
			}
		}
		return m, nil
	}

	return m, nil
}

// handleBodyField handles input when body field is focused
func handleBodyField(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Auto-enable editing when focused on body field
	if !m.editing {
		m.editing = true
	}

	switch msg.String() {
	case "enter":
		// Add newline
		m.body = m.body[:m.cursorPos] + "\n" + m.body[m.cursorPos:]
		m.cursorPos++
		return m, nil

	case "backspace":
		if len(m.body) > 0 && m.cursorPos > 0 {
			m.body = m.body[:m.cursorPos-1] + m.body[m.cursorPos:]
			m.cursorPos--
		}
		return m, nil

	case "delete":
		if m.cursorPos < len(m.body) {
			m.body = m.body[:m.cursorPos] + m.body[m.cursorPos+1:]
		}
		return m, nil

	case "left":
		if m.cursorPos > 0 {
			m.cursorPos--
		}
		return m, nil

	case "right":
		if m.cursorPos < len(m.body) {
			m.cursorPos++
		}
		return m, nil

	case "up":
		// Move cursor to previous line
		lines := strings.Split(m.body[:m.cursorPos], "\n")
		if len(lines) > 1 {
			currentLinePos := len(lines[len(lines)-1])
			prevLineLen := len(lines[len(lines)-2])
			m.cursorPos -= currentLinePos + 1 // +1 for newline
			if currentLinePos > prevLineLen {
				m.cursorPos -= currentLinePos - prevLineLen
			}
		}
		return m, nil

	case "down":
		// Move cursor to next line
		remaining := m.body[m.cursorPos:]
		lines := strings.Split(remaining, "\n")
		if len(lines) > 1 {
			currentLineRemaining := len(lines[0])
			m.cursorPos += currentLineRemaining + 1 // +1 for newline
			nextLineLen := len(lines[1])
			if currentLineRemaining > nextLineLen {
				m.cursorPos += nextLineLen
			}
		}
		return m, nil

	case "home", "ctrl+a":
		m.cursorPos = 0
		return m, nil

	case "end":
		m.cursorPos = len(m.body)
		return m, nil

	default:
		// Add character
		if len(msg.String()) == 1 {
			m.body = m.body[:m.cursorPos] + msg.String() + m.body[m.cursorPos:]
			m.cursorPos++
		}
		return m, nil
	}
}
