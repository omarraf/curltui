package tui

import (
	"github.com/username/curltui/internal/curl"
	"github.com/username/curltui/internal/http"
)

// NewModel creates and initializes a new Model
// TODO: Implement initialization logic
//   - Create a new Request with default values (GET method, empty URL)
//   - If importCmd is provided, parse it using curl.ParseCurlCommand
//   - Handle parsing errors gracefully (store in model.err)
//   - Set focusedField to 0 (first field - method selector)
//   - Initialize empty Response
//   - Set initial dimensions (will be updated on first WindowSizeMsg)
func NewModel(importCmd string) Model {
	req := curl.NewRequest()
	resp := http.NewResponse()

	// TODO: If importCmd is not empty, parse it and populate req
	// TODO: Handle any parsing errors

	return Model{
		request:      req,
		response:     resp,
		focusedField: 0,
		err:          nil,
		width:        0,
		height:       0,
	}
}
