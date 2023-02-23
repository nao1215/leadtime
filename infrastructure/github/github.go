// Package github is http client for GitHub API.
// This package is wrapper for google/go-github package.
package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v50/github"
	"github.com/nao1215/leadtime/domain/model"
	"golang.org/x/oauth2"
)

// Client is http client for GitHub API.
type Client struct {
	// client is http client
	client *github.Client
}

// NewClient return http client for GitHub API.
func NewClient(token string) *Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	client := oauth2.NewClient(context.Background(), tokenSource)

	return &Client{client: github.NewClient(client)}
}

// ListRepositories return List the repositories for a user.
func (c *Client) ListRepositories(ctx context.Context) ([]*model.Repository, error) {
	repos, resp, err := c.client.Repositories.List(ctx, "", nil)
	if resp != nil {
		defer func() error {
			if err := resp.Body.Close(); err != nil {
				return fmt.Errorf("failed to close response body: %w", err)
			}

			return nil // nolint
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
func (c *Client) ListPullRequests(ctx context.Context, owner, repo string) ([]*model.PullRequest, error) {
	prs, resp, err := c.client.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{})
	if resp != nil {
		defer func() error {
			if err := resp.Body.Close(); err != nil {
				return fmt.Errorf("failed to close response body: %w", err)
			}

			return nil // nolint
		}()
	}
	if err != nil {
		return nil, &APIError{StatusCode: resp.StatusCode, Message: "failed to get pull request list"}
	}

	pullReqs := make([]*model.PullRequest, 0)
	for _, v := range prs {
		pullReqs = append(pullReqs, toDomainModelPR(v))
	}

	return pullReqs, nil
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
