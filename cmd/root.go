package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "leadtime",
		Short: "leadtime statistics on the time it took for PRs to be merged",
		Long:  "leadtime statistics on the time it took for PRs to be merged",
	}
}

// Execute run leadtime process.
func Execute() int {
	rootCmd := newRootCmd()
	rootCmd.AddCommand(newStatCmd())

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)

		return 1
	}

	return 0
}
