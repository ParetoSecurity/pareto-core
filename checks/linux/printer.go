package checks

import (
	"fmt"

	sharedchecks "github.com/ParetoSecurity/agent/checks/shared"
	"github.com/caarlos0/log"
)

type Printer struct {
	passed bool
	ports  map[int]string
}

// Name returns the name of the check
func (f *Printer) Name() string {
	return "Sharing printers is off"
}

// Run executes the check
func (f *Printer) Run() error {
	f.passed = true
	f.ports = make(map[int]string)

	// Samba, NFS and CUPS ports to check
	printService := map[int]string{
		631: "CUPS",
	}

	for port, service := range printService {
		if sharedchecks.CheckPort(port, "tcp") {
			log.WithField("check", f.Name()).WithField("port", port).WithField("service", service).Debug("Port open")
			f.passed = false
			f.ports[port] = service
		}
	}

	return nil
}

// Passed returns the status of the check
func (f *Printer) Passed() bool {
	return f.passed
}

// CanRun returns whether the check can run
func (f *Printer) IsRunnable() bool {
	return true
}

// UUID returns the UUID of the check
func (f *Printer) UUID() string {
	return "b96524e0-150b-4bb8-abc7-517051b6c14e"
}

// PassedMessage returns the message to return if the check passed
func (f *Printer) PassedMessage() string {
	return "Sharing printers is off"
}

// FailedMessage returns the message to return if the check failed
func (f *Printer) FailedMessage() string {
	return "Sharing printers is on"
}

// RequiresRoot returns whether the check requires root access
func (f *Printer) RequiresRoot() bool {
	return false
}

// Status returns the status of the check
func (f *Printer) Status() string {
	if !f.Passed() {
		msg := "Printer sharing services found running on ports:"
		for port, service := range f.ports {
			msg += fmt.Sprintf(" %s(%d)", service, port)
		}
		return msg
	}
	return f.PassedMessage()
}
