package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/nao1215/leadtime/domain/model"
	"github.com/nao1215/leadtime/domain/repository"
	"github.com/nao1215/leadtime/infrastructure/github"
	"github.com/shogo82148/pointer"
)

// LeadTimeUsecase is use cases for stat leadtime
type LeadTimeUsecase interface {
	Stat(ctx context.Context, input *LeadTimeUsecaseStatInput) (*LeadTimeUsecaseStatOutput, error)
}

// LeadTimeUsecaseStatInput is input data for LeadTimeUsecase.Stat().
type LeadTimeUsecaseStatInput struct {
	// Owner is GitHub account name
	Owner string
	// Repository is GitHub repository name
	Repository string
}

// Valid is input data validation
func (lt *LeadTimeUsecaseStatInput) Valid() error {
	if lt.Owner == "" {
		return ErrEmptyGitHubOwnerName
	}
	if lt.Repository == "" {
		return ErrEmptyRepositoryName
	}
	return nil
}

// LeadTimeUsecaseStatOutput is output data for LeadTimeUsecase.Stat().
type LeadTimeUsecaseStatOutput struct {
	LeadTime *LeadTime
}

// LTUsecase implement LeadTimeUsecase
type LTUsecase struct {
	gitHubRepo repository.GitHubRepository
}

// NewLeadTimeUsecase initialize LTUsecase
func NewLeadTimeUsecase(gitHubRepo repository.GitHubRepository) LeadTimeUsecase {
	return &LTUsecase{
		gitHubRepo: gitHubRepo,
	}
}

// PullRequest is PR information for presentation layer.
type PullRequest struct {
	Number           int         `json:"number,omitempty"`
	State            string      `json:"state,omitempty"`
	Title            string      `json:"title,omitempty"`
	FirstCommitAt    time.Time   `json:"first_commit_at,omitempty"`
	CreatedAt        time.Time   `json:"created_at,omitempty"`
	ClosedAt         time.Time   `json:"closed_at,omitempty"`
	MergedAt         time.Time   `json:"merged_at,omitempty"`
	User             *model.User `json:"user,omitempty"`
	MergeTimeMinutes int         `json:"merge_time_minutes,omitempty"`
}

func (p *PullRequest) toUsecasePullRequest(domainModelPR *model.PullRequest, firstCommitAt time.Time) *PullRequest {
	p.Number = pointer.IntValue(domainModelPR.Number)
	p.Title = pointer.StringValue(domainModelPR.Title)
	p.State = pointer.StringValue(domainModelPR.State)
	p.FirstCommitAt = firstCommitAt

	if domainModelPR.CreatedAt != nil {
		p.CreatedAt = pointer.TimeValue(&domainModelPR.CreatedAt.Time)
	}
	if domainModelPR.ClosedAt != nil {
		p.ClosedAt = pointer.TimeValue(&domainModelPR.ClosedAt.Time)
	}
	if domainModelPR.MergedAt != nil {
		p.MergedAt = pointer.TimeValue(&domainModelPR.MergedAt.Time)
	}
	if domainModelPR.User != nil {
		p.User = domainModelPR.User
	}

	if p.MergedAt != (time.Time{}) {
		p.MergeTimeMinutes = MinuteDiff(p.MergedAt, p.FirstCommitAt)
	} else if p.ClosedAt != (time.Time{}) {
		p.MergeTimeMinutes = MinuteDiff(p.ClosedAt, p.FirstCommitAt)
	} else {
		p.MergeTimeMinutes = MinuteDiff(time.Now(), p.FirstCommitAt)
	}

	return p
}

type LeadTime struct {
	PullRequests []*PullRequest `json:"pull_requests,omitempty"`
}

// Stat return lead time statistics
func (lt *LTUsecase) Stat(ctx context.Context, input *LeadTimeUsecaseStatInput) (*LeadTimeUsecaseStatOutput, error) {
	prs, err := lt.gitHubRepo.ListPullRequests(ctx, input.Owner, input.Repository)
	if err != nil {
		return nil, err
	}

	pullReqs := make([]*PullRequest, 0)
	for _, v := range prs {
		if v.Number == nil {
			continue
		}

		commit, err := lt.gitHubRepo.GetFirstCommit(ctx, input.Owner, input.Repository, *v.Number)
		if err != nil {
			if errors.Is(err, github.ErrNoCommit) {
				continue
			}
			return nil, err
		}

		pr := &PullRequest{}
		pullReqs = append(pullReqs, pr.toUsecasePullRequest(v, commit.Date.Time))
	}

	return &LeadTimeUsecaseStatOutput{
		LeadTime: &LeadTime{
			PullRequests: pullReqs,
		},
	}, nil
}

func MinuteDiff(after, before time.Time) int {
	diff := after.Sub(before)
	return int(diff.Minutes())
}
