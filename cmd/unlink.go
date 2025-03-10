package cmd

import (
	"os"
	"testing"

	"github.com/ParetoSecurity/agent/shared"
	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

var unlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Unlink this device from the team",
	Run: func(cc *cobra.Command, args []string) {
		log.Info("Unlinking device ...")
		shared.Config.TeamID = ""
		shared.Config.AuthToken = ""
		if err := shared.SaveConfig(); err != nil {
			log.WithError(err).Warn("failed to save config")
			if testing.Testing() {
				return
			}
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(unlinkCmd)
}
