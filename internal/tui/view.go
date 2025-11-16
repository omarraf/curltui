package tui

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/username/curltui/internal/curl"
)

// renderView renders the complete TUI layout
func renderView(m Model) string {
	if m.width == 0 {
		return "Loading..."
	}

	// Handle paste mode
	if m.pasteMode {
		return renderPasteMode(m)
	}

	var b strings.Builder

	// Header
	b.WriteString(titleStyle.Render("curltui - HTTP Request Builder"))
	b.WriteString("\n\n")

	// Calculate layout dimensions
	leftWidth := m.width / 2
	rightWidth := m.width - leftWidth - 2
	contentHeight := m.height - 10 // Reserve space for header and footer

	// Split screen layout
	leftSide := renderLeftSide(m, leftWidth, contentHeight)
	rightSide := renderRightSide(m, rightWidth, contentHeight)

	// Combine sides
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, leftSide, rightSide))
	b.WriteString("\n\n")

	// Footer with help
	footer := renderFooter(m)
	b.WriteString(footer)

	return b.String()
}

// renderLeftSide renders the request builder section
func renderLeftSide(m Model, width, height int) string {
	var sections []string

	// Method and URL row
	methodStr := renderMethodField(m)
	urlStr := renderURLField(m, width-20)
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, methodStr, " ", urlStr)
	sections = append(sections, topRow)

	// Headers section
	headersStr := renderHeadersField(m, width, height/3)
	sections = append(sections, headersStr)

	// Body section
	bodyStr := renderBodyField(m, width, height/3)
	sections = append(sections, bodyStr)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderRightSide renders the curl preview and response viewer
func renderRightSide(m Model, width, height int) string {
	var sections []string

	// Curl preview
	curlStr := renderCurlPreview(m, width, height/3)
	sections = append(sections, curlStr)

	// Response viewer
	responseStr := renderResponse(m, width, height*2/3)
	sections = append(sections, responseStr)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderMethodField renders the method selector
func renderMethodField(m Model) string {
	focused := m.focusedField == MethodField
	style := sectionStyle
	if focused {
		style = focusedSectionStyle
	}

	content := fmt.Sprintf("< %s >", m.method)
	if focused {
		content = fmt.Sprintf("◄ %s ►", m.method)
	}

	return style.Width(15).Render(content)
}

// renderURLField renders the URL input field
func renderURLField(m Model, width int) string {
	focused := m.focusedField == URLField
	style := sectionStyle
	if focused {
		style = focusedSectionStyle
	}

	url := m.url
	if url == "" && !focused {
		url = dimStyle.Render("https://api.example.com/endpoint")
	}

	// Show cursor if editing
	if focused && m.editing && m.cursorPos <= len(m.url) {
		if m.cursorPos < len(m.url) {
			url = m.url[:m.cursorPos] + "│" + m.url[m.cursorPos:]
		} else {
			url = m.url + "│"
		}
	}

	title := "URL"
	if focused {
		title = "URL (type to edit, Tab to navigate)"
	}

	content := fmt.Sprintf("%s\n%s", dimStyle.Render(title), url)
	return style.Width(width).Render(content)
}

// renderHeadersField renders the headers editor
func renderHeadersField(m Model, width, height int) string {
	focused := m.focusedField == HeadersField
	style := sectionStyle
	if focused {
		style = focusedSectionStyle
	}

	var content strings.Builder
	content.WriteString(dimStyle.Render("Headers") + "\n\n")

	if len(m.headers) == 0 {
		content.WriteString(dimStyle.Render("No headers. Press 'a' to add."))
	} else {
		for i, h := range m.headers {
			prefix := "  "
			if focused && i == m.headerIndex {
				prefix = "> "
			}

			if focused && i == m.headerIndex && m.editing {
				key := h.Key
				value := h.Value
				if m.headerEditField == 0 {
					key = key + "│"
				} else {
					value = value + "│"
				}
				content.WriteString(fmt.Sprintf("%s%s: %s\n", prefix, key, value))
			} else {
				content.WriteString(fmt.Sprintf("%s%s: %s\n", prefix, h.Key, h.Value))
			}
		}
	}

	if focused && !m.editing {
		content.WriteString("\n" + dimStyle.Render("a:add e:edit d:delete"))
	}

	return style.Width(width - 2).Height(height).Render(content.String())
}

// renderBodyField renders the body editor
func renderBodyField(m Model, width, height int) string {
	focused := m.focusedField == BodyField
	style := sectionStyle
	if focused {
		style = focusedSectionStyle
	}

	var content strings.Builder
	title := "Body"
	if focused {
		title = "Body (type to edit, Enter for newlines)"
	}
	content.WriteString(dimStyle.Render(title) + "\n\n")

	body := m.body
	if body == "" && !focused {
		body = dimStyle.Render("{\"key\": \"value\"}")
	} else if focused && m.editing && m.cursorPos <= len(m.body) {
		// Show cursor
		if m.cursorPos < len(m.body) {
			body = m.body[:m.cursorPos] + "│" + m.body[m.cursorPos:]
		} else {
			body = m.body + "│"
		}
	}

	content.WriteString(body)

	return style.Width(width - 2).Height(height).Render(content.String())
}

// renderCurlPreview renders the curl command preview
func renderCurlPreview(m Model, width, height int) string {
	var content strings.Builder
	content.WriteString(dimStyle.Render("Curl Command") + "\n\n")

	req := m.toRequest()
	curlCmd := curl.BuildCurlCommand(req)

	content.WriteString(codeStyle.Render(curlCmd))

	return sectionStyle.Width(width - 2).Height(height).Render(content.String())
}

// renderResponse renders the response viewer
func renderResponse(m Model, width, height int) string {
	var content strings.Builder
	content.WriteString(dimStyle.Render("Response") + "\n\n")

	if m.loading {
		content.WriteString(primaryStyle.Render("Loading..."))
	} else if m.err != nil {
		content.WriteString(errorStyle.Render("Error: " + m.err.Error()))
	} else if m.response.StatusCode == 0 {
		content.WriteString(dimStyle.Render("No response yet. Press Ctrl+E to execute."))
	} else {
		// Status code with color
		statusStyle := successStyle
		if m.response.StatusCode >= 400 {
			statusStyle = errorStyle
		} else if m.response.StatusCode >= 300 {
			statusStyle = warningStyle
		}

		content.WriteString(statusStyle.Render(fmt.Sprintf("Status: %d", m.response.StatusCode)))
		content.WriteString(dimStyle.Render(fmt.Sprintf(" | Duration: %s", m.response.Duration)))
		content.WriteString("\n\n")

		// Response body
		body := m.response.Body
		// Try to pretty print JSON
		if isJSON(body) {
			var formatted interface{}
			if err := json.Unmarshal([]byte(body), &formatted); err == nil {
				if pretty, err := json.MarshalIndent(formatted, "", "  "); err == nil {
					body = highlightJSON(string(pretty))
				}
			}
		}

		content.WriteString(body)
	}

	return sectionStyle.Width(width - 2).Height(height).Render(content.String())
}

// renderPasteMode renders the paste mode UI
func renderPasteMode(m Model) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Paste curl command"))
	b.WriteString("\n\n")

	content := m.pasteBuffer
	if content == "" {
		content = dimStyle.Render("Paste your curl command here...")
	}

	b.WriteString(sectionStyle.Width(m.width - 4).Render(content))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Press Enter to parse, Esc to cancel"))

	return b.String()
}

// renderFooter renders the help footer
func renderFooter(m Model) string {
	var helpText string
	if m.focusedField == MethodField {
		helpText = "←/→:change method • Tab:next field • Ctrl+E:execute • Ctrl+P:paste curl • q:quit"
	} else if m.focusedField == URLField {
		helpText = "Type to edit URL • Tab:next field • Ctrl+E:execute • Ctrl+P:paste curl • q:quit"
	} else if m.focusedField == HeadersField {
		helpText = "a:add header • e:edit • d:delete • j/k:navigate • Tab:next field • Ctrl+E:execute • q:quit"
	} else if m.focusedField == BodyField {
		helpText = "Type to edit body • Enter:newline • Tab:next field • Ctrl+E:execute • q:quit"
	}

	return helpStyle.Render(helpText)
}

// isJSON checks if a string is valid JSON
func isJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

// highlightJSON applies simple syntax highlighting to JSON
func highlightJSON(s string) string {
	// Simple highlighting: this could be much more sophisticated
	s = strings.ReplaceAll(s, "\"", jsonStringStyle.Render("\""))
	return s
}

// Additional styles for JSON highlighting
var (
	jsonStringStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	primaryStyle    = lipgloss.NewStyle().Foreground(primaryColor)
	dimStyle        = lipgloss.NewStyle().Foreground(mutedColor)
	warningStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
)
