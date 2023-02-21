// Package config get setting from environment variable or configuration file.
package config

// Config represents configuration information for the leadtime CLI command.
// For example, include GitHub API access information.
type Config struct {
	// GitHubAccessToken is access token for GitHub API.
	GitHubAccessToken string
}

// Argument represents the argument for leadtime CLI startup.
type Argument struct {
	// GitHubOwner is GitHun owner name(e.g. nao1215)
	GitHubOwner string
	// GitHubRepository is GitHub repository name(e.g. leadtime)
	GitHubRepository string
}

// NewArgument initialize Argument struct.
func NewArgument(owner, repo string) *Argument {
	return &Argument{
		GitHubOwner:      owner,
		GitHubRepository: repo,
	}
}
