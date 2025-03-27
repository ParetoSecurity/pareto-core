package systemd

import (
	"github.com/ParetoSecurity/agent/shared"
)

func isEnabled(service string) bool {
	state, err := shared.RunCommand("systemctl", "--user", "is-enabled", service)
	if state == "enabled" && err == nil {
		return true
	}
	return false
}

func enable(service string) error {
	_, err := shared.RunCommand("systemctl", "--user", "enable", service)
	return err
}

func disable(service string) error {
	_, err := shared.RunCommand("systemctl", "--user", "disable", service)
	return err
}
