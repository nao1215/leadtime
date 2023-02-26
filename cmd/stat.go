package cmd

import (
	"context"
	"fmt"

	"github.com/nao1215/leadtime/di"
	"github.com/nao1215/leadtime/domain/usecase"
	"github.com/shogo82148/pointer"
	"github.com/spf13/cobra"
)

func newStatCmd() *cobra.Command {
	statCmd := &cobra.Command{
		Use:   "stat",
		Short: "Print GitHub pull request leadtime statics",
		Long: `Print GitHub pull request leadtime statics.

leadtime calculates statistics for PRs already in Closed/Merged status.`,
		Example: "  LT_GITHUB_ACCESS_TOKEN=XXX leadtime stat --owner=nao1215 --repo=sqly",
		RunE:    stat,
	}

	statCmd.Flags().StringP("owner", "o", "", "Specify GitHub owner name")
	statCmd.Flags().StringP("repo", "r", "", "Specify GitHub repository name")
	statCmd.Flags().BoolP("markdown", "m", false, "Output markdown")

	return statCmd
}

func stat(cmd *cobra.Command, args []string) error { //nolint
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

	markdown, err := cmd.Flags().GetBool("markdown")
	if err != nil {
		return err
	}

	input := &usecase.LeadTimeUsecaseStatInput{
		Owner:      owner,
		Repository: repo,
	}
	if err := input.Valid(); err != nil {
		return err
	}

	output, err := leadTime.LeadTimeUsecase.Stat(context.Background(), input)
	if err != nil {
		return err
	}
	output.LeadTime.RemoveOpenPR()

	if markdown {
		outputMarkdown(output.LeadTime)
		return nil
	}
	outputDefault(output.LeadTime)

	return nil
}

func outputMarkdown(lt *usecase.LeadTime) {
	fmt.Println("# Pull Request Lead Time")
	fmt.Println("## Statistics")
	fmt.Printf("Statistics were calculated for %d closed PRs.  \n", len(lt.PRs))
	fmt.Println("| Item | Result |")
	fmt.Println("|:-----|:-------|")
	fmt.Printf("| Lead Time(Max)|%d[min]|\n", lt.Max())
	fmt.Printf("| Lead Time(Min)|%d[min]|\n", lt.Min())
	fmt.Printf("| Lead Time(Sum)|%d[min]|\n", lt.Sum())
	fmt.Printf("| Lead Time(Ave)|%.2f[min]|\n", lt.Average())
	fmt.Printf("| Lead Time(MN )|%.2f[min]|\n", lt.Median())
	fmt.Println()
	fmt.Println("## Pull Request Detail")
	fmt.Println("| Number | Author | Bot | LeadTime[min] | Title |")
	fmt.Println("|:-------|:-------|:----|:--------------|:------|")
	for _, v := range lt.PRs {
		if v.User.Bot {
			fmt.Printf("|#%d|%s|%s|%d|%s|\n", v.Number, pointer.StringValue(v.User.Name), "yes", v.MergeTimeMinutes, v.Title)
			continue
		}
		fmt.Printf("|#%d|%s|%s|%d|%s|\n", v.Number, pointer.StringValue(v.User.Name), "no", v.MergeTimeMinutes, v.Title)
	}
}

func outputDefault(lt *usecase.LeadTime) {
	fmt.Printf("PR\tAuthor\tBot\tLeadTime[min]\tTitle\n")
	for _, v := range lt.PRs {
		if v.User.Bot {
			fmt.Printf("#%d\t%s\t%s\t%d\t%s\n", v.Number, pointer.StringValue(v.User.Name), "yes", v.MergeTimeMinutes, v.Title)
			continue
		}
		fmt.Printf("#%d\t%s\t%s\t%d\t%s\n", v.Number, pointer.StringValue(v.User.Name), "no", v.MergeTimeMinutes, v.Title)
	}

	fmt.Println("")
	fmt.Println("[statistics]")
	fmt.Printf(" Total PR       = %d\n", len(lt.PRs))
	fmt.Printf(" Lead Time(Max) = %d[min]\n", lt.Max())
	fmt.Printf(" Lead Time(Min) = %d[min]\n", lt.Min())
	fmt.Printf(" Lead Time(Sum) = %d[min]\n", lt.Sum())
	fmt.Printf(" Lead Time(Ave) = %.2f[min]\n", lt.Average())
	fmt.Printf(" Lead Time(Median) = %.2f[min]\n", lt.Median())
}
