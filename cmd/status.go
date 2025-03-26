package cmd

import (
	"os"

	"github.com/ParetoSecurity/agent/shared"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print the status of the checks",
	Run: func(cmd *cobra.Command, args []string) {
		shared.PrintStates()
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
