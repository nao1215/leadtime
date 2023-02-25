package cmd

import (
	"context"

	"github.com/nao1215/leadtime/di"
	"github.com/nao1215/leadtime/domain/usecase"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{ //nolint
	Use:   "list",
	Short: "List up GitHub pull request information",
	Long: `List up GitHub repository or pull request information.
	
If you not specify repository name, leadtime command list up repository.
If you sepcify repository, leadtime command list up pull requests in the
repository.`,
	Example: "  leadtime list --owner=nao1215 --repo=leadtime",
	RunE:    list,
}

func init() { //nolint
	/*
		listCmd.Flags().StringP("owner", "o", "", "Specify GitHub owner name")
		listCmd.Flags().StringP("repo", "r", "", "Specify GitHub repository name")
		rootCmd.AddCommand(listCmd)
	*/
}

func list(cmd *cobra.Command, args []string) error { //nolint
	leadTime, err := di.NewLeadTime()
	if err != nil {
		return err
	}

	owner, err := cmd.Flags().GetString("owner")
	if err != nil {
		return err
	}

	repo, err := cmd.Flags().GetString("repo")
	if err != nil {
		return err
	}

	input := &usecase.PRUsecaseListInput{
		Owner:      owner,
		Repository: repo,
	}
	if err := input.Valid(); err != nil {
		return err
	}

	_, err = leadTime.PullRequestUsecase.List(context.Background(), input)
	if err != nil {
		return err
	}

	return nil
}
