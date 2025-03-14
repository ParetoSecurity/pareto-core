//go:build linux
// +build linux

package cmd

import (
	"os"
	"path/filepath"
	"testing"

	shared "github.com/ParetoSecurity/agent/shared"
	"github.com/stretchr/testify/assert"
)

func TestIsUserTimerInstalled(t *testing.T) {
	// Setup: Create a temporary directory and necessary files
	tempDir, err := os.MkdirTemp("", "test-systemd")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after the test

	systemdPath := filepath.Join(tempDir, ".config", "systemd", "user")
	err = os.MkdirAll(systemdPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create systemd user directory: %v", err)
	}

	timerFilePath := filepath.Join(systemdPath, "pareto-core.timer")

	// Test case 1: Timer file exists
	_, err = os.Create(timerFilePath)
	if err != nil {
		t.Fatalf("Failed to create timer file: %v", err)
	}

	// Mock os.UserHomeDir to return the temporary directory
	originalUserHomeDir := userHomeDir
	userHomeDir = func() (string, error) {
		return tempDir, nil
	}
	defer func() { userHomeDir = originalUserHomeDir }() // Restore original function

	assert.True(t, isUserTimerInstalled(), "Expected true when timer file exists")

	// Test case 2: Timer file does not exist
	err = os.Remove(timerFilePath)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to remove timer file: %v", err)
	}

	assert.False(t, isUserTimerInstalled(), "Expected false when timer file does not exist")

	// Test case 3: os.UserHomeDir returns an error
	userHomeDir = func() (string, error) {
		return "", os.ErrNotExist // Simulate an error
	}

	assert.False(t, isUserTimerInstalled(), "Expected false when UserHomeDir fails")
}

func TestInstallUserTimer(t *testing.T) {
	// Setup: Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test-install-timer")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Mock userHomeDir
	originalUserHomeDir := userHomeDir
	userHomeDir = func() (string, error) {
		return tempDir, nil
	}
	defer func() { userHomeDir = originalUserHomeDir }()

	// Mock shared.RunCommand to prevent actual systemctl calls
	shared.RunCommandMocks = map[string]string{
		"systemctl --user daemon-reload":                  "",
		"systemctl --user enable --now pareto-core.timer": "",
	}

	// Call the function to test
	installUserTimer()

	// Verify timer file was created correctly
	timerPath := filepath.Join(tempDir, ".config", "systemd", "user", "pareto-core.timer")
	timerContent, err := os.ReadFile(timerPath)
	assert.NoError(t, err, "Timer file should exist")
	assert.Contains(t, string(timerContent), "Description=Timer for pareto-core hourly execution")

	// Verify service file was created correctly
	servicePath := filepath.Join(tempDir, ".config", "systemd", "user", "pareto-core.service")
	serviceContent, err := os.ReadFile(servicePath)
	assert.NoError(t, err, "Service file should exist")
	assert.Contains(t, string(serviceContent), "Description=Service for pareto-core")

}
func TestUninstallUserTimer(t *testing.T) {
	// Setup: Create a temporary directory and necessary files
	tempDir, err := os.MkdirTemp("", "test-uninstall-timer")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the systemd directory structure
	systemdPath := filepath.Join(tempDir, ".config", "systemd", "user")
	err = os.MkdirAll(systemdPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create systemd user directory: %v", err)
	}

	// Create test files that should be removed
	timerPath := filepath.Join(systemdPath, "pareto-coretimer")
	servicePath := filepath.Join(systemdPath, "pareto-coreservice")

	if err := os.WriteFile(timerPath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test timer file: %v", err)
	}
	if err := os.WriteFile(servicePath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test service file: %v", err)
	}

	// Mock userHomeDir
	originalUserHomeDir := userHomeDir
	userHomeDir = func() (string, error) {
		return tempDir, nil
	}
	defer func() { userHomeDir = originalUserHomeDir }()

	// Mock shared.RunCommand to prevent actual systemctl calls
	shared.RunCommandMocks = map[string]string{
		"systemctl --user daemon-reload":                   "",
		"systemctl --user disable --now pareto-core.timer": "",
	}
	// Call the function to test
	uninstallUserTimer()

	// Verify files were removed
	_, timerErr := os.Stat(timerPath)
	assert.True(t, os.IsNotExist(timerErr), "Timer file should be removed")

	_, serviceErr := os.Stat(servicePath)
	assert.True(t, os.IsNotExist(serviceErr), "Service file should be removed")
}
