package curl

import (
	"strings"
	"testing"
)

func TestBuildCurlCommand(t *testing.T) {
	tests := []struct {
		name     string
		req      Request
		contains []string
	}{
		{
			name: "GET request",
			req: Request{
				Method: "GET",
				URL:    "https://api.example.com",
			},
			contains: []string{"curl", "-X GET", "https://api.example.com"},
		},
		{
			name: "POST request with body",
			req: Request{
				Method: "POST",
				URL:    "https://api.example.com",
				Body:   `{"key":"value"}`,
			},
			contains: []string{"curl", "-X POST", "-d", "https://api.example.com"},
		},
		{
			name: "Request with headers",
			req: Request{
				Method:  "GET",
				URL:     "https://api.example.com",
				Headers: map[string]string{"Content-Type": "application/json"},
			},
			contains: []string{"curl", "-X GET", "-H", "Content-Type: application/json"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildCurlCommand(tt.req)
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("BuildCurlCommand() = %v, should contain %v", result, substr)
				}
			}
		})
	}
}
