package models

// SwaggerInfo contains additional documentation for the API
// This file should not contain any model definitions, only documentation and helpers

// Example responses for documentation purposes
var exampleResponses = struct {
	// Example successful authentication response
	AuthSuccess string
	// Example error response
	Error string
}{
	AuthSuccess: `{
		"token": "eyJhbGciOiJIUzI1NiIs..."
	}`,
	Error: `{
		"error": "Invalid request format"
	}`,
}

// Additional documentation for complex types or endpoints can be added here
