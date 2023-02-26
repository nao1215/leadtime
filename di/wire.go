//go:build wireinject
// +build wireinject

// Package di Inject dependence by wire command.
package di

import (
	"github.com/google/wire"
	"github.com/nao1215/leadtime/config"
	"github.com/nao1215/leadtime/domain/usecase"
	"github.com/nao1215/leadtime/infrastructure/github"
)

//go:generate wire

// LeadTime is usecase set.
type LeadTime struct {
	GithubConfig    *config.GitHubConfig
	LeadTimeUsecase usecase.LeadTimeUsecase
}

// newLeadTime initialize LeadTime struct
func newLeadTime(githubConfig *config.GitHubConfig, leadTimeUsecase usecase.LeadTimeUsecase) *LeadTime {
	return &LeadTime{
		GithubConfig:    githubConfig,
		LeadTimeUsecase: leadTimeUsecase,
	}
}

func NewLeadTime() (*LeadTime, error) {
	wire.Build(
		config.NewGitHubConfig,
		config.NewGitHubAccessToken,
		usecase.NewLeadTimeUsecase,
		github.NewClient,
		github.NewGitHubRepository,
		newLeadTime,
	)
	return &LeadTime{}, nil
}
