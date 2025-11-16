package curl

// Request represents an HTTP request
type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
}

// NewRequest creates a new Request with default values
func NewRequest() Request {
	return Request{
		Method:  "GET",
		URL:     "",
		Headers: make(map[string]string),
		Body:    "",
	}
}
