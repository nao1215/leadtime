package usecase

import "errors"

var (
	// ErrEmptyGitHubAccessToken means "github access token is empty"
	ErrEmptyGitHubAccessToken = errors.New("github access token is empty")
	// ErrEmptyGitHubOwnerName means "github owner name is empty"
	ErrEmptyGitHubOwnerName = errors.New("github owner name is empty")
)
