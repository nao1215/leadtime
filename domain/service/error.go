package service

import "errors"

var (
	// ErrNoPullRequest means "there is no pull request in this repository"
	ErrNoPullRequest = errors.New("there is no pull request in this repository")
)
