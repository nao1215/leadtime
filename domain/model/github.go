// Package model is domain model and business logic.
package model

import "time"

// Token is token (e.g. github access token)
type Token string

func (t Token) String() string {
	return string(t)
}

// Repository represents GitHub repository information
type Repository struct {
	// ID is repository id
	ID *int64 `json:"id,omitempty"`
	// Owner is repository owner
	Owner *User `json:"owner,omitempty"`
	// Name is repository name
	Name *string `json:"name,omitempty"`
	// FullName is repository full name
	FullName *string `json:"full_name,omitempty"`
	// Description is repository description
	Description *string `json:"description,omitempty"`
}

// User represents a GitHub user.
type User struct {
	// Name is user name.
	Name *string `json:"name,omitempty"`
}

// Timestamp represents a time.
type Timestamp struct {
	Time time.Time
}

// PullRequest represents a GitHub pull request on a repository.
type PullRequest struct {
	// ID is PR's id.
	ID *int64
	// Number is PR number
	Number *int
	// State is PR state(e.g. closed)
	State *string
	// Title is PR title
	Title *string
	// CreatedAt is date of PR creation
	CreatedAt *Timestamp
	// ClosedAt is date of PR close
	ClosedAt *Timestamp
	// MergedAt is date of PR merged
	MergedAt *Timestamp
	// User is user information
	User *User
	// Comments is PR comment count
	Comments *int
	// Additions is number of addition lines
	Additions *int
	// Deletions is number of deletions line
	Deletions *int
	// ChangedFiles is number of changed files
	ChangedFiles *int
}

// IsClosed check whether pull request is closed or not.
func (pr *PullRequest) IsClosed() bool {
	if pr.State == nil {
		return false
	}

	return *pr.State == "closed"
}

// Commit is git commit information
type Commit struct {
	// Author is author user
	Author *User
	// Committer is commiter user
	Committer *User
	// Date is commit date
	Date *Timestamp
}
