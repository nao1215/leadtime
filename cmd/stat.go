package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"os"
	"sort"

	"github.com/nao1215/leadtime/di"
	"github.com/nao1215/leadtime/domain/usecase"
	"github.com/shogo82148/pointer"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func newStatCmd() *cobra.Command {
	statCmd := &cobra.Command{
		Use:   "stat",
		Short: "Print GitHub pull request leadtime statics",
		Long: `Print GitHub pull request leadtime statics.
leadtime calculates statistics for PRs already in Closed/Merged status.
|------------- lead time -------------|
|               |--- time to merge ---|
---------------------------------------
^               ^                     ^
first commit    create PR          merge PR
`,
		Example: "  LT_GITHUB_ACCESS_TOKEN=XXX leadtime stat --owner=nao1215 --repo=sqly",
		RunE:    stat,
	}

	statCmd.Flags().StringP("owner", "o", "", "Specify GitHub owner name")
	statCmd.Flags().StringP("repo", "r", "", "Specify GitHub repository name")
	statCmd.Flags().BoolP("markdown", "m", false, "Output markdown")
	statCmd.Flags().BoolP("exclude-bot", "B", false, "Exclude Pull Requests created by bots")
	statCmd.Flags().IntSliceP("exclude-pr", "P", []int{}, "Exclude specified Pull Requests (e.g. '-P 1,3,19')")
	statCmd.Flags().StringSliceP("exclude-user", "U", []string{}, "Exclude Pull Requests created by specified user (e.g. '-U nao,alice')")
	statCmd.Flags().BoolP("all", "a", false, "Print all data used for statistics")
	statCmd.Flags().BoolP("json", "j", false, "Output json")

	return statCmd
}

type option struct {
	// all is flag whether output statistical data instead of statistical information or not
	all bool
	// excludeBot is whether PRs created by bots exclude or not
	excludeBot bool
	// excludePRs is PR number list for exclusion
	excludePRs []int
	// excludeUsers is user list for exclusion
	excludeUsers []string
	// gitHubOwner is owner name
	gitHubOwner string
	// gitHubRepo is github repository
	gitHubRepo string
	// json is json output mode flag
	json bool
	// markdown is markdown output mode flag
	markdown bool
}

func (o *option) valid() error {
	if o.json && o.markdown {
		return ErrMultipleOutputFlag
	}
	return nil
}

func newOption(cmd *cobra.Command) (*option, error) {
	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		return nil, err
	}

	bot, err := cmd.Flags().GetBool("exclude-bot")
	if err != nil {
		return nil, err
	}

	excludePRs, err := cmd.Flags().GetIntSlice("exclude-pr")
	if err != nil {
		return nil, err
	}

	excludeUsers, err := cmd.Flags().GetStringSlice("exclude-user")
	if err != nil {
		return nil, err
	}

	owner, err := cmd.Flags().GetString("owner")
	if err != nil {
		return nil, err
	}

	repo, err := cmd.Flags().GetString("repo")
	if err != nil {
		return nil, err
	}

	json, err := cmd.Flags().GetBool("json")
	if err != nil {
		return nil, err
	}

	markdown, err := cmd.Flags().GetBool("markdown")
	if err != nil {
		return nil, err
	}

	return &option{
		all:          all,
		excludeBot:   bot,
		excludePRs:   excludePRs,
		excludeUsers: excludeUsers,
		gitHubOwner:  owner,
		gitHubRepo:   repo,
		markdown:     markdown,
		json:         json,
	}, nil
}

func stat(cmd *cobra.Command, args []string) error { //nolint
	leadTime, err := di.NewLeadTime()
	if err != nil {
		return err
	}

	opt, err := newOption(cmd)
	if err != nil {
		return err
	}

	if err := opt.valid(); err != nil {
		return err
	}

	input := &usecase.LeadTimeUsecaseStatInput{
		Owner:      opt.gitHubOwner,
		Repository: opt.gitHubRepo,
	}
	if err := input.Valid(); err != nil {
		return err
	}

	output, err := leadTime.LeadTimeUsecase.Stat(context.Background(), input)
	if err != nil {
		return err
	}

	dlts := newDetailLeadTimeStat(output.LeadTime)
	dlts.removePRs(opt)
	dlts.stat()

	return dlts.print(opt)
}

func (dlts *DetailLeadTimeStat) print(opt *option) error {
	if opt.markdown {
		if err := dlts.drawGraph(); err != nil {
			return err
		}
		dlts.markdown(opt.all)
		return nil
	}

	if opt.json {
		if err := dlts.json(os.Stdout, opt.all); err != nil {
			return err
		}
		return nil
	}

	dlts.stdout(opt.all)

	return nil
}

func (dlts *DetailLeadTimeStat) drawGraph() error {
	p := plot.New()
	p.X.Label.Text = "PR number"
	p.Y.Label.Text = "Lead Time[min]"

	data := make(plotter.XYs, 0, len(dlts.PullRequests))
	for _, v := range dlts.PullRequests {
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
	p.Y.Max = float64(dlts.max()) + 100
	line.Color = color.RGBA{R: 226, G: 45, B: 60, A: 255}
	line.Width = vg.Points(1.5)
	p.Add(line)

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "leadtime.png"); err != nil {
		return err
	}
	return nil
}

func (dlts *DetailLeadTimeStat) markdown(all bool) {
	fmt.Println("# Pull Request Lead Time")
	fmt.Println("## Statistics")
	fmt.Printf("Statistics were calculated for %d closed PRs.  \n", len(dlts.PullRequests))
	fmt.Println("| Item | Result |")
	fmt.Println("|:-----|:-------|")
	fmt.Printf("| Lead Time(Max)|%d[min]|\n", dlts.max())
	fmt.Printf("| Lead Time(Min)|%d[min]|\n", dlts.min())
	fmt.Printf("| Lead Time(Sum)|%d[min]|\n", dlts.sum())
	fmt.Printf("| Lead Time(Ave)|%.2f[min]|\n", dlts.average())
	fmt.Printf("| Lead Time(MN )|%.2f[min]|\n", dlts.median())
	fmt.Println()
	fmt.Println("![PR Lead Time](./leadtime.png)")
	fmt.Println()

	if all {
		fmt.Println("## Pull Request Detail")
		fmt.Println("| Number | Author | Bot | LeadTime[min] | Title |")
		fmt.Println("|:-------|:-------|:----|:--------------|:------|")
		for _, v := range dlts.PullRequests {
			if v.User.Bot {
				fmt.Printf("|#%d|%s|%s|%d|%s|\n", v.Number, pointer.StringValue(v.User.Name), "yes", v.MergeTimeMinutes, v.Title)
				continue
			}
			fmt.Printf("|#%d|%s|%s|%d|%s|\n", v.Number, pointer.StringValue(v.User.Name), "no", v.MergeTimeMinutes, v.Title)
		}
	}
}

func (dlts *DetailLeadTimeStat) stdout(all bool) {
	if all {
		fmt.Printf("PR\tAuthor\tBot\tLeadTime[min]\tTitle\n")
		for _, v := range dlts.PullRequests {
			if v.User.Bot {
				fmt.Printf("#%d\t%s\t%s\t%d\t%s\n", v.Number, pointer.StringValue(v.User.Name), "yes", v.MergeTimeMinutes, v.Title)
				continue
			}
			fmt.Printf("#%d\t%s\t%s\t%d\t%s\n", v.Number, pointer.StringValue(v.User.Name), "no", v.MergeTimeMinutes, v.Title)
		}
		fmt.Println("")
	}
	fmt.Println("[statistics]")
	fmt.Printf(" Total PR       = %d\n", len(dlts.PullRequests))
	fmt.Printf(" Lead Time(Max) = %d[min]\n", dlts.max())
	fmt.Printf(" Lead Time(Min) = %d[min]\n", dlts.min())
	fmt.Printf(" Lead Time(Sum) = %d[min]\n", dlts.sum())
	fmt.Printf(" Lead Time(Ave) = %.2f[min]\n", dlts.average())
	fmt.Printf(" Lead Time(Median) = %.2f[min]\n", dlts.median())
}

// LeadTimeStat is Lead time statistics.
type LeadTimeStat struct {
	TotalPR           int     `json:"total_pr,omitempty"`
	LeadTimeMaximum   int     `json:"lead_time_maximum,omitempty"`
	LeadTimeMinimum   int     `json:"lead_time_minimum,omitempty"`
	LeadTimeSummation int     `json:"lead_time_summation,omitempty"`
	LeadTimeAverage   float64 `json:"lead_time_average,omitempty"`
	LeadTimeMedian    float64 `json:"lead_time_median,omitempty"`
}

type DetailLeadTimeStat struct {
	LeadTimeStatistics *LeadTimeStat          `json:"lead_time_statistics,omitempty"`
	PullRequests       []*usecase.PullRequest `json:"pull_requests,omitempty"`
}

func newDetailLeadTimeStat(lt *usecase.LeadTime) *DetailLeadTimeStat {
	return &DetailLeadTimeStat{
		LeadTimeStatistics: &LeadTimeStat{},
		PullRequests:       lt.PullRequests,
	}
}

func (dlts *DetailLeadTimeStat) json(w io.Writer, all bool) error {
	var bytes []byte
	var err error
	if all {
		bytes, err = json.Marshal(dlts)
	} else {
		bytes, err = json.Marshal(dlts.LeadTimeStatistics)
	}
	if err != nil {
		return err
	}

	fmt.Fprintln(w, string(bytes))

	return nil
}

func (dlts *DetailLeadTimeStat) stat() {
	dlts.LeadTimeStatistics = &LeadTimeStat{
		TotalPR:           len(dlts.PullRequests),
		LeadTimeMaximum:   dlts.max(),
		LeadTimeMinimum:   dlts.min(),
		LeadTimeSummation: dlts.sum(),
		LeadTimeAverage:   dlts.average(),
		LeadTimeMedian:    dlts.median(),
	}
}

func (dlts *DetailLeadTimeStat) removePRs(opt *option) {
	dlts.removeOpenPR()
	if opt.excludeBot {
		dlts.removePRCreatedByBot()
	}
	if len(opt.excludePRs) != 0 {
		dlts.removeSpecifiedPRs(opt.excludePRs)
	}
	if len(opt.excludeUsers) != 0 {
		dlts.removePRsCreatedByTargetUser(opt.excludeUsers)
	}
}

func (dlts *DetailLeadTimeStat) removeOpenPR() {
	prs := make([]*usecase.PullRequest, 0, len(dlts.PullRequests))
	for _, v := range dlts.PullRequests {
		if v.State == "open" {
			continue
		}
		prs = append(prs, v)
	}
	dlts.PullRequests = prs
}

func (dlts *DetailLeadTimeStat) removePRCreatedByBot() {
	prs := make([]*usecase.PullRequest, 0, len(dlts.PullRequests))
	for _, v := range dlts.PullRequests {
		if v.User.IsBot() {
			continue
		}
		prs = append(prs, v)
	}
	dlts.PullRequests = prs
}

func (dlts *DetailLeadTimeStat) removeSpecifiedPRs(removeTargetPRs []int) {
	prs := make([]*usecase.PullRequest, 0, len(dlts.PullRequests))
	for _, v := range dlts.PullRequests {
		if slices.Contains(removeTargetPRs, v.Number) {
			continue
		}
		prs = append(prs, v)
	}
	dlts.PullRequests = prs
}

func (dlts *DetailLeadTimeStat) removePRsCreatedByTargetUser(target []string) {
	prs := make([]*usecase.PullRequest, 0, len(dlts.PullRequests))
	for _, v := range dlts.PullRequests {
		if slices.Contains(target, pointer.StringValue(v.User.Name)) {
			continue
		}
		prs = append(prs, v)
	}
	dlts.PullRequests = prs
}

func (dlts *DetailLeadTimeStat) min() int {
	if len(dlts.PullRequests) == 0 {
		return 0
	}

	min := dlts.PullRequests[0].MergeTimeMinutes
	for _, v := range dlts.PullRequests[1:] {
		if v.MergeTimeMinutes < min {
			min = v.MergeTimeMinutes
		}
	}

	return min
}

func (dlts *DetailLeadTimeStat) max() int {
	if len(dlts.PullRequests) == 0 {
		return 0
	}

	max := dlts.PullRequests[0].MergeTimeMinutes
	for _, v := range dlts.PullRequests[1:] {
		if v.MergeTimeMinutes > max {
			max = v.MergeTimeMinutes
		}
	}

	return max
}

func (dlts *DetailLeadTimeStat) average() float64 {
	if len(dlts.PullRequests) == 0 {
		return 0
	}

	sum := float64(0)
	for _, v := range dlts.PullRequests {
		sum += float64(v.MergeTimeMinutes)
	}

	return sum / float64(len(dlts.PullRequests))
}

func (dlts *DetailLeadTimeStat) sum() int {
	if len(dlts.PullRequests) == 0 {
		return 0
	}

	sum := 0
	for _, v := range dlts.PullRequests {
		sum += v.MergeTimeMinutes
	}

	return sum
}

func (dlts *DetailLeadTimeStat) median() float64 {
	if len(dlts.PullRequests) == 0 {
		return 0
	}

	nums := make([]int, 0, len(dlts.PullRequests))
	for _, v := range dlts.PullRequests {
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
