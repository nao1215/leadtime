package repository

import (
	"context"

	"github.com/nao1215/leadtime/domain/model"
)

// GitHubRepository is interface for manipulating GitHub.
type GitHubRepository interface {
	// ListRepositories return repository list
	ListRepositories(ctx context.Context) ([]*model.Repository, error)
	// ListRepositories return pull request list
	ListPullRequests(ctx context.Context, owner, repo string) ([]*model.PullRequest, error)
	// ListCommitsInPR return commits in PR.
	ListCommitsInPR(ctx context.Context, owner, repo string, number int) ([]*model.Commit, error)
	// GetFirstCommit return first commit in PR.
	GetFirstCommit(ctx context.Context, owner, repository string, number int) (*model.Commit, error)
}
