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

		token := "good_token"
		client := NewClient(token)
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
		client.client.BaseURL = testURL
		if !strings.HasSuffix(client.client.BaseURL.Path, "/") {
			client.client.BaseURL.Path += "/"
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
		got, err := client.ListRepositories(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Return status code 500 from GitHub", func(t *testing.T) {
		t.Parallel()

		token := "test_token"
		client := NewClient(token)
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
		client.client.BaseURL = testURL
		if !strings.HasSuffix(client.client.BaseURL.Path, "/") {
			client.client.BaseURL.Path += "/"
		}

		// test start
		_, err = client.ListRepositories(ctx)
		if err == nil {
			t.Fatal("expect error occurred, however got nil")
		}

		var apiError *APIError
		if !errors.As(err, &apiError) {
			t.Fatalf("mismatch expect=%T, got=%T", &APIError{}, err)
		}

		apiErr := err.(*APIError)
		if apiErr.StatusCode != http.StatusInternalServerError {
			t.Errorf("mismatch expect=%d, got=%d", apiErr.StatusCode, http.StatusInternalServerError)
		}
	})

	t.Run("Return status code 401 from GitHub", func(t *testing.T) {
		t.Parallel()

		token := "bad_token"
		client := NewClient(token)
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
		client.client.BaseURL = testURL
		if !strings.HasSuffix(client.client.BaseURL.Path, "/") {
			client.client.BaseURL.Path += "/"
		}

		// test start
		_, err = client.ListRepositories(ctx)
		if err == nil {
			t.Fatal("expect error occurred, however got nil")
		}

		var apiError *APIError
		if !errors.As(err, &apiError) {
			t.Fatalf("mismatch expect=%T, got=%T", &APIError{}, err)
		}

		apiErr := err.(*APIError)
		if apiErr.StatusCode != http.StatusUnauthorized {
			t.Errorf("mismatch expect=%d, got=%d", apiErr.StatusCode, http.StatusUnauthorized)
		}
	})
}

func TestClient_ListPullRequests(t *testing.T) {
	t.Parallel()

	t.Run("Get PR list", func(t *testing.T) {
		t.Parallel()

		now := time.Date(2023, 2, 24, 12, 34, 56, 0, time.UTC)
		testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			wantURL := "/repos/owner/repo/pulls"
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
						Name: github.String("test_user1"),
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
						Name: github.String("test_user2"),
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

		token := "token"
		client := NewClient(token)
		ctx := context.Background()

		testURL, err := url.Parse(testServer.URL)
		if err != nil {
			t.Fatal(err)
		}
		client.client.BaseURL = testURL
		if !strings.HasSuffix(client.client.BaseURL.Path, "/") {
			client.client.BaseURL.Path += "/"
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
		gotPRs, err := client.ListPullRequests(ctx, "owner", "repo")
		if diff := cmp.Diff(wantPRs, gotPRs); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Return status code 500 from GitHub", func(t *testing.T) {
		t.Parallel()

		token := "test_token"
		client := NewClient(token)
		ctx := context.Background()

		testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			wantURL := "/repos/owner/repo/pulls"
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
		client.client.BaseURL = testURL
		if !strings.HasSuffix(client.client.BaseURL.Path, "/") {
			client.client.BaseURL.Path += "/"
		}

		// test start
		_, err = client.ListPullRequests(ctx, "owner", "repo")
		if err == nil {
			t.Fatal("expect error occurred, however got nil")
		}

		var apiError *APIError
		if !errors.As(err, &apiError) {
			t.Fatalf("mismatch expect=%T, got=%T", &APIError{}, err)
		}

		apiErr := err.(*APIError)
		if apiErr.StatusCode != http.StatusInternalServerError {
			t.Errorf("mismatch expect=%d, got=%d", apiErr.StatusCode, http.StatusInternalServerError)
		}
	})

	t.Run("Return status code 401 from GitHub", func(t *testing.T) {
		t.Parallel()

		token := "bad_token"
		client := NewClient(token)
		ctx := context.Background()

		testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			wantURL := "/repos/owner/repo/pulls"
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
		client.client.BaseURL = testURL
		if !strings.HasSuffix(client.client.BaseURL.Path, "/") {
			client.client.BaseURL.Path += "/"
		}

		// test start
		_, err = client.ListPullRequests(ctx, "owner", "repo")
		if err == nil {
			t.Fatal("expect error occurred, however got nil")
		}

		var apiError *APIError
		if !errors.As(err, &apiError) {
			t.Fatalf("mismatch expect=%T, got=%T", &APIError{}, err)
		}

		apiErr := err.(*APIError)
		if apiErr.StatusCode != http.StatusUnauthorized {
			t.Errorf("mismatch expect=%d, got=%d", apiErr.StatusCode, http.StatusUnauthorized)
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
					Name: github.String("test_user1"),
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
