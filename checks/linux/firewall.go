package checks

import (
	"os/exec"
	"strings"

	"github.com/ParetoSecurity/agent/shared"
	"github.com/caarlos0/log"
)

// Firewall checks the system firewall.
type Firewall struct {
	passed bool
	status string
}

// Name returns the name of the check
func (f *Firewall) Name() string {
	return "Firewall is on"
}

func (f *Firewall) checkUFW() bool {
	output, err := shared.RunCommand("ufw", "status")
	if err != nil {
		log.WithError(err).WithField("output", output).Warn("Failed to check UFW status")
		return false
	}
	log.WithField("output", output).Debug("UFW status")
	return strings.Contains(output, "Status: active")
}

func (f *Firewall) checkFirewalld() bool {
	output, err := shared.RunCommand("systemctl", "is-active", "firewalld")
	if err != nil {
		log.WithError(err).WithField("output", output).Warn("Failed to check firewalld status")
		return false
	}
	log.WithField("output", output).Debug("Firewalld status")
	return output == "active"
}

// Run executes the check
func (f *Firewall) Run() error {
	if f.RequiresRoot() && !shared.IsRoot() {
		log.Debug("Running check via helper")
		// Run as root
		passed, err := shared.RunCheckViaHelper(f.UUID())
		if err != nil {
			log.WithError(err).Warn("Failed to run check via helper")
			return err
		}
		f.passed = passed
		return nil
	}

	log.Debug("Running check directly")
	f.passed = false

	if !f.passed {
		f.passed = f.checkUFW()
	}

	if !f.passed {
		f.passed = f.checkFirewalld()
	}

	if !f.passed {
		f.status = f.FailedMessage()
	}

	return nil
}

// Passed returns the status of the check
func (f *Firewall) Passed() bool {
	return f.passed
}

// IsRunnable returns whether Firewall is runnable.
func (f *Firewall) IsRunnable() bool {

	can := shared.IsSocketServicePresent()
	if !can {
		f.status = "Root helper is not available, check cannot run. See https://paretosecurity.com/root-helper for more information."
		return false
	}

	// Check if ufw or firewalld are present
	_, errUFW := exec.LookPath("ufw")
	_, errFirewalld := exec.LookPath("firewalld")
	if errUFW != nil && errFirewalld != nil {
		f.status = "Neither ufw nor firewalld are present, check cannot run"
		return false
	}

	return true
}

// UUID returns the UUID of the check
func (f *Firewall) UUID() string {
	return "2e46c89a-5461-4865-a92e-3b799c12034a"
}

// PassedMessage returns the message to return if the check passed
func (f *Firewall) PassedMessage() string {
	return "Firewall is on"
}

// FailedMessage returns the message to return if the check failed
func (f *Firewall) FailedMessage() string {
	return "Firewall is off"
}

// RequiresRoot returns whether the check requires root access
func (f *Firewall) RequiresRoot() bool {
	return true
}

// Status returns the status of the check
func (f *Firewall) Status() string {
	if f.Passed() {
		return f.PassedMessage()
	}
	return f.status
}
