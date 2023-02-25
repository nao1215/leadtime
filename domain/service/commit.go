package service

import (
	"context"

	"github.com/nao1215/leadtime/domain/model"
	"github.com/nao1215/leadtime/domain/repository"
)

// CommitService is service that manipulate commit in the pr.
type CommitService struct {
	githubRepository repository.GitHubRepository
}

// NewCommitRequestService initialize CommitRequestService
func NewCommitRequestService(gihub repository.GitHubRepository) *CommitService {
	return &CommitService{
		githubRepository: gihub,
	}
}

// GetFirstCommit return first commit in the pull request.
func (c *CommitService) GetFirstCommit(ctx context.Context, owner, repository string, number int) (*model.Commit, error) {
	list, err := c.githubRepository.ListCommitsInPR(ctx, owner, repository, number)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, ErrNoCommit
	}

	return list[0], nil
}
