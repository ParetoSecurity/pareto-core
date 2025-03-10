// Package main provides the entry point for the application.
package main

import (
	"github.com/caarlos0/log"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/ParetoSecurity/agent/cmd"
	shared "github.com/ParetoSecurity/agent/shared"
)

func init() {

	// automatically set GOMAXPROCS to match available CPUs.
	// GOMAXPROCS will be used as the default value for the --parallelism flag.
	if _, err := maxprocs.Set(); err != nil {
		log.WithError(err).Fatal("failed to set GOMAXPROCS")
	}

}
func main() {
	if err := shared.LoadConfig(); err != nil {
		log.WithError(err).Warn("failed to load config")
	}
	cmd.Execute()
}
