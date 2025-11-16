package tui

import (
	"github.com/username/curltui/internal/curl"
	"github.com/username/curltui/internal/http"
)

// NewModel creates and initializes a new Model
func NewModel(importCmd string) Model {
	m := Model{
		method:          "GET",
		url:             "",
		headers:         []HeaderPair{},
		body:            "",
		focusedField:    URLField, // Start on URL field
		cursorPos:       0,
		headerIndex:     0,
		editing:         true, // Start in editing mode so user can type immediately
		loading:         false,
		pasteMode:       false,
		pasteBuffer:     "",
		responseTab:     0,
		bodyScroll:      0,
		responseScroll:  0,
		headerEditField: 0,
		response:        http.NewResponse(),
		err:             nil,
		width:           0,
		height:          0,
	}

	// If importCmd is provided, parse it and populate fields
	if importCmd != "" {
		req, err := curl.ParseCurlCommand(importCmd)
		if err != nil {
			m.err = err
		} else {
			m.method = req.Method
			m.url = req.URL
			m.body = req.Body

			// Convert headers map to slice
			for key, value := range req.Headers {
				m.headers = append(m.headers, HeaderPair{Key: key, Value: value})
			}
		}
	}

	return m
}
