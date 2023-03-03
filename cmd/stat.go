package cmd

import (
	"context"
	"fmt"
	"image/color"

	"github.com/nao1215/leadtime/di"
	"github.com/nao1215/leadtime/domain/usecase"
	"github.com/shogo82148/pointer"
	"github.com/spf13/cobra"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
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
		if err := drawGraph(output.LeadTime); err != nil {
			return err
		}
		outputMarkdown(output.LeadTime)
		return nil
	}
	outputDefault(output.LeadTime)

	return nil
}

func drawGraph(lt *usecase.LeadTime) error {
	p := plot.New()
	p.X.Label.Text = "PR number"
	p.Y.Label.Text = "Lead Time[min]"

	data := plotter.XYs{}
	for _, v := range lt.PRs {
		data = append(data, plotter.XY{
			X: float64(v.Number),
			Y: float64(v.MergeTimeMinutes),
		})
	}

	line, err := plotter.NewLine(data)
	if err != nil {
		return err
	}
	p.Add(plotter.NewGrid())
	line.Color = color.RGBA{R: 226, G: 45, B: 60, A: 255}
	line.Width = vg.Points(1.5)
	p.Add(line)

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "leadtime.png"); err != nil {
		return err
	}
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
	fmt.Println("![PR Lead Time](./leadtime.png)")
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
