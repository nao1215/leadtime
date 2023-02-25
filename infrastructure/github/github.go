// Package github is http client for GitHub API.
// This package is wrapper for google/go-github package.
package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v50/github"
	"github.com/nao1215/leadtime/domain/model"
	"github.com/nao1215/leadtime/domain/repository"
	"golang.org/x/oauth2"
)

// Client is http client for GitHub API.
type Client struct {
	*github.Client
}

// NewClient return http client for GitHub API.
func NewClient(token model.Token) *Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.String()},
	)
	client := oauth2.NewClient(context.Background(), tokenSource)

	return &Client{Client: github.NewClient(client)}
}

// GitHubRepository is http client for GitHub API
type GitHubRepository struct {
	client *Client
}

// NewGitHubRepository initialize repository.GitHubRepository
func NewGitHubRepository(client *Client) repository.GitHubRepository {
	return &GitHubRepository{client: client}
}

// ListRepositories return List the repositories for a user.
func (c *GitHubRepository) ListRepositories(ctx context.Context) ([]*model.Repository, error) {
	repos, resp, err := c.client.Repositories.List(ctx, "", nil)
	if resp != nil {
		defer func() error {
			if err := resp.Body.Close(); err != nil {
				return fmt.Errorf("failed to close response body: %w", err)
			}

			return nil
		}()
	}
	if err != nil {
		return nil, &APIError{StatusCode: resp.StatusCode, Message: "failed to get repository list"}
	}

	repoList := make([]*model.Repository, 0)
	for _, v := range repos {
		var user *model.User
		if v.Owner != nil {
			user = &model.User{
				Name: v.Owner.Name,
			}
		}

		repo := &model.Repository{
			ID:          v.ID,
			Owner:       user,
			Name:        v.Name,
			FullName:    v.FullName,
			Description: v.Description,
		}
		repoList = append(repoList, repo)
	}

	return repoList, nil
}

// ListPullRequests return List the pull requests.
func (c *GitHubRepository) ListPullRequests(ctx context.Context, owner, repo string) ([]*model.PullRequest, error) {

	pullReqs := make([]*model.PullRequest, 0)
	opts := &github.PullRequestListOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 20},
	}

	for {
		prs, resp, err := c.client.PullRequests.List(ctx, owner, repo, opts)
		if resp != nil {
			defer func() error {
				if err := resp.Body.Close(); err != nil {
					return fmt.Errorf("failed to close response body: %w", err)
				}

				return nil
			}()
		}
		if err != nil {
			return nil, &APIError{StatusCode: resp.StatusCode, Message: "failed to get pull request list"}
		}

		for _, v := range prs {
			pullReqs = append(pullReqs, toDomainModelPR(v))
		}

		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}

	return pullReqs, nil
}

// ListCommitsInPR return List the commits in the PR.
// oreder is newest to oldest.
func (c *GitHubRepository) ListCommitsInPR(ctx context.Context, owner, repo string, number int) ([]*model.Commit, error) {
	opts := &github.ListOptions{PerPage: 20}

	commitsInPR := make([]*model.Commit, 0)
	for {
		commits, resp, err := c.client.PullRequests.ListCommits(ctx, owner, repo, number, opts)
		if resp != nil {
			defer func() error {
				if err := resp.Body.Close(); err != nil {
					return fmt.Errorf("failed to close response body: %w", err)
				}

				return nil
			}()
		}
		if err != nil {
			return nil, &APIError{StatusCode: resp.StatusCode, Message: "failed to get git commit list"}
		}

		for _, v := range commits {
			commitsInPR = append(commitsInPR, toDomainModelCommit(v))
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return commitsInPR, nil
}

// toDomainModelPR convert *github.PullRequest to *model.PullRequest
func toDomainModelPR(githubPR *github.PullRequest) *model.PullRequest {
	var createdAt *model.Timestamp
	if githubPR.ClosedAt != nil {
		createdAt = &model.Timestamp{
			Time: githubPR.ClosedAt.Time,
		}
	}

	var closedAt *model.Timestamp
	if githubPR.ClosedAt != nil {
		closedAt = &model.Timestamp{
			Time: githubPR.ClosedAt.Time,
		}
	}

	var mergedAt *model.Timestamp
	if githubPR.MergedAt != nil {
		mergedAt = &model.Timestamp{
			Time: githubPR.MergedAt.Time,
		}
	}

	var user *model.User
	if githubPR.User != nil {
		user = &model.User{
			Name: githubPR.User.Name,
		}
	}

	pr := &model.PullRequest{
		ID:           githubPR.ID,
		Number:       githubPR.Number,
		State:        githubPR.State,
		Title:        githubPR.Title,
		CreatedAt:    createdAt,
		ClosedAt:     closedAt,
		MergedAt:     mergedAt,
		User:         user,
		Comments:     githubPR.Comments,
		Additions:    githubPR.Additions,
		Deletions:    githubPR.Deletions,
		ChangedFiles: githubPR.ChangedFiles,
	}

	return pr
}

// toDomainModelCommit convert *github.RepositoryCommit to *model.Commit.
func toDomainModelCommit(commit *github.RepositoryCommit) *model.Commit {
	var author *model.User
	if commit.Author != nil {
		author = &model.User{
			Name: commit.Author.Name,
		}
	}

	var committer *model.User
	if commit.Committer != nil {
		committer = &model.User{
			Name: commit.Committer.Name,
		}
	}

	var date *model.Timestamp
	if commit.Commit != nil && commit.Commit.Committer != nil {
		date = &model.Timestamp{
			Time: commit.Commit.Committer.GetDate().Time,
		}
	}

	domainModelCommit := &model.Commit{
		Author:    author,
		Committer: committer,
		Date:      date,
	}

	return domainModelCommit
}
