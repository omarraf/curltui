package curl

import "strings"

// BuildCurlCommand converts a Request struct to a valid curl command string
// TODO: Implement conversion of Request struct to curl command string
//   - Start with "curl" base command
//   - Add -X flag for HTTP method (if not GET)
//   - Add -H flags for each header in the format "Key: Value"
//   - Add -d, --data, or --data-raw flag for request body
//   - Properly escape/quote arguments for shell safety
//   - Handle special characters in URL, headers, and body
//   - Example output: curl -X POST -H "Content-Type: application/json" -d '{"key":"value"}' https://api.example.com
func BuildCurlCommand(req Request) string {
	var parts []string
	parts = append(parts, "curl")

	// Basic stub implementation - returns minimal curl command
	if req.URL != "" {
		parts = append(parts, req.URL)
	}

	return strings.Join(parts, " ")
}
