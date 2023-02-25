package service

import "errors"

var (
	// ErrNoPullRequest means "there is no pull request in this repository"
	ErrNoPullRequest = errors.New("there is no pull request in this repository")
	// ErrNoCommit means "there is no commit in this repository"
	ErrNoCommit = errors.New("there is no commit in this repository")
)
