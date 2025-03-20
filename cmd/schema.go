package cmd

import (
	"github.com/ParetoSecurity/agent/claims"
	"github.com/ParetoSecurity/agent/runner"
	"github.com/spf13/cobra"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Output schema for all checks",
	Long:  "Output schema for all checks in JSON format.",
	Run: func(cc *cobra.Command, args []string) {
		runner.PrintSchemaJSON(claims.All)
	},
}

func init() {
	rootCmd.AddCommand(schemaCmd)
}
