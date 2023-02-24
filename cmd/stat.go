package cmd

import "github.com/spf13/cobra"

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "Print GitHub pull request leadtime statics",
	Long:  `Print GitHub pull request leadtime statics`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return stat(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(statCmd)
}

func stat(cmd *cobra.Command, args []string) error {
	return nil
}
