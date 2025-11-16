package http

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Response represents an HTTP response
type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       string
	Duration   time.Duration
}

// NewResponse creates a new empty Response
func NewResponse() Response {
	return Response{
		StatusCode: 0,
		Headers:    make(map[string]string),
		Body:       "",
		Duration:   0,
	}
}

// Execute executes a curl command and returns the response
func Execute(curlCmd string) (Response, error) {
	resp := NewResponse()

	if curlCmd == "" {
		return resp, fmt.Errorf("empty curl command")
	}

	// Parse the curl command to extract parts
	// We need to append flags to the command
	// Remove "curl" from beginning if present
	cmd := strings.TrimSpace(curlCmd)
	if strings.HasPrefix(cmd, "curl ") {
		cmd = strings.TrimSpace(cmd[5:])
	}

	// Build the full command with our flags
	// -i: include headers
	// -s: silent mode
	// -w: custom output format for timing
	start := time.Now()

	fullCmd := fmt.Sprintf("curl -i -s -w '\\n___TIME___:%%{time_total}' %s", cmd)
	shellCmd := exec.Command("sh", "-c", fullCmd)

	output, err := shellCmd.CombinedOutput()
	duration := time.Since(start)

	if err != nil {
		// Try to get useful error info from output
		if len(output) > 0 {
			return resp, fmt.Errorf("curl failed: %s", string(output))
		}
		return resp, fmt.Errorf("curl failed: %v", err)
	}

	// Parse the output
	outputStr := string(output)

	// Extract timing info
	timeRegex := regexp.MustCompile(`___TIME___:(\d+\.?\d*)`)
	if matches := timeRegex.FindStringSubmatch(outputStr); len(matches) > 1 {
		if seconds, err := strconv.ParseFloat(matches[1], 64); err == nil {
			resp.Duration = time.Duration(seconds * float64(time.Second))
		}
	}

	// Remove timing info from output
	outputStr = timeRegex.ReplaceAllString(outputStr, "")

	// Split headers and body
	parts := strings.SplitN(outputStr, "\r\n\r\n", 2)
	if len(parts) < 2 {
		// Try \n\n separator
		parts = strings.SplitN(outputStr, "\n\n", 2)
	}

	if len(parts) >= 1 {
		// Parse headers
		headerLines := strings.Split(parts[0], "\n")

		// First line should be status line
		if len(headerLines) > 0 {
			statusLine := headerLines[0]
			// Extract status code from "HTTP/1.1 200 OK"
			statusRegex := regexp.MustCompile(`HTTP/[\d.]+\s+(\d+)`)
			if matches := statusRegex.FindStringSubmatch(statusLine); len(matches) > 1 {
				if code, err := strconv.Atoi(matches[1]); err == nil {
					resp.StatusCode = code
				}
			}
		}

		// Parse remaining headers
		for i := 1; i < len(headerLines); i++ {
			line := strings.TrimSpace(headerLines[i])
			if line == "" {
				continue
			}
			// Split on first colon
			colonIdx := strings.Index(line, ":")
			if colonIdx > 0 {
				key := strings.TrimSpace(line[:colonIdx])
				value := strings.TrimSpace(line[colonIdx+1:])
				resp.Headers[key] = value
			}
		}
	}

	if len(parts) >= 2 {
		resp.Body = strings.TrimSpace(parts[1])
	}

	if resp.Duration == 0 {
		resp.Duration = duration
	}

	return resp, nil
}
