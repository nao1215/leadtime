package config

import "errors"

var (
	// ErrNotSetGitHubAccessToken : for security concerns, set the environment variable
	// LT_GITHUB_ACCESS_TOKEN to the GitHub access token. The token should not set by command argument.
	ErrNotSetGitHubAccessToken = errors.New("GitHub access token is not set in the environment variable LT_GITHUB_ACCESS_TOKEN")
)
