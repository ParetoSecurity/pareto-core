//go:build linux
// +build linux

package cmd

import (
	"os"
	"path/filepath"

	"github.com/ParetoSecurity/agent/shared"
	"github.com/caarlos0/log"
)

// Define a variable to hold the original os.UserHomeDir function
var userHomeDir = os.UserHomeDir

const timerContent = `[Unit]
Description=Timer for pareto-core hourly execution

[Timer]
OnCalendar=hourly
Persistent=true

[Install]
WantedBy=timers.target`

const localServiceContent = `[Unit]
Description=Service for pareto-core

[Service]
Type=oneshot
ExecStart=/usr/bin/paretosecurity check
StandardInput=null

[Install]
WantedBy=timers.target`

func isUserTimerInstalled() bool {
	homeDir, err := userHomeDir()
	if err != nil {
		log.WithError(err).Warn("Failed to get home directory")
		return false
	}

	systemdPath := filepath.Join(homeDir, ".config", "systemd", "user")
	if _, err := os.Stat(filepath.Join(systemdPath, "pareto-core.timer")); err == nil {
		return true
	}
	return false
}

func installUserTimer() {
	homeDir, err := userHomeDir()
	if err != nil {
		log.WithError(err).Fatal("Failed to get home directory")
		return
	}

	systemdPath := filepath.Join(homeDir, ".config", "systemd", "user")
	if err := os.MkdirAll(systemdPath, 0755); err != nil {
		log.WithError(err).Fatal("Failed to create systemd user directory")
		return
	}

	// Create timer file
	timerPath := filepath.Join(systemdPath, "pareto-core.timer")
	if err := os.WriteFile(timerPath, []byte(timerContent), 0644); err != nil {
		log.WithError(err).Fatal("Failed to create timer file")
		return
	}

	// Create service file
	servicePath := filepath.Join(systemdPath, "pareto-core.service")
	if err := os.WriteFile(servicePath, []byte(localServiceContent), 0644); err != nil {
		log.WithError(err).Fatal("Failed to create service file")
		return
	}

	// Execute commands
	if out, err := shared.RunCommand("systemctl", "--user", "daemon-reload"); err != nil {
		log.WithError(err).WithField("out", out).Fatal("Failed to reload systemd")
		return
	}
	if out, err := shared.RunCommand("systemctl", "--user", "enable", "--now", "pareto-core.timer"); err != nil {
		log.WithError(err).WithField("out", out).Fatal("Failed to enable and start timer")
		return
	}

	log.Info("Timer installed successfully, to enable it run:")
	log.Infof("sudo loginctl enable-linger %s", os.Getenv("USER"))
}

func uninstallUserTimer() {
	// Logic to uninstall the user timer
	homeDir, err := userHomeDir()
	if err != nil {
		log.WithError(err).Fatal("Failed to get home directory:")
		return
	}

	systemdPath := filepath.Join(homeDir, ".config", "systemd", "user")
	if err := os.MkdirAll(systemdPath, 0755); err != nil {
		log.WithError(err).Fatal("Failed to create systemd user directory")
		return
	}

	// Create timer file
	timerPath := filepath.Join(systemdPath, "pareto-coretimer")
	if err := os.Remove(timerPath); err != nil {
		log.WithError(err).Fatal("Failed to remove timer file")
		return
	}

	// Create service file
	servicePath := filepath.Join(systemdPath, "pareto-coreservice")
	if err := os.Remove(servicePath); err != nil {
		log.WithError(err).Fatal("Failed to remove service file")
		return
	}
	// Execute commands
	if out, err := shared.RunCommand("systemctl", "--user", "daemon-reload"); err != nil {
		log.WithError(err).WithField("out", out).Fatal("Failed to reload systemd")
		return
	}
}
