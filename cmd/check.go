package cmd

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ParetoSecurity/pareto-core/claims"
	"github.com/ParetoSecurity/pareto-core/runner"
	shared "github.com/ParetoSecurity/pareto-core/shared"
	team "github.com/ParetoSecurity/pareto-core/team"
	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check [--json] [--schema] [--install] [--uninstall]",
	Short: "Run checks",
	Run: func(cc *cobra.Command, args []string) {
		jsonOutput, _ := cc.Flags().GetBool("json")
		schemaOutput, _ := cc.Flags().GetBool("schema")
		installFlag, _ := cc.Flags().GetBool("install")
		uninstallFlag, _ := cc.Flags().GetBool("uninstall")
		if testing.Testing() {
			return
		}
		checkCommand(jsonOutput, schemaOutput, installFlag, uninstallFlag)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().Bool("json", false, "output JSON")
	checkCmd.Flags().Bool("schema", false, "output schema for all checks")
	checkCmd.Flags().Bool("install", false, "setup hourly checks")
	checkCmd.Flags().Bool("uninstall", false, "remove hourly checks")
}

func showLinkingMessage() {
	log.Info("To link your account with the team, please run `paretosecurity link`.")
	log.Info("For more information, please visit https://paretosecurity.com/dashboard")
}

func checkCommand(jsonOutput bool, schemaOutput bool, installFlag bool, uninstallFlag bool) {

	if shared.IsRoot() {
		log.Warn("Please run this command as a normal user, as it won't report all checks correctly.")
	}

	if installFlag {
		installUserTimer()
		return
	}
	if uninstallFlag {
		uninstallUserTimer()
		return
	}
	if schemaOutput {
		runner.PrintSchemaJSON(claims.All)
		return
	}

	if jsonOutput {
		runner.CheckJSON(claims.All)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	done := make(chan struct{})
	go func() {
		runner.Check(ctx, claims.All)
		close(done)
	}()

	select {
	case <-done:
		if shared.IsLinked() {
			err := team.ReportToTeam(false)
			if err != nil {
				log.WithError(err).Warn("failed to report to team")
			}
		} else {
			showLinkingMessage()
		}
		if !isUserTimerInstalled() {
			log.Info("To ensure your system is checked every hour, please run `paretosecurity check --install` to set it up.")
		}
	case <-ctx.Done():
		log.Warn("Check run timed out")
		os.Exit(1)
	}
}
