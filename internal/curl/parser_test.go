package curl

import (
	"testing"
)

func TestParseCurlCommand(t *testing.T) {
	tests := []struct {
		name    string
		cmd     string
		want    Request
		wantErr bool
	}{
		{
			name: "Simple GET",
			cmd:  "curl https://api.example.com",
			want: Request{
				Method:  "GET",
				URL:     "https://api.example.com",
				Headers: make(map[string]string),
				Body:    "",
			},
			wantErr: false,
		},
		{
			name: "POST with header and body",
			cmd:  `curl -X POST -H "Content-Type: application/json" -d '{"key":"value"}' https://api.example.com`,
			want: Request{
				Method:  "POST",
				URL:     "https://api.example.com",
				Headers: map[string]string{"Content-Type": "application/json"},
				Body:    `{"key":"value"}`,
			},
			wantErr: false,
		},
		{
			name:    "Empty command",
			cmd:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCurlCommand(tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCurlCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Method != tt.want.Method {
					t.Errorf("ParseCurlCommand() Method = %v, want %v", got.Method, tt.want.Method)
				}
				if got.URL != tt.want.URL {
					t.Errorf("ParseCurlCommand() URL = %v, want %v", got.URL, tt.want.URL)
				}
				if got.Body != tt.want.Body {
					t.Errorf("ParseCurlCommand() Body = %v, want %v", got.Body, tt.want.Body)
				}
			}
		})
	}
}
