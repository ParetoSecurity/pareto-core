package checks

import (
	"strings"

	"github.com/ParetoSecurity/agent/shared"
	"github.com/samber/lo"
)

type DockerAccess struct {
	passed bool
	status string
}

// Name returns the name of the check
func (f *DockerAccess) Name() string {
	return "Access to Docker is restricted"
}

// Run executes the check
func (f *DockerAccess) Run() error {
	output, err := shared.RunCommand("docker", "info", "--format", "{{.SecurityOptions}}")
	if err != nil || lo.IsEmpty(output) {
		f.passed = false
		f.status = "Failed to get Docker info"
		return err
	}

	if !strings.Contains(output, "rootless") {
		f.passed = false
		f.status = f.FailedMessage()
		return nil
	}

	f.passed = true

	return nil
}

// Passed returns the status of the check
func (f *DockerAccess) Passed() bool {
	return f.passed
}

// CanRun returns whether the check can run
func (f *DockerAccess) IsRunnable() bool {

	out, _ := shared.RunCommand("docker", "version")
	if !strings.Contains(out, "Version") {
		f.status = "Docker is not installed"
		return false
	}

	return true
}

// UUID returns the UUID of the check
func (f *DockerAccess) UUID() string {
	return "25443ceb-c1ec-408c-b4f3-2328ea0c84e1"
}

// PassedMessage returns the message to return if the check passed
func (f *DockerAccess) PassedMessage() string {
	return "Docker is running in rootless mode"
}

// FailedMessage returns the message to return if the check failed
func (f *DockerAccess) FailedMessage() string {
	return "Docker is not running in rootless mode"
}

// RequiresRoot returns whether the check requires root access
func (f *DockerAccess) RequiresRoot() bool {
	return false
}

// Status returns the status of the check
func (f *DockerAccess) Status() string {
	if !f.Passed() {
		return f.status
	}
	return f.PassedMessage()
}
