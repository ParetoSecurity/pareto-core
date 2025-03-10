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

var (
	shortVersion bool
	versionCmd   = &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(cmd *cobra.Command, args []string) {
			if shortVersion {
				log.Infof("%s@%s %s", shared.Version, shared.Commit, shared.Date)
				os.Exit(0)
			}

			log.Infof("%s@%s %s", shared.Version, shared.Commit, shared.Date)
			log.Infof("Built with %s", runtime.Version())

			device := shared.CurrentReportingDevice()
			log.Infof("Machine UUID: %s", device.MachineUUID)
			log.Infof("Name: %s", device.MachineName)
			log.Infof("OS Version: %s", device.OSVersion)
			log.Infof("Model Name: %s", device.ModelName)
			log.Infof("Model Serial: %s", device.ModelSerial)

			hostInfo, err := sysinfo.Host()
			if err != nil {
				log.Warn("Failed to get process information")
			}
			envInfo := hostInfo.Info()
			envInfo.IPs = []string{}
			envInfo.MACs = []string{}
			jsonOutput, err := json.MarshalIndent(envInfo, "", "  ")
			if err != nil {
				log.Warn("Failed to marshal host info")
			}
			log.Infof("Host Info: %s", string(jsonOutput))

			os.Exit(0)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&shortVersion, "short", false, "Only print the version")
}
