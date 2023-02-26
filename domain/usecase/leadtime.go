package usecase

import (
	"context"
	"errors"
	"sort"
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
	Number           int
	State            string
	Title            string
	FirstCommitAt    time.Time
	CreatedAt        time.Time
	ClosedAt         time.Time
	MergedAt         time.Time
	UserName         string
	MergeTimeMinutes int
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
		p.UserName = pointer.StringValue(domainModelPR.User.Name)
	}

	if p.MergedAt != (time.Time{}) {
		p.MergeTimeMinutes = MinuteDiff(p.MergedAt, p.FirstCommitAt)
	} else if p.ClosedAt != (time.Time{}) {
		p.MergeTimeMinutes = MinuteDiff(p.ClosedAt, p.FirstCommitAt)
	}

	return p
}

type LeadTime struct {
	PRs []*PullRequest
}

func (lt *LeadTime) Min() int {
	if len(lt.PRs) == 0 {
		return 0
	}

	min := lt.PRs[0].MergeTimeMinutes
	for _, v := range lt.PRs[1:] {
		if v.MergeTimeMinutes < min {
			min = v.MergeTimeMinutes
		}
	}

	return min
}

func (lt *LeadTime) Max() int {
	if len(lt.PRs) == 0 {
		return 0
	}

	max := lt.PRs[0].MergeTimeMinutes
	for _, v := range lt.PRs[1:] {
		if v.MergeTimeMinutes > max {
			max = v.MergeTimeMinutes
		}
	}

	return max
}

func (lt *LeadTime) Ave() float64 {
	if len(lt.PRs) == 0 {
		return 0
	}

	sum := float64(0)
	for _, v := range lt.PRs {
		sum += float64(v.MergeTimeMinutes)
	}

	return sum / float64(len(lt.PRs))
}

func (lt *LeadTime) Sum() int {
	if len(lt.PRs) == 0 {
		return 0
	}

	sum := 0
	for _, v := range lt.PRs {
		sum += v.MergeTimeMinutes
	}

	return sum
}

func (lt *LeadTime) Median() float64 {
	if len(lt.PRs) == 0 {
		return 0
	}

	nums := make([]int, 0, len(lt.PRs))
	for _, v := range lt.PRs {
		nums = append(nums, v.MergeTimeMinutes)
	}
	sort.Ints(nums)

	var median float64
	mid := len(nums) / 2
	if len(nums)%2 == 0 {
		median = float64(nums[mid-1]+nums[mid]) / 2
	} else {
		median = float64(nums[mid])
	}

	return median
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
			PRs: pullReqs,
		},
	}, nil
}

func MinuteDiff(after, before time.Time) int {
	diff := after.Sub(before)
	return int(diff.Minutes())
}
