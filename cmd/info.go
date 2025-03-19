package cmd

import (
	"encoding/json"
	"os"
	"runtime"

	"github.com/ParetoSecurity/agent/shared"
	"github.com/caarlos0/log"
	"github.com/elastic/go-sysinfo"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print the system and reports information",
	Run: func(cmd *cobra.Command, args []string) {

		log.Infof("%s@%s %s", shared.Version, shared.Commit, shared.Date)
		log.Infof("Built with %s", runtime.Version())

		device := shared.CurrentReportingDevice()
		jsonOutput, err := json.MarshalIndent(device, "", "  ")
		if err != nil {
			log.Warn("Failed to marshal host info")
		}
		log.Infof("Device Info: %s", string(jsonOutput))

		hostInfo, err := sysinfo.Host()
		if err != nil {
			log.Warn("Failed to get process information")
		}
		envInfo := hostInfo.Info()
		envInfo.IPs = []string{}  // Exclude IPs for privacy
		envInfo.MACs = []string{} // Exclude MACs for privacy
		jsonOutput, err = json.MarshalIndent(envInfo, "", "  ")
		if err != nil {
			log.Warn("Failed to marshal host info")
		}
		log.Infof("Host Info: %s", string(jsonOutput))

		// Print the status of the checks
		log.Infof("Checks Status:")
		shared.PrintStates()

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
