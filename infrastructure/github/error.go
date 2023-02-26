package github

import (
	"errors"
	"fmt"
)

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

var (
	// ErrNoPullRequest means "there is no pull request in this repository"
	ErrNoPullRequest = errors.New("there is no pull request in this repository")
	// ErrNoCommit means "there is no commit in this repository"
	ErrNoCommit = errors.New("there is no commit in this repository")
)
