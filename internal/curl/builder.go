package curl

import "strings"

// BuildCurlCommand converts a Request struct to a valid curl command string
func BuildCurlCommand(req Request) string {
	var parts []string
	parts = append(parts, "curl")

	// Always add -X flag with method
	parts = append(parts, "-X", req.Method)

	// Add headers
	for key, value := range req.Headers {
		// Escape double quotes in header values
		escapedValue := strings.ReplaceAll(value, `"`, `\"`)
		parts = append(parts, "-H", `"`+key+": "+escapedValue+`"`)
	}

	// Add body if non-empty
	if req.Body != "" {
		// Use single quotes for body and escape single quotes within
		escapedBody := strings.ReplaceAll(req.Body, `'`, `'\''`)
		parts = append(parts, "-d", `'`+escapedBody+`'`)
	}

	// Add URL (quoted)
	if req.URL != "" {
		parts = append(parts, `"`+req.URL+`"`)
	}

	return strings.Join(parts, " ")
}
