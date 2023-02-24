package usecase

import (
	"context"
	"fmt"

	"github.com/nao1215/leadtime/domain/model"
	"github.com/nao1215/leadtime/domain/service"
)

// PullRequestUsecase is use cases for obtaining PR information.
type PullRequestUsecase interface {
	List(ctx context.Context, input *PRUsecaseListInput) (*PRUsecaseListOutput, error)
}

// PRUsecaseListInput is input data for PRUsecase.List().
type PRUsecaseListInput struct {
	// Owner is GitHub account name
	Owner string
	// Repository is GitHub repository name
	Repository string
}

// Valid is input data validation
func (pr *PRUsecaseListInput) Valid() error {
	if pr.Owner == "" {
		return ErrEmptyGitHubOwnerName
	}
	// If repository name is empty, it is ok.
	return nil
}

// isEmptyRepositoryName check whether repository name is empty.
// true means empty, false means not empty.
func (pr *PRUsecaseListInput) isEmptyRepositoryName() bool {
	return pr.Repository == ""
}

// PRUsecaseListOutput is output data for PRUsecase.List().
type PRUsecaseListOutput struct {
}

// PullRequest is PR information for presentation layer.
type PullRequest struct {
	Number    int64
	State     string
	Title     string
	CreatedAt *model.Timestamp
	ClosedAt  *model.Timestamp
	MergedAt  *model.Timestamp
	User      model.User
}

// PRUsecase implement PullRequestUsecase
type PRUsecase struct {
	prService *service.PullRequestService
}

// NewPullRequestUsecase initialize pullRequestUsecase
func NewPullRequestUsecase(prService *service.PullRequestService) PullRequestUsecase {
	return &PRUsecase{
		prService: prService,
	}
}

// List return all pull request information in a repository.
func (pr *PRUsecase) List(ctx context.Context, input *PRUsecaseListInput) (*PRUsecaseListOutput, error) {
	prs, err := pr.prService.List(ctx, input.Owner, input.Repository)
	if err != nil {
		return nil, err
	}

	for _, v := range prs {
		/*
			fmt.Printf("ID:%d\n", v.ID)
			fmt.Printf("Number:%d\n", *v.Number)
			fmt.Printf("State:%s\n", *v.State)
			fmt.Printf("Title:%s\n", *v.Title)
		*/
		fmt.Println(v)
	}

	return nil, nil
}
