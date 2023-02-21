// Package config get setting from environment variable or configuration file.
package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewArgument(t *testing.T) {
	t.Parallel()
	t.Run("new argument", func(t *testing.T) {
		t.Parallel()

		want := &Argument{
			GitHubOwner:      "nao1215",
			GitHubRepository: "leadtime",
		}

		got := NewArgument(want.GitHubOwner, want.GitHubRepository)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})
}
