package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/username/curltui/internal/curl"
	"github.com/username/curltui/internal/http"
)

// Field represents which field is currently focused
type Field int

const (
	MethodField Field = iota
	URLField
	HeadersField
	BodyField
)

// HeaderPair represents a key-value pair for headers (to maintain order)
type HeaderPair struct {
	Key   string
	Value string
}

// Model represents the application state
type Model struct {
	// Request state
	method  string
	url     string
	headers []HeaderPair
	body    string

	// UI state
	focusedField    Field
	cursorPos       int // cursor position within text fields
	headerIndex     int // which header is selected
	editing         bool
	loading         bool
	pasteMode       bool
	pasteBuffer     string
	responseTab     int // 0 = body, 1 = headers
	bodyScroll      int
	responseScroll  int
	headerEditField int // 0 = key, 1 = value

	// Response state
	response http.Response
	err      error

	// Window dimensions
	width  int
	height int
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

// toRequest converts the model state to a curl.Request
func (m Model) toRequest() curl.Request {
	req := curl.NewRequest()
	req.Method = m.method
	req.URL = m.url
	req.Body = m.body

	for _, h := range m.headers {
		if h.Key != "" {
			req.Headers[h.Key] = h.Value
		}
	}

	return req
}
