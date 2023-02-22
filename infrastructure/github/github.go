// Package github is http client for GitHub API.
// This package is wrapper for google/go-github package.
package github

import (
	"context"

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
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, &APIError{StatusCode: resp.StatusCode, Message: "failed to gey repository list"}
	}

	repoList := make([]*model.Repository, 0)
	for _, v := range repos {
		repo := &model.Repository{
			ID:          v.ID,
			Owner:       &model.User{Name: v.Owner.Name},
			Name:        v.Name,
			FullName:    v.FullName,
			Description: v.Description,
		}
		repoList = append(repoList, repo)
	}

	return repoList, nil
}
