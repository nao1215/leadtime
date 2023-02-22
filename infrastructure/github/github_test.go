package github

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-github/v50/github"
	"github.com/nao1215/leadtime/domain/model"

	"github.com/google/go-cmp/cmp"
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
			w.Write(respBody)
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
			w.Write([]byte("error message"))
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
			t.Fatal("expect error occured, however got nil")
		}

		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("mismatch expect=%T, got=%T", &APIError{}, apiErr)
		}

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
			w.Write(respBody)
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
			t.Fatal("expect error occured, however got nil")
		}

		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("mismatch expect=%T, got=%T", &APIError{}, apiErr)
		}

		if apiErr.StatusCode != http.StatusUnauthorized {
			t.Errorf("mismatch expect=%d, got=%d", apiErr.StatusCode, http.StatusUnauthorized)
		}
	})
}
