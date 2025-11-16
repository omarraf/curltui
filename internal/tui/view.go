package tui

import (
	"encoding/json"
	"fmt"
	"regexp"
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

	// Status bar at top
	b.WriteString(renderStatusBar(m))
	b.WriteString("\n")

	// Calculate layout dimensions
	leftWidth := m.width / 2
	rightWidth := m.width - leftWidth - 1
	contentHeight := m.height - 6 // Reserve space for status bar and help footer

	// Split screen layout
	leftSide := renderLeftSide(m, leftWidth, contentHeight)
	rightSide := renderRightSide(m, rightWidth, contentHeight)

	// Combine sides
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, leftSide, rightSide))
	b.WriteString("\n")

	// Footer with help
	footer := renderFooter(m)
	b.WriteString(footer)

	return b.String()
}

// renderStatusBar renders the top status bar
func renderStatusBar(m Model) string {
	var status string

	mode := "NAVIGATE"
	if m.editing {
		mode = "EDIT"
	}
	if m.loading {
		mode = "SENDING REQUEST..."
	}

	status = fmt.Sprintf(" %s ", mode)

	// Add last request info if available
	if m.response.StatusCode != 0 {
		statusCode := fmt.Sprintf("Last: %dms • %d", m.response.Duration.Milliseconds(), m.response.StatusCode)
		status += dimStyle.Render(" | ") + statusCode
	}

	// Pad to full width
	padding := m.width - lipgloss.Width(status)
	if padding > 0 {
		status += strings.Repeat(" ", padding)
	}

	return statusBarStyle.Render(status)
}

// renderLeftSide renders the request builder section
func renderLeftSide(m Model, width, height int) string {
	var sections []string

	// Method and URL row
	methodStr := renderMethodField(m)
	urlStr := renderURLField(m, width-18)
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, methodStr, urlStr)
	sections = append(sections, topRow)

	// Headers section (dynamic height)
	headersHeight := max(6, min(height/4, len(m.headers)+4))
	headersStr := renderHeadersField(m, width-1, headersHeight)
	sections = append(sections, headersStr)

	// Body section (remaining space)
	bodyHeight := height - headersHeight - 5
	bodyStr := renderBodyField(m, width-1, bodyHeight)
	sections = append(sections, bodyStr)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderRightSide renders the curl preview and response viewer
func renderRightSide(m Model, width, height int) string {
	var sections []string

	// Curl preview (dynamic height based on command length)
	curlHeight := 6
	curlStr := renderCurlPreview(m, width-1, curlHeight)
	sections = append(sections, curlStr)

	// Response viewer (remaining space)
	responseHeight := height - curlHeight - 1
	responseStr := renderResponse(m, width-1, responseHeight)
	sections = append(sections, responseStr)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderMethodField renders the method selector showing all methods
func renderMethodField(m Model) string {
	focused := m.focusedField == MethodField
	style := sectionStyle
	if focused {
		style = focusedSectionStyle
	}

	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

	var methodsDisplay strings.Builder
	for _, method := range methods {
		if method == m.method {
			methodsDisplay.WriteString(methodActiveStyle.Render(method))
		} else {
			methodsDisplay.WriteString(methodInactiveStyle.Render(method))
		}
		methodsDisplay.WriteString(" ")
	}

	// Add arrow indicator for focused field
	indicator := ""
	if focused {
		for i, method := range methods {
			if method == m.method {
				spaces := ""
				for j := 0; j < i; j++ {
					spaces += strings.Repeat(" ", len(methods[j])+1)
				}
				indicator = "\n" + spaces + dimStyle.Render("^^^")
				break
			}
		}
	}

	content := methodsDisplay.String() + indicator

	return style.Width(16).Render(content)
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
		url = placeholderStyle.Render("https://api.example.com/endpoint")
	}

	// Show cursor if focused
	if focused && m.cursorPos <= len(m.url) {
		if m.cursorPos < len(m.url) {
			url = m.url[:m.cursorPos] + cursorStyle.Render("│") + m.url[m.cursorPos:]
		} else {
			url = m.url + cursorStyle.Render("│")
		}
	}

	title := dimStyle.Render("URL")
	content := title + "\n" + url

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
	content.WriteString(dimStyle.Render("Headers"))
	content.WriteString("\n")

	if len(m.headers) == 0 {
		content.WriteString("\n")
		content.WriteString(placeholderStyle.Render("No headers"))
		if focused {
			content.WriteString("\n")
			content.WriteString(dimStyle.Render("'a' to add • 'i' for common headers"))
		}
	} else {
		for i, h := range m.headers {
			prefix := "  "
			if focused && i == m.headerIndex {
				prefix = highlightStyle.Render("> ")
			}

			if focused && i == m.headerIndex && m.editing {
				key := h.Key
				value := h.Value
				if m.headerEditField == 0 {
					key = key + cursorStyle.Render("│")
				} else {
					value = value + cursorStyle.Render("│")
				}
				content.WriteString(fmt.Sprintf("\n%s%s: %s", prefix, jsonKeyStyle.Render(key), value))
			} else {
				content.WriteString(fmt.Sprintf("\n%s%s: %s", prefix, dimStyle.Render(h.Key), h.Value))
			}
		}
	}

	return style.Width(width).Height(height).Render(content.String())
}

// renderBodyField renders the body editor
func renderBodyField(m Model, width, height int) string {
	focused := m.focusedField == BodyField
	style := sectionStyle
	if focused {
		style = focusedSectionStyle
	}

	var content strings.Builder

	// Title with validation indicator
	title := dimStyle.Render("Body")
	if m.body != "" {
		if isJSON(m.body) {
			title += " " + successStatusStyle.Render("✓")
		} else {
			title += " " + errorMessageStyle.Render("✗")
		}
	}
	content.WriteString(title)
	content.WriteString("\n")

	body := m.body
	if body == "" && !focused {
		content.WriteString("\n")
		content.WriteString(placeholderStyle.Render(`{"key": "value"}`))
	} else if focused && m.cursorPos <= len(m.body) {
		// Show cursor
		if m.cursorPos < len(m.body) {
			body = m.body[:m.cursorPos] + cursorStyle.Render("│") + m.body[m.cursorPos:]
		} else {
			body = m.body + cursorStyle.Render("│")
		}
		content.WriteString("\n")
		content.WriteString(body)
	} else {
		// Syntax highlight JSON
		content.WriteString("\n")
		if isJSON(body) {
			content.WriteString(highlightJSON(body))
		} else {
			content.WriteString(body)
		}
	}

	return style.Width(width).Height(height).Render(content.String())
}

// renderCurlPreview renders the curl command preview with syntax highlighting
func renderCurlPreview(m Model, width, height int) string {
	var content strings.Builder
	content.WriteString(dimStyle.Render("Curl Command"))
	content.WriteString(dimStyle.Render(" (Ctrl+C to copy)"))
	content.WriteString("\n\n")

	req := m.toRequest()
	curlCmd := curl.BuildCurlCommand(req)

	// Syntax highlight the curl command
	highlighted := highlightCurlCommand(curlCmd)
	content.WriteString(highlighted)

	return sectionStyle.Width(width).Height(height).Render(content.String())
}

// renderResponse renders the response viewer with tabs
func renderResponse(m Model, width, height int) string {
	var content strings.Builder

	// Tab bar
	tabs := []string{"Body", "Headers", "Info"}
	var tabBar strings.Builder
	for i, tab := range tabs {
		if i == m.responseTab {
			tabBar.WriteString(tabActiveStyle.Render(tab))
		} else {
			tabBar.WriteString(tabInactiveStyle.Render(tab))
		}
		if i < len(tabs)-1 {
			tabBar.WriteString("  ")
		}
	}
	content.WriteString(tabBar.String())
	content.WriteString("\n")
	content.WriteString(strings.Repeat("─", width-4))
	content.WriteString("\n")

	if m.loading {
		content.WriteString("\n")
		content.WriteString(loadingStyle.Render("⣾ Sending request..."))
	} else if m.err != nil {
		content.WriteString("\n")
		content.WriteString(errorMessageStyle.Render("Error"))
		content.WriteString("\n")
		content.WriteString(errorMessageStyle.Render(m.err.Error()))
	} else if m.response.StatusCode == 0 {
		content.WriteString("\n")
		content.WriteString(placeholderStyle.Render("No response yet. Press Ctrl+E to execute."))
	} else {
		// Status line
		statusStyle := getStatusStyle(m.response.StatusCode)
		statusLine := fmt.Sprintf("%d", m.response.StatusCode)
		content.WriteString("\n")
		content.WriteString(statusStyle.Render(statusLine))
		content.WriteString(dimStyle.Render(fmt.Sprintf(" • %s • %d bytes",
			m.response.Duration, len(m.response.Body))))
		content.WriteString("\n\n")

		// Tab content
		switch m.responseTab {
		case 0: // Body
			body := m.response.Body
			if isJSON(body) {
				var formatted interface{}
				if err := json.Unmarshal([]byte(body), &formatted); err == nil {
					if pretty, err := json.MarshalIndent(formatted, "", "  "); err == nil {
						body = highlightJSON(string(pretty))
					}
				}
			}
			content.WriteString(body)

		case 1: // Headers
			for key, value := range m.response.Headers {
				content.WriteString(jsonKeyStyle.Render(key))
				content.WriteString(": ")
				content.WriteString(value)
				content.WriteString("\n")
			}

		case 2: // Info
			contentType := m.response.Headers["Content-Type"]
			if contentType == "" {
				contentType = "unknown"
			}
			content.WriteString(dimStyle.Render("Content-Type: "))
			content.WriteString(contentType)
			content.WriteString("\n")
			content.WriteString(dimStyle.Render("Size: "))
			content.WriteString(fmt.Sprintf("%d bytes", len(m.response.Body)))
			content.WriteString("\n")
			content.WriteString(dimStyle.Render("Duration: "))
			content.WriteString(fmt.Sprintf("%s", m.response.Duration))
		}
	}

	return sectionStyle.Width(width).Height(height).Render(content.String())
}

// renderPasteMode renders the paste mode UI
func renderPasteMode(m Model) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Paste curl command"))
	b.WriteString("\n\n")

	content := m.pasteBuffer
	if content == "" {
		content = placeholderStyle.Render("Paste your curl command here...")
	}

	b.WriteString(sectionStyle.Width(m.width - 4).Render(content))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Press Enter to parse, Esc to cancel"))

	return b.String()
}

// renderFooter renders the help footer
func renderFooter(m Model) string {
	var helpText string

	sep := dimStyle.Render(" • ")

	if m.focusedField == MethodField {
		helpText = highlightStyle.Render("←/→") + ":method" + sep +
			"Tab:next" + sep +
			"Ctrl+E:execute" + sep +
			"Ctrl+P:paste" + sep +
			"q:quit"
	} else if m.focusedField == URLField {
		helpText = "Type URL" + sep +
			"Tab:next" + sep +
			"Ctrl+E:execute" + sep +
			"Ctrl+P:paste" + sep +
			"q:quit"
	} else if m.focusedField == HeadersField {
		helpText = highlightStyle.Render("a") + ":add" + sep +
			highlightStyle.Render("e") + ":edit" + sep +
			highlightStyle.Render("d") + ":delete" + sep +
			highlightStyle.Render("j/k") + ":navigate" + sep +
			"Tab:next" + sep +
			"Ctrl+E:execute"
	} else if m.focusedField == BodyField {
		helpText = "Type body" + sep +
			highlightStyle.Render("Enter") + ":newline" + sep +
			"Tab:next" + sep +
			"Ctrl+E:execute"
	}

	return helpStyle.Render(helpText)
}

// getStatusStyle returns the appropriate style for HTTP status code
func getStatusStyle(code int) lipgloss.Style {
	if code >= 200 && code < 300 {
		return successStatusStyle
	} else if code >= 300 && code < 400 {
		return redirectStatusStyle
	} else if code >= 400 && code < 500 {
		return clientErrorStatusStyle
	} else {
		return serverErrorStatusStyle
	}
}

// highlightCurlCommand applies syntax highlighting to curl command
func highlightCurlCommand(cmd string) string {
	// Highlight curl command
	cmd = regexp.MustCompile(`^curl\b`).ReplaceAllStringFunc(cmd, func(s string) string {
		return curlCmdStyle.Render(s)
	})

	// Highlight HTTP methods
	cmd = regexp.MustCompile(`\b(GET|POST|PUT|PATCH|DELETE|HEAD|OPTIONS)\b`).ReplaceAllStringFunc(cmd, func(s string) string {
		return curlMethodStyle.Render(s)
	})

	// Highlight flags
	cmd = regexp.MustCompile(`\s(-[A-Za-z]|--[a-z-]+)\b`).ReplaceAllStringFunc(cmd, func(s string) string {
		return " " + curlFlagStyle.Render(strings.TrimSpace(s))
	})

	// Highlight quoted strings (both single and double quotes)
	cmd = regexp.MustCompile(`"[^"]*"`).ReplaceAllStringFunc(cmd, func(s string) string {
		return curlStringStyle.Render(s)
	})
	cmd = regexp.MustCompile(`'[^']*'`).ReplaceAllStringFunc(cmd, func(s string) string {
		return curlStringStyle.Render(s)
	})

	// Highlight URLs
	cmd = regexp.MustCompile(`https?://[^\s"']+`).ReplaceAllStringFunc(cmd, func(s string) string {
		return curlURLStyle.Render(s)
	})

	return cmd
}

// isJSON checks if a string is valid JSON
func isJSON(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

// highlightJSON applies syntax highlighting to JSON
func highlightJSON(s string) string {
	// Highlight keys
	s = regexp.MustCompile(`"([^"]+)":`).ReplaceAllStringFunc(s, func(match string) string {
		return jsonKeyStyle.Render(match[:len(match)-1]) + ":"
	})

	// Highlight string values
	s = regexp.MustCompile(`: "([^"]*)"`).ReplaceAllStringFunc(s, func(match string) string {
		return ": " + jsonStringStyle.Render(match[2:])
	})

	// Highlight numbers
	s = regexp.MustCompile(`:\s*(-?\d+\.?\d*)`).ReplaceAllStringFunc(s, func(match string) string {
		parts := strings.Split(match, ":")
		return ": " + jsonNumberStyle.Render(strings.TrimSpace(parts[1]))
	})

	// Highlight booleans
	s = regexp.MustCompile(`:\s*(true|false)`).ReplaceAllStringFunc(s, func(match string) string {
		parts := strings.Split(match, ":")
		return ": " + jsonBoolStyle.Render(strings.TrimSpace(parts[1]))
	})

	// Highlight null
	s = regexp.MustCompile(`:\s*null`).ReplaceAllStringFunc(s, func(match string) string {
		return ": " + jsonNullStyle.Render("null")
	})

	return s
}

// Helper functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
