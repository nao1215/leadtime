// Package config get setting from environment variable or configuration file.
package config

import (
	"github.com/caarlos0/env/v7"
	"github.com/nao1215/leadtime/domain/model"
)

// GitHubConfig represents configuration for GitHub.
type GitHubConfig struct {
	// AccessToken is access token for GitHub API.
	AccessToken model.Token `env:"LT_GITHUB_ACCESS_TOKEN,required"`
}

// NewGitHubConfig initialize github config.
// If user does not set environment variable LT_GITHUB_ACCESS_TOKEN,
// return error.
func NewGitHubConfig() (*GitHubConfig, error) {
	cfg := &GitHubConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, ErrNotSetGitHubAccessToken
	}

	return cfg, nil
}

// NewGitHubAccessToken return github access token
func NewGitHubAccessToken(config *GitHubConfig) model.Token {
	return config.AccessToken
}
