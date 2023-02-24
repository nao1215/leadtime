// Package config get setting from environment variable or configuration file.
package config

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/leadtime/domain/model"
)

func TestNewGitHubConfig(t *testing.T) {
	const token = model.Token("test_token")

	t.Run("Get github config", func(t *testing.T) {
		t.Setenv("LT_GITHUB_ACCESS_TOKEN", token.String())

		want := &GitHubConfig{
			AccessToken: token,
		}
		got, err := NewGitHubConfig()
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("if user does not set github access token", func(t *testing.T) {
		_, got := NewGitHubConfig()
		if !errors.Is(got, ErrNotSetGitHubAccessToken) {
			t.Errorf("mismatch want=%v, got=%v", ErrNotSetGitHubAccessToken, got)
		}
	})
}
