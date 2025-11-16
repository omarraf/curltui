package curl

import (
	"fmt"
	"strings"
)

// ParseCurlCommand parses a curl command string into a Request struct
func ParseCurlCommand(cmd string) (Request, error) {
	req := NewRequest()

	if cmd == "" {
		return req, fmt.Errorf("empty curl command")
	}

	// Tokenize the command
	tokens, err := tokenize(cmd)
	if err != nil {
		return req, err
	}

	if len(tokens) == 0 {
		return req, fmt.Errorf("empty curl command")
	}

	// Skip "curl" or "curl.exe" prefix
	if tokens[0] == "curl" || tokens[0] == "curl.exe" {
		tokens = tokens[1:]
	}

	// Parse flags and arguments
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		// Handle method flags
		if token == "-X" || token == "--request" {
			if i+1 >= len(tokens) {
				return req, fmt.Errorf("missing value for %s flag", token)
			}
			req.Method = tokens[i+1]
			i++
			continue
		}

		// Handle header flags
		if token == "-H" || token == "--header" {
			if i+1 >= len(tokens) {
				return req, fmt.Errorf("missing value for %s flag", token)
			}
			header := tokens[i+1]
			// Split on first colon
			parts := strings.SplitN(header, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				req.Headers[key] = value
			}
			i++
			continue
		}

		// Handle data flags
		if token == "-d" || token == "--data" || token == "--data-raw" || token == "--data-binary" {
			if i+1 >= len(tokens) {
				return req, fmt.Errorf("missing value for %s flag", token)
			}
			req.Body = tokens[i+1]
			i++
			continue
		}

		// Handle combined flags like -XPOST
		if strings.HasPrefix(token, "-X") && len(token) > 2 {
			req.Method = token[2:]
			continue
		}

		if strings.HasPrefix(token, "-H") && len(token) > 2 {
			header := token[2:]
			parts := strings.SplitN(header, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				req.Headers[key] = value
			}
			continue
		}

		if strings.HasPrefix(token, "-d") && len(token) > 2 {
			req.Body = token[2:]
			continue
		}

		// Skip other flags we don't care about
		if strings.HasPrefix(token, "-") {
			// Skip flag and its value if it takes one
			if i+1 < len(tokens) && !strings.HasPrefix(tokens[i+1], "-") {
				i++
			}
			continue
		}

		// Assume it's a URL if it doesn't start with a dash
		if !strings.HasPrefix(token, "-") {
			req.URL = token
		}
	}

	return req, nil
}

// tokenize splits a curl command into tokens, respecting quotes
func tokenize(cmd string) ([]string, error) {
	var tokens []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	escaped := false

	cmd = strings.TrimSpace(cmd)

	for i := 0; i < len(cmd); i++ {
		ch := cmd[i]

		if escaped {
			current.WriteByte(ch)
			escaped = false
			continue
		}

		if ch == '\\' {
			escaped = true
			continue
		}

		if ch == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			continue
		}

		if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			continue
		}

		if (ch == ' ' || ch == '\t' || ch == '\n') && !inSingleQuote && !inDoubleQuote {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteByte(ch)
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	if inSingleQuote || inDoubleQuote {
		return nil, fmt.Errorf("unclosed quote in curl command")
	}

	return tokens, nil
}
