package shared

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml"
)

func TestSaveConfig_Success(t *testing.T) {
	// Create a temporary directory for testing.
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Set configPath to a temporary file.
	configPath = filepath.Join(tempDir, "pareto.toml")

	// Prepare a test configuration.

	Config = ParetoConfig{
		TeamID:    "team1",
		AuthToken: "token1",
	}

	// Call SaveConfig.
	if err := SaveConfig(); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Read the written file.
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	// Unmarshal the file content.
	var loadedConfig ParetoConfig
	if err := toml.Unmarshal(data, &loadedConfig); err != nil {
		t.Fatalf("failed to decode config file: %v", err)
	}

	// Validate the saved configuration.
	if loadedConfig.TeamID != Config.TeamID {
		t.Errorf("expected TeamID %q, got %q", Config.TeamID, loadedConfig.TeamID)
	}
	if loadedConfig.AuthToken != Config.AuthToken {
		t.Errorf("expected AuthToken %q, got %q", Config.AuthToken, loadedConfig.AuthToken)
	}

}

func TestSaveConfig_Failure(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "config-test-failure")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Set configPath to a directory to simulate a failure (os.Create should fail).
	configPath = tempDir

	// Call SaveConfig expecting an error.
	if err := SaveConfig(); err == nil {
		t.Errorf("expected error when configPath is a directory, got nil")
	}
}
func TestLoadConfig_NonExistent(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config-test-load")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Set configPath to a non-existent file
	configPath = filepath.Join(tempDir, "non-existent.toml")

	// Reset Config
	Config = ParetoConfig{}

	// LoadConfig should create a default config file
	if err := LoadConfig(); err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Check if file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("config file was not created")
	}
}

func TestLoadConfig_ExistingFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config-test-load-existing")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Set configPath to a temporary file
	configPath = filepath.Join(tempDir, "pareto.toml")

	// Create a test config file
	testConfig := ParetoConfig{
		TeamID:    "testteam",
		AuthToken: "testtoken",
	}

	// Write test config to file
	data, err := toml.Marshal(testConfig)
	if err != nil {
		t.Fatalf("failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	// Reset Config
	Config = ParetoConfig{}

	// Load the config
	if err := LoadConfig(); err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Validate loaded configuration
	if Config.TeamID != testConfig.TeamID {
		t.Errorf("expected TeamID %q, got %q", testConfig.TeamID, Config.TeamID)
	}
	if Config.AuthToken != testConfig.AuthToken {
		t.Errorf("expected AuthToken %q, got %q", testConfig.AuthToken, Config.AuthToken)
	}
}

func TestLoadConfig_DefaultValues(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config-test-default-values")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Set configPath to a temporary file
	configPath = filepath.Join(tempDir, "pareto.toml")

	// Create a test config file with missing ReportURL
	testConfig := ParetoConfig{
		TeamID:    "testteam",
		AuthToken: "testtoken",
		// ReportURL intentionally omitted
	}

	// Write test config to file
	data, err := toml.Marshal(testConfig)
	if err != nil {
		t.Fatalf("failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	// Reset Config
	Config = ParetoConfig{}

	// Load the config
	if err := LoadConfig(); err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Validate default values were set
	if Config.ReportURL != "https://dash.paretosecurity.com" {
		t.Errorf("expected default ReportURL %q, got %q", "https://dash.paretosecurity.com", Config.ReportURL)
	}
}

func TestLoadConfig_InvalidFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config-test-invalid")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Set configPath to a temporary file
	configPath = filepath.Join(tempDir, "pareto.toml")

	// Write invalid TOML content
	if err := os.WriteFile(configPath, []byte("invalid toml content"), 0644); err != nil {
		t.Fatalf("failed to write invalid config file: %v", err)
	}

	// LoadConfig should return an error
	if err := LoadConfig(); err == nil {
		t.Errorf("expected error when loading invalid config, got nil")
	}
}
