// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/nao1215/leadtime/config"
	"github.com/nao1215/leadtime/domain/service"
	"github.com/nao1215/leadtime/domain/usecase"
	"github.com/nao1215/leadtime/infrastructure/github"
)

// Injectors from wire.go:

func NewLeadTime() (*LeadTime, error) {
	gitHubConfig, err := config.NewGitHubConfig()
	if err != nil {
		return nil, err
	}
	token := config.NewGitHubAccessToken(gitHubConfig)
	client := github.NewClient(token)
	gitHubRepository := github.NewGitHubRepository(client)
	pullRequestService := service.NewPullRequestService(gitHubRepository)
	pullRequestUsecase := usecase.NewPullRequestUsecase(pullRequestService)
	commitService := service.NewCommitRequestService(gitHubRepository)
	leadTimeUsecase := usecase.NewLeadTimeUsecase(pullRequestService, commitService)
	leadTime := newLeadTime(gitHubConfig, pullRequestUsecase, leadTimeUsecase)
	return leadTime, nil
}

// wire.go:

// LeadTime is usecase set.
type LeadTime struct {
	GithubConfig       *config.GitHubConfig
	PullRequestUsecase usecase.PullRequestUsecase
	LeadTimeUsecase    usecase.LeadTimeUsecase
}

// newLeadTime initialize LeadTime struct
func newLeadTime(githubConfig *config.GitHubConfig, pullRequestUsecase usecase.PullRequestUsecase,
	leadTimeUsecase usecase.LeadTimeUsecase) *LeadTime {
	return &LeadTime{
		GithubConfig:       githubConfig,
		PullRequestUsecase: pullRequestUsecase,
		LeadTimeUsecase:    leadTimeUsecase,
	}
}
