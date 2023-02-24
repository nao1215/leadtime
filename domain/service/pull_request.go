package service

import (
	"context"

	"github.com/nao1215/leadtime/domain/model"
	"github.com/nao1215/leadtime/domain/repository"
)

// PullRequestService is service that github repository
type PullRequestService struct {
	githubRepository repository.GitHubRepository
}

// NewPullRequestService initialize PullRequestService
func NewPullRequestService(gihub repository.GitHubRepository) *PullRequestService {
	return &PullRequestService{
		githubRepository: gihub,
	}
}

// List return pull request list in repository
func (pr *PullRequestService) List(ctx context.Context, owner, repository string) ([]*model.PullRequest, error) {
	list, err := pr.githubRepository.ListPullRequests(ctx, owner, repository)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, ErrNoPullRequest
	}

	return list, nil
}
