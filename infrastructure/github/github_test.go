package github

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v50/github"
	"github.com/nao1215/leadtime/domain/model"
)

func TestListRepositories(t *testing.T) {
	t.Parallel()

	t.Run("Get repository list", func(t *testing.T) {
		t.Parallel()

		token := model.Token("good_token")
		client := NewClient(token)
		repo := NewGitHubRepository(client)
		ctx := context.Background()

		// Test server
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respBody, err := json.Marshal([]github.Repository{
				{
					ID:          github.Int64(1),
					Owner:       &github.User{Name: github.String("user1")},
					Name:        github.String("repo1"),
					FullName:    github.String("user1/repo1"),
					Description: github.String("repo1 description"),
				},
				{
					ID:          github.Int64(2),
					Owner:       &github.User{Name: github.String("user2")},
					Name:        github.String("repo2"),
					FullName:    github.String("user2/repo2"),
					Description: github.String("repo2 description"),
				},
			})
			if err != nil {
				t.Fatal(err)
			}
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(respBody); err != nil {
				t.Fatal(err)
			}
		}))
		defer testServer.Close()

		testURL, err := url.Parse(testServer.URL)
		if err != nil {
			t.Fatal(err)
		}
		client.BaseURL = testURL
		if !strings.HasSuffix(client.BaseURL.Path, "/") {
			client.BaseURL.Path += "/"
		}

		want := []*model.Repository{
			{
				ID:          github.Int64(1),
				Owner:       &model.User{Name: github.String("user1")},
				Name:        github.String("repo1"),
				FullName:    github.String("user1/repo1"),
				Description: github.String("repo1 description"),
			},
			{
				ID:          github.Int64(2),
				Owner:       &model.User{Name: github.String("user2")},
				Name:        github.String("repo2"),
				FullName:    github.String("user2/repo2"),
				Description: github.String("repo2 description"),
			},
		}
		// test start
		got, err := repo.ListRepositories(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Return status code 500 from GitHub", func(t *testing.T) {
		t.Parallel()

		token := model.Token("test_token")
		client := NewClient(token)
		repo := NewGitHubRepository(client)
		ctx := context.Background()

		// Test server
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte("error message")); err != nil {
				t.Fatal(err)
			}
		}))
		defer testServer.Close()

		testURL, err := url.Parse(testServer.URL)
		if err != nil {
			t.Fatal(err)
		}
		client.BaseURL = testURL
		if !strings.HasSuffix(client.BaseURL.Path, "/") {
			client.BaseURL.Path += "/"
		}

		// test start
		_, err = repo.ListRepositories(ctx)
		if err == nil {
			t.Fatal("expect error occurred, however got nil")
		}

		var apiError *APIError
		if !errors.As(err, &apiError) {
			t.Fatalf("mismatch expect=%T, got=%T", &APIError{}, err)
		}

		apiErr := err.(*APIError) //nolint
		if apiErr.StatusCode != http.StatusInternalServerError {
			t.Errorf("mismatch expect=%d, got=%d", apiErr.StatusCode, http.StatusInternalServerError)
		}
	})

	t.Run("Return status code 401 from GitHub", func(t *testing.T) {
		t.Parallel()

		token := model.Token("bad_token")
		client := NewClient(token)
		repo := NewGitHubRepository(client)
		ctx := context.Background()

		// Test server
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respBody, err := json.Marshal([]github.Repository{
				{
					ID:          github.Int64(1),
					Owner:       &github.User{Name: github.String("user1")},
					Name:        github.String("repo1"),
					FullName:    github.String("user1/repo1"),
					Description: github.String("repo1 description"),
				},
			})
			if err != nil {
				t.Fatal(err)
			}
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write(respBody); err != nil {
				t.Fatal(err)
			}
		}))
		defer testServer.Close()

		testURL, err := url.Parse(testServer.URL)
		if err != nil {
			t.Fatal(err)
		}
		client.BaseURL = testURL
		if !strings.HasSuffix(client.BaseURL.Path, "/") {
			client.BaseURL.Path += "/"
		}

		// test start
		_, err = repo.ListRepositories(ctx)
		if err == nil {
			t.Fatal("expect error occurred, however got nil")
		}

		var apiError *APIError
		if !errors.As(err, &apiError) {
			t.Fatalf("mismatch expect=%T, got=%T", &APIError{}, err)
		}

		apiErr := err.(*APIError) //nolint
		if apiErr.StatusCode != http.StatusUnauthorized {
			t.Errorf("mismatch expect=%d, got=%d", apiErr.StatusCode, http.StatusUnauthorized)
		}
	})
}

func TestClient_ListPullRequests(t *testing.T) {
	t.Parallel()

	const apiURL = "/repos/owner/repo/pulls"

	t.Run("Get PR list", func(t *testing.T) {
		t.Parallel()

		now := time.Date(2023, 2, 24, 12, 34, 56, 0, time.UTC)
		testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			wantURL := apiURL
			if wantURL != req.URL.Path {
				t.Errorf("mismatch want=%v, got=%s", wantURL, req.URL.Path)
			}

			wantHTTPMethod := http.MethodGet
			if wantHTTPMethod != req.Method {
				t.Errorf("mismatch want=%v, got=%s", wantHTTPMethod, req.Method)
			}

			respBody, err := json.Marshal([]github.PullRequest{
				{
					ID:     github.Int64(1),
					Number: github.Int(1),
					State:  github.String("open"),
					Title:  github.String("test_pr1"),
					CreatedAt: &github.Timestamp{
						Time: now,
					},
					ClosedAt: &github.Timestamp{
						Time: now,
					},
					MergedAt: &github.Timestamp{
						Time: now,
					},
					User: &github.User{
						Login: github.String("test_user1"),
					},
					Comments:     github.Int(0),
					Additions:    github.Int(10),
					Deletions:    github.Int(5),
					ChangedFiles: github.Int(2),
				},
				{
					ID:     github.Int64(2),
					Number: github.Int(2),
					State:  github.String("closed"),
					Title:  github.String("test_pr2"),
					CreatedAt: &github.Timestamp{
						Time: now,
					},
					ClosedAt: &github.Timestamp{
						Time: now,
					},
					MergedAt: &github.Timestamp{
						Time: now,
					},
					User: &github.User{
						Login: github.String("test_user2"),
					},
					Comments:     github.Int(2),
					Additions:    github.Int(5),
					Deletions:    github.Int(2),
					ChangedFiles: github.Int(1),
				},
			})
			if err != nil {
				t.Fatal(err)
			}
			rw.WriteHeader(http.StatusOK)
			if _, err := rw.Write(respBody); err != nil {
				t.Fatal(err)
			}
		}))
		defer testServer.Close()

		token := model.Token("token")
		client := NewClient(token)
		repo := NewGitHubRepository(client)
		ctx := context.Background()

		testURL, err := url.Parse(testServer.URL)
		if err != nil {
			t.Fatal(err)
		}
		client.BaseURL = testURL
		if !strings.HasSuffix(client.BaseURL.Path, "/") {
			client.BaseURL.Path += "/"
		}

		wantPRs := []*model.PullRequest{
			{
				ID:           github.Int64(1),
				Number:       github.Int(1),
				State:        github.String("open"),
				Title:        github.String("test_pr1"),
				CreatedAt:    &model.Timestamp{Time: now},
				ClosedAt:     &model.Timestamp{Time: now},
				MergedAt:     &model.Timestamp{Time: now},
				User:         &model.User{Name: github.String("test_user1")},
				Comments:     github.Int(0),
				Additions:    github.Int(10),
				Deletions:    github.Int(5),
				ChangedFiles: github.Int(2),
			},
			{
				ID:           github.Int64(2),
				Number:       github.Int(2),
				State:        github.String("closed"),
				Title:        github.String("test_pr2"),
				CreatedAt:    &model.Timestamp{Time: now},
				ClosedAt:     &model.Timestamp{Time: now},
				MergedAt:     &model.Timestamp{Time: now},
				User:         &model.User{Name: github.String("test_user2")},
				Comments:     github.Int(2),
				Additions:    github.Int(5),
				Deletions:    github.Int(2),
				ChangedFiles: github.Int(1),
			},
		}
		gotPRs, err := repo.ListPullRequests(ctx, "owner", "repo")
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(wantPRs, gotPRs); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Get empty PR list", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			wantURL := apiURL
			if wantURL != req.URL.Path {
				t.Errorf("mismatch want=%v, got=%s", wantURL, req.URL.Path)
			}

			wantHTTPMethod := http.MethodGet
			if wantHTTPMethod != req.Method {
				t.Errorf("mismatch want=%v, got=%s", wantHTTPMethod, req.Method)
			}

			respBody, err := json.Marshal([]github.PullRequest{})
			if err != nil {
				t.Fatal(err)
			}
			rw.WriteHeader(http.StatusOK)
			if _, err := rw.Write(respBody); err != nil {
				t.Fatal(err)
			}
		}))
		defer testServer.Close()

		token := model.Token("token")
		client := NewClient(token)
		repo := NewGitHubRepository(client)
		ctx := context.Background()

		testURL, err := url.Parse(testServer.URL)
		if err != nil {
			t.Fatal(err)
		}
		client.BaseURL = testURL
		if !strings.HasSuffix(client.BaseURL.Path, "/") {
			client.BaseURL.Path += "/"
		}

		want := ErrNoPullRequest
		_, got := repo.ListPullRequests(ctx, "owner", "repo")
		if !errors.Is(got, want) {
			t.Errorf("mismatch want=%v, got=%v", want, got)
		}
	})

	t.Run("Return status code 500 from GitHub", func(t *testing.T) {
		t.Parallel()

		token := model.Token("test_token")
		client := NewClient(token)
		repo := NewGitHubRepository(client)
		ctx := context.Background()

		testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			wantURL := apiURL
			if wantURL != req.URL.Path {
				t.Errorf("mismatch want=%v, got=%s", wantURL, req.URL.Path)
			}

			wantHTTPMethod := http.MethodGet
			if wantHTTPMethod != req.Method {
				t.Errorf("mismatch want=%v, got=%s", wantHTTPMethod, req.Method)
			}

			rw.WriteHeader(http.StatusInternalServerError)
			if _, err := rw.Write([]byte("error message")); err != nil {
				t.Fatal(err)
			}
		}))
		defer testServer.Close()

		testURL, err := url.Parse(testServer.URL)
		if err != nil {
			t.Fatal(err)
		}
		client.BaseURL = testURL
		if !strings.HasSuffix(client.BaseURL.Path, "/") {
			client.BaseURL.Path += "/"
		}

		// test start
		_, err = repo.ListPullRequests(ctx, "owner", "repo")
		if err == nil {
			t.Fatal("expect error occurred, however got nil")
		}

		var apiError *APIError
		if !errors.As(err, &apiError) {
			t.Fatalf("mismatch expect=%T, got=%T", &APIError{}, err)
		}

		apiErr := err.(*APIError) //nolint
		if apiErr.StatusCode != http.StatusInternalServerError {
			t.Errorf("mismatch expect=%d, got=%d", apiErr.StatusCode, http.StatusInternalServerError)
		}
	})

	t.Run("Return status code 401 from GitHub", func(t *testing.T) {
		t.Parallel()

		token := model.Token("bad_token")
		client := NewClient(token)
		repo := NewGitHubRepository(client)
		ctx := context.Background()

		testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			wantURL := apiURL
			if wantURL != req.URL.Path {
				t.Errorf("mismatch want=%v, got=%s", wantURL, req.URL.Path)
			}

			wantHTTPMethod := http.MethodGet
			if wantHTTPMethod != req.Method {
				t.Errorf("mismatch want=%v, got=%s", wantHTTPMethod, req.Method)
			}
			respBody, err := json.Marshal([]github.PullRequest{})
			if err != nil {
				t.Fatal(err)
			}
			rw.WriteHeader(http.StatusUnauthorized)
			if _, err := rw.Write(respBody); err != nil {
				t.Fatal(err)
			}
		}))
		defer testServer.Close()

		testURL, err := url.Parse(testServer.URL)
		if err != nil {
			t.Fatal(err)
		}
		client.BaseURL = testURL
		if !strings.HasSuffix(client.BaseURL.Path, "/") {
			client.BaseURL.Path += "/"
		}

		// test start
		_, err = repo.ListPullRequests(ctx, "owner", "repo")
		if err == nil {
			t.Fatal("expect error occurred, however got nil")
		}

		var apiError *APIError
		if !errors.As(err, &apiError) {
			t.Fatalf("mismatch expect=%T, got=%T", &APIError{}, err)
		}

		apiErr := err.(*APIError) //nolint
		if apiErr.StatusCode != http.StatusUnauthorized {
			t.Errorf("mismatch expect=%d, got=%d", apiErr.StatusCode, http.StatusUnauthorized)
		}
	})
}

func TestClient_ListCommitsInPR(t *testing.T) {
	t.Parallel()

	const (
		apiURL = "/repos/owner/repo/pulls/123/commits"
		token  = "token"
	)
	now := time.Date(2023, 2, 24, 12, 34, 56, 0, time.UTC)

	t.Run("Get all commit in the PR", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			wantURL := apiURL
			if wantURL != req.URL.Path {
				t.Errorf("mismatch want=%v, got=%s", wantURL, req.URL.Path)
			}

			wantHTTPMethod := http.MethodGet
			if wantHTTPMethod != req.Method {
				t.Errorf("mismatch want=%v, got=%s", wantHTTPMethod, req.Method)
			}

			respBody, err := json.Marshal([]github.RepositoryCommit{
				{
					Commit: &github.Commit{
						Committer: &github.CommitAuthor{
							Date: &github.Timestamp{Time: now},
						},
					},
					Author: &github.User{
						Name: github.String("author1"),
					},
					Committer: &github.User{
						Name: github.String("comitter1"),
					},
				},
				{
					Commit: &github.Commit{
						Committer: &github.CommitAuthor{
							Date: &github.Timestamp{Time: now},
						},
					},
					Author: &github.User{
						Name: github.String("author2"),
					},
					Committer: &github.User{
						Name: github.String("comitter2"),
					},
				},
			})
			if err != nil {
				t.Fatal(err)
			}
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(respBody); err != nil {
				t.Fatal(err)
			}
		}))
		defer testServer.Close()

		client := NewClient(token)
		repo := NewGitHubRepository(client)
		ctx := context.Background()

		testURL, err := url.Parse(testServer.URL)
		if err != nil {
			t.Fatal(err)
		}
		client.BaseURL = testURL
		if !strings.HasSuffix(client.BaseURL.Path, "/") {
			client.BaseURL.Path += "/"
		}

		wantCommits := []*model.Commit{
			{
				Author:    &model.User{Name: github.String("author1")},
				Committer: &model.User{Name: github.String("comitter1")},
				Date:      &model.Timestamp{Time: now},
			},
			{
				Author:    &model.User{Name: github.String("author2")},
				Committer: &model.User{Name: github.String("comitter2")},
				Date:      &model.Timestamp{Time: now},
			},
		}
		gotCommits, err := repo.ListCommitsInPR(ctx, "owner", "repo", 123)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(wantCommits, gotCommits); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Return status code 500 from GitHub", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			wantURL := apiURL
			if wantURL != req.URL.Path {
				t.Errorf("mismatch want=%v, got=%s", wantURL, req.URL.Path)
			}

			wantHTTPMethod := http.MethodGet
			if wantHTTPMethod != req.Method {
				t.Errorf("mismatch want=%v, got=%s", wantHTTPMethod, req.Method)
			}
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte("error message")); err != nil {
				t.Fatal(err)
			}
		}))
		defer testServer.Close()

		client := NewClient(token)
		repo := NewGitHubRepository(client)
		ctx := context.Background()

		testURL, err := url.Parse(testServer.URL)
		if err != nil {
			t.Fatal(err)
		}
		client.BaseURL = testURL
		if !strings.HasSuffix(client.BaseURL.Path, "/") {
			client.BaseURL.Path += "/"
		}

		_, err = repo.ListCommitsInPR(ctx, "owner", "repo", 123)
		if err == nil {
			t.Fatal("expect error occurred, however got nil")
		}

		var apiError *APIError
		if !errors.As(err, &apiError) {
			t.Fatalf("mismatch expect=%T, got=%T", &APIError{}, err)
		}

		apiErr := err.(*APIError) //nolint
		if apiErr.StatusCode != http.StatusInternalServerError {
			t.Errorf("mismatch expect=%d, got=%d", apiErr.StatusCode, http.StatusInternalServerError)
		}
	})

	t.Run("Return status code 401 from GitHub", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			wantURL := apiURL
			if wantURL != req.URL.Path {
				t.Errorf("mismatch want=%v, got=%s", wantURL, req.URL.Path)
			}

			wantHTTPMethod := http.MethodGet
			if wantHTTPMethod != req.Method {
				t.Errorf("mismatch want=%v, got=%s", wantHTTPMethod, req.Method)
			}
			respBody, err := json.Marshal([]github.PullRequest{})
			if err != nil {
				t.Fatal(err)
			}
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write(respBody); err != nil {
				t.Fatal(err)
			}
		}))
		defer testServer.Close()

		client := NewClient(token)
		repo := NewGitHubRepository(client)
		ctx := context.Background()

		testURL, err := url.Parse(testServer.URL)
		if err != nil {
			t.Fatal(err)
		}
		client.BaseURL = testURL
		if !strings.HasSuffix(client.BaseURL.Path, "/") {
			client.BaseURL.Path += "/"
		}

		_, err = repo.ListCommitsInPR(ctx, "owner", "repo", 123)
		if err == nil {
			t.Fatal("expect error occurred, however got nil")
		}

		var apiError *APIError
		if !errors.As(err, &apiError) {
			t.Fatalf("mismatch expect=%T, got=%T", &APIError{}, err)
		}

		apiErr := err.(*APIError) //nolint
		if apiErr.StatusCode != http.StatusUnauthorized {
			t.Errorf("mismatch expect=%d, got=%d", apiErr.StatusCode, http.StatusUnauthorized)
		}
	})
}

func TestGitHubRepository_GetFirstCommit(t *testing.T) {
	t.Parallel()

	const (
		apiURL = "/repos/owner/repo/pulls/123/commits"
		token  = "token"
	)
	now := time.Date(2023, 2, 24, 12, 34, 56, 0, time.UTC)

	t.Run("Get first commit in the PR", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			wantURL := apiURL
			if wantURL != req.URL.Path {
				t.Errorf("mismatch want=%v, got=%s", wantURL, req.URL.Path)
			}

			wantHTTPMethod := http.MethodGet
			if wantHTTPMethod != req.Method {
				t.Errorf("mismatch want=%v, got=%s", wantHTTPMethod, req.Method)
			}

			respBody, err := json.Marshal([]github.RepositoryCommit{
				{
					Commit: &github.Commit{
						Committer: &github.CommitAuthor{
							Date: &github.Timestamp{Time: now},
						},
					},
					Author: &github.User{
						Name: github.String("author1"),
					},
					Committer: &github.User{
						Name: github.String("comitter1"),
					},
				},
				{
					Commit: &github.Commit{
						Committer: &github.CommitAuthor{
							Date: &github.Timestamp{Time: now},
						},
					},
					Author: &github.User{
						Name: github.String("author2"),
					},
					Committer: &github.User{
						Name: github.String("comitter2"),
					},
				},
			})
			if err != nil {
				t.Fatal(err)
			}
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(respBody); err != nil {
				t.Fatal(err)
			}
		}))
		defer testServer.Close()

		client := NewClient(token)
		repo := NewGitHubRepository(client)
		ctx := context.Background()

		testURL, err := url.Parse(testServer.URL)
		if err != nil {
			t.Fatal(err)
		}
		client.BaseURL = testURL
		if !strings.HasSuffix(client.BaseURL.Path, "/") {
			client.BaseURL.Path += "/"
		}

		want := &model.Commit{
			Author:    &model.User{Name: github.String("author1")},
			Committer: &model.User{Name: github.String("comitter1")},
			Date:      &model.Timestamp{Time: now},
		}
		got, err := repo.GetFirstCommit(ctx, "owner", "repo", 123)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("No commit in the PR", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			wantURL := apiURL
			if wantURL != req.URL.Path {
				t.Errorf("mismatch want=%v, got=%s", wantURL, req.URL.Path)
			}

			wantHTTPMethod := http.MethodGet
			if wantHTTPMethod != req.Method {
				t.Errorf("mismatch want=%v, got=%s", wantHTTPMethod, req.Method)
			}

			respBody, err := json.Marshal([]github.RepositoryCommit{})
			if err != nil {
				t.Fatal(err)
			}
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(respBody); err != nil {
				t.Fatal(err)
			}
		}))
		defer testServer.Close()

		client := NewClient(token)
		repo := NewGitHubRepository(client)
		ctx := context.Background()

		testURL, err := url.Parse(testServer.URL)
		if err != nil {
			t.Fatal(err)
		}
		client.BaseURL = testURL
		if !strings.HasSuffix(client.BaseURL.Path, "/") {
			client.BaseURL.Path += "/"
		}

		want := ErrNoCommit
		_, got := repo.GetFirstCommit(ctx, "owner", "repo", 123)
		if !errors.Is(got, want) {
			t.Errorf("mismatch want=%v, got=%v", want, got)
		}
	})
}

func Test_toDomainModelPR(t *testing.T) {
	t.Parallel()

	now := time.Date(2023, 2, 24, 12, 34, 56, 0, time.UTC)

	tests := []struct {
		name     string
		githubPR *github.PullRequest
		want     *model.PullRequest
	}{
		{
			name:     "convert empty PR",
			githubPR: &github.PullRequest{},
			want:     &model.PullRequest{},
		},
		{
			name: "convert PR without user/time information",
			githubPR: &github.PullRequest{
				ID:     github.Int64(1),
				Number: github.Int(1),
				State:  github.String("open"),
				Title:  github.String("test_pr1"),
				// Without CreatedAt, ClosedAt, MergedAt, User
				Comments:     github.Int(0),
				Additions:    github.Int(10),
				Deletions:    github.Int(5),
				ChangedFiles: github.Int(2),
			},
			want: &model.PullRequest{
				ID:           github.Int64(1),
				Number:       github.Int(1),
				State:        github.String("open"),
				Title:        github.String("test_pr1"),
				Comments:     github.Int(0),
				Additions:    github.Int(10),
				Deletions:    github.Int(5),
				ChangedFiles: github.Int(2),
			},
		},
		{
			name: "convert PR",
			githubPR: &github.PullRequest{
				ID:     github.Int64(1),
				Number: github.Int(1),
				State:  github.String("open"),
				Title:  github.String("test_pr1"),
				CreatedAt: &github.Timestamp{
					Time: now,
				},
				ClosedAt: &github.Timestamp{
					Time: now,
				},
				MergedAt: &github.Timestamp{
					Time: now,
				},
				User: &github.User{
					Login: github.String("test_user1"),
				},
				Comments:     github.Int(0),
				Additions:    github.Int(10),
				Deletions:    github.Int(5),
				ChangedFiles: github.Int(2),
			},
			want: &model.PullRequest{
				ID:           github.Int64(1),
				Number:       github.Int(1),
				State:        github.String("open"),
				Title:        github.String("test_pr1"),
				CreatedAt:    &model.Timestamp{Time: now},
				ClosedAt:     &model.Timestamp{Time: now},
				MergedAt:     &model.Timestamp{Time: now},
				User:         &model.User{Name: github.String("test_user1")},
				Comments:     github.Int(0),
				Additions:    github.Int(10),
				Deletions:    github.Int(5),
				ChangedFiles: github.Int(2),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := toDomainModelPR(tt.githubPR)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_toDomainModelCommit(t *testing.T) {
	t.Parallel()

	now := time.Date(2023, 2, 24, 12, 34, 56, 0, time.UTC)

	tests := []struct {
		name   string
		commit *github.RepositoryCommit
		want   *model.Commit
	}{
		{
			name:   "PR without commit",
			commit: &github.RepositoryCommit{},
			want:   &model.Commit{},
		},
		{
			name: "convert git commit to domain model commit",
			commit: &github.RepositoryCommit{
				Commit: &github.Commit{
					Committer: &github.CommitAuthor{
						Date: &github.Timestamp{Time: now},
					},
				},
				Author: &github.User{
					Name: github.String("author"),
				},
				Committer: &github.User{
					Name: github.String("comitter"),
				},
			},
			want: &model.Commit{
				Author:    &model.User{Name: github.String("author")},
				Committer: &model.User{Name: github.String("comitter")},
				Date:      &model.Timestamp{Time: now},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := toDomainModelCommit(tt.commit)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
