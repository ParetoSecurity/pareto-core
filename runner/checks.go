package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ParetoSecurity/agent/check"
	"github.com/ParetoSecurity/agent/claims"
	"github.com/ParetoSecurity/agent/shared"
	"github.com/caarlos0/log"
	"github.com/pterm/pterm"
)

// Check runs a series of checks concurrently for a list of claims.
//
// It iterates over each claim provided in claimsTorun and, for each claim,
// over its associated checks. Each check is executed in its own goroutine.
func Check(ctx context.Context, claimsTorun []claims.Claim) {
	multi := pterm.DefaultMultiPrinter
	var wg sync.WaitGroup
	log.Info("Starting checks...")
	if _, err := multi.Start(); err != nil {
		log.WithError(err).Warn("failed to stop multi printer")
	}
	for _, claim := range claimsTorun {
		for _, chk := range claim.Checks {
			wg.Add(1)
			go func(claim claims.Claim, chk check.Check) {
				defer wg.Done()
				select {
				case <-ctx.Done():
					return
				default:
					spinner, _ := pterm.DefaultSpinner.WithWriter(multi.NewWriter()).Start(fmt.Sprintf("%s: %s", claim.Title, chk.Name()))
					spinner.FailPrinter = &pterm.PrefixPrinter{
						MessageStyle: &pterm.Style{pterm.FgLightRed},
						Prefix: pterm.Prefix{
							Style: &pterm.Style{pterm.BgRed, pterm.FgLightRed},
							Text:  "✗",
						},
					}
					spinner.SuccessPrinter = &pterm.PrefixPrinter{
						MessageStyle: &pterm.Style{pterm.FgLightGreen},
						Prefix: pterm.Prefix{
							Style: &pterm.Style{pterm.BgGreen, pterm.FgLightGreen},
							Text:  "✓",
						},
					}

					// Skip checks that are not runnable
					if !chk.IsRunnable() {
						spinner.Warning(pterm.White(claim.Title), pterm.White(": "), pterm.Blue(fmt.Sprintf("%s > ", chk.Name())), pterm.Yellow("skipped"))
						return
					}

					if err := chk.Run(); err != nil {
						spinner.Fail(pterm.White(claim.Title), pterm.White(": "), pterm.Blue(fmt.Sprintf("%s > ", chk.Name())), pterm.Red(err.Error()))
					}

					if chk.Passed() {
						spinner.Success(pterm.White(claim.Title), pterm.White(": "), pterm.Green(chk.Status()))
					} else {
						spinner.Fail(pterm.White(claim.Title), pterm.White(": "), pterm.Blue(fmt.Sprintf("%s > ", chk.Name())), pterm.Red(chk.Status()))
					}
					shared.UpdateLastState(shared.LastState{
						UUID:    chk.UUID(),
						State:   chk.Passed(),
						Details: chk.Status(),
					})
				}
			}(claim, chk)
		}
	}
	wg.Wait()
	if err := shared.CommitLastState(); err != nil {
		log.WithError(err).Warn("failed to commit last state")
	}
	if _, err := multi.Stop(); err != nil {
		log.WithError(err).Warn("failed to stop multi printer")
	}

	log.Info("Checks completed.")
	if err := shared.SaveConfig(); err != nil {
		log.WithError(err).Warn("cannot save config")
	}
}

// CheckJSON validates and executes each check defined in the provided slice of claims.
// It iterates over each claim in claimsTorun and for every check within a claim:
//   - If the check is not runnable (determined by IsRunnable), it tags the check with a "skipped" status.
//   - Otherwise, it runs the check using Run. If Run returns an error, the check's status is set to the error message.
//   - If no error occurs, the check's Pass/Fail result (via Passed) is recorded as "passed" or "failed".
//
// After processing all checks, the function attempts to save the configuration using shared.SaveConfig,
// logs a warning if saving fails, marshals the status map into indented JSON,
// logs any marshalling error, and prints the final JSON.
func CheckJSON(claimsTorun []claims.Claim) {
	status := make(map[string]string)
	for _, claim := range claimsTorun {
		for _, chk := range claim.Checks {

			if !chk.IsRunnable() {
				status[chk.UUID()] = "skipped"
				continue
			}

			if err := chk.Run(); err != nil {
				status[chk.UUID()] = err.Error()
				continue
			}
			if chk.Passed() {
				status[chk.UUID()] = "passed"
			} else {
				status[chk.UUID()] = "failed"
			}
		}
	}
	if err := shared.SaveConfig(); err != nil {
		log.WithError(err).Warn("cannot save config")
	}
	out, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		log.WithError(err).Warn("cannot marshal status")
	}
	fmt.Println(string(out))
}

// PrintSchemaJSON constructs and prints a JSON schema generated from a slice of claims.
// For each claim, the function builds a nested map where the claim's title is the key and its
// value is another map. This inner map associates each check's UUID with a slice that contains
// the check's passed message and failed message.
// The resulting schema is marshalled into an indented JSON string and printed to standard output.
// In case of an error during marshalling, the function logs a warning with the error details.
func PrintSchemaJSON(claimsTorun []claims.Claim) {
	schema := make(map[string]map[string][]string)
	for _, claim := range claimsTorun {
		checks := make(map[string][]string)
		for _, chk := range claim.Checks {
			checks[chk.UUID()] = []string{chk.PassedMessage(), chk.FailedMessage()}
		}
		schema[claim.Title] = checks
	}
	out, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		log.WithError(err).Warn("cannot marshal schema")
	}
	fmt.Println(string(out))
}
