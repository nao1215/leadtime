package github

import "fmt"

// APIError is error for GitHub API.
type APIError struct {
	// StatusCode is HTTP status code from GitHub.
	StatusCode int
	// Message is error message
	Message string
}

// Error return string that represents a GitHub API error
func (e *APIError) Error() string {
	return fmt.Sprintf("GitHub API error: status code %d, message: %s", e.StatusCode, e.Message)
}
