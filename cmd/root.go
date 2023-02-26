package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "leadtime",
		Short:   "leadtime statistics on the time it took for PRs to be merged",
		Long:    "leadtime statistics on the time it took for PRs to be merged",
		Example: "  LT_GITHUB_ACCESS_TOKEN=XXX leadtime stat --owner=nao1215 --repo=sqly",
	}
}

// Execute run leadtime process.
func Execute() int {
	rootCmd := newRootCmd()
	rootCmd.AddCommand(newStatCmd())
	rootCmd.AddCommand(newVersionCmd())

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)

		return 1
	}

	return 0
}
