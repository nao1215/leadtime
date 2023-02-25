package usecase

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/nao1215/leadtime/domain/service"
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

type LeadTime struct {
	PRstats []*PRStat
}

func (lt *LeadTime) Min() int {
	if len(lt.PRstats) == 0 {
		return 0
	}

	min := lt.PRstats[0].MergeTimeMinutes
	for _, v := range lt.PRstats[1:] {
		if v.MergeTimeMinutes < min {
			min = v.MergeTimeMinutes
		}
	}

	return min
}

func (lt *LeadTime) Max() int {
	if len(lt.PRstats) == 0 {
		return 0
	}

	max := lt.PRstats[0].MergeTimeMinutes
	for _, v := range lt.PRstats[1:] {
		if v.MergeTimeMinutes > max {
			max = v.MergeTimeMinutes
		}
	}

	return max
}

func (lt *LeadTime) Ave() float64 {
	if len(lt.PRstats) == 0 {
		return 0
	}

	sum := float64(0)
	for _, v := range lt.PRstats {
		sum += float64(v.MergeTimeMinutes)
	}

	return sum / float64(len(lt.PRstats))
}

func (lt *LeadTime) Sum() int {
	if len(lt.PRstats) == 0 {
		return 0
	}

	sum := 0
	for _, v := range lt.PRstats {
		sum += v.MergeTimeMinutes
	}

	return sum
}

func (lt *LeadTime) Median() float64 {
	if len(lt.PRstats) == 0 {
		return 0
	}

	nums := make([]int, 0, len(lt.PRstats))
	for _, v := range lt.PRstats {
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

type PRStat struct {
	Number           int
	Title            string
	MergeTimeMinutes int
}

// LTUsecase implement LeadTimeUsecase
type LTUsecase struct {
	prService     *service.PullRequestService
	commitService *service.CommitService
}

// NewLeadTimeUsecase initialize LTUsecase
func NewLeadTimeUsecase(prService *service.PullRequestService, commitService *service.CommitService) LeadTimeUsecase {
	return &LTUsecase{
		prService:     prService,
		commitService: commitService,
	}
}

// Stat return lead time statistics
func (lt *LTUsecase) Stat(ctx context.Context, input *LeadTimeUsecaseStatInput) (*LeadTimeUsecaseStatOutput, error) {
	prs, err := lt.prService.List(ctx, input.Owner, input.Repository)
	if err != nil {
		return nil, err
	}

	pullReqs := make([]*PullRequest, 0)
	for _, v := range prs {
		if !v.IsClosed() {
			continue
		}
		pr := &PullRequest{}
		pullReqs = append(pullReqs, pr.toUsecasePullRequest(v))
	}

	prStats := make([]*PRStat, 0)
	for _, v := range pullReqs {
		prStat := &PRStat{}
		commit, err := lt.commitService.GetFirstCommit(ctx, input.Owner, input.Repository, v.Number)
		if err != nil {
			if errors.Is(err, service.ErrNoCommit) {
				continue
			}
			return nil, err
		}

		prStat.Number = v.Number
		prStat.Title = v.Title
		if v.MergedAt != (time.Time{}) {
			prStat.MergeTimeMinutes = MinuteDiff(v.MergedAt, commit.Date.Time)
		} else if v.ClosedAt != (time.Time{}) {
			prStat.MergeTimeMinutes = MinuteDiff(v.ClosedAt, commit.Date.Time)
		} else {
			continue
		}

		prStats = append(prStats, prStat)
	}

	return &LeadTimeUsecaseStatOutput{
		LeadTime: &LeadTime{
			PRstats: prStats,
		},
	}, nil
}

func MinuteDiff(after, before time.Time) int {
	diff := after.Sub(before)
	return int(diff.Minutes())
}
