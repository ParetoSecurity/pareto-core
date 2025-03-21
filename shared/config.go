package shared

import (
	"os"
	"path/filepath"

	"github.com/caarlos0/log"
	"github.com/pelletier/go-toml"
	"github.com/samber/lo"
)

var Config ParetoConfig
var configPath string

type CheckStatus struct {
	Disabled bool
}

type ParetoConfig struct {
	TeamID    string
	AuthToken string
	ReportURL string
}

func init() {
	states = make(map[string]LastState)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Warn("failed to get user home directory, using current directory instead")
		homeDir = "."
	}
	configPath = filepath.Join(homeDir, ".config", "pareto.toml")
	log.Debugf("configPath: %s", configPath)
}

func SaveConfig() error {

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	return encoder.Encode(Config)
}

func LoadConfig() error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := SaveConfig(); err != nil {
			return err
		}
		return nil
	}
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := toml.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		return err
	}

	// Set default values if not set
	if lo.IsEmpty(Config.ReportURL) {
		Config.ReportURL = "https://dash.paretosecurity.com"
	}

	return nil
}
