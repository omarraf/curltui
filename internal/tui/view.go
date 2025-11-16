package tui

import (
	"fmt"
	"strings"
)

// renderView renders the complete TUI layout
// TODO: Implement method selector
//   - Display current HTTP method (GET, POST, PUT, PATCH, DELETE, etc.)
//   - Highlight when focused
//   - Show indicator for cycling through methods
//   - Style: [GET ▼] or similar
//
// TODO: Implement URL input field
//   - Display current URL
//   - Show cursor when focused
//   - Support horizontal scrolling for long URLs
//   - Placeholder text when empty: "https://api.example.com/endpoint"
//
// TODO: Implement headers list editor
//   - Display list of headers as "Key: Value" pairs
//   - Allow adding new headers
//   - Allow editing existing headers
//   - Allow deleting headers
//   - Highlight focused header
//   - Show "Add Header" option when focused
//
// TODO: Implement body editor
//   - Multi-line text input area
//   - Support for JSON formatting/syntax highlighting
//   - Show line numbers
//   - Scroll support for large bodies
//   - Highlight when focused
//
// TODO: Implement curl preview pane
//   - Display the generated curl command
//   - Use curl.BuildCurlCommand to generate
//   - Syntax highlighting for better readability
//   - Allow copying to clipboard
//   - Update in real-time as request changes
//
// TODO: Implement response viewer pane
//   - Display HTTP status code and reason
//   - Display response headers
//   - Display response body (formatted if JSON)
//   - Show request duration
//   - Syntax highlighting for JSON/XML/HTML
//   - Scroll support for large responses
//   - Show loading state during request execution
//
// TODO: Layout structure (split screen)
//   - Left side: Request builder (method, URL, headers, body)
//   - Right side: Top: curl preview, Bottom: response viewer
//   - Use borders to separate sections
//   - Add section titles
//   - Show keyboard shortcuts in footer
func renderView(m Model) string {
	if m.width == 0 {
		return "Loading..."
	}

	var b strings.Builder

	// Header
	b.WriteString(titleStyle.Render("curltui - HTTP Request Builder"))
	b.WriteString("\n\n")

	// Basic placeholder UI
	b.WriteString("Method: " + m.request.Method + "\n")
	b.WriteString("URL: " + m.request.URL + "\n")
	b.WriteString("\n")

	// TODO: Implement full UI layout with all sections

	// Footer with help
	footer := helpStyle.Render("Press q or Ctrl+C to quit")
	b.WriteString("\n")
	b.WriteString(footer)

	return b.String()
}

// Helper function to render a section with border
// TODO: Implement proper section rendering with lipgloss
func renderSection(title, content string, focused bool) string {
	return fmt.Sprintf("[%s]\n%s", title, content)
}
