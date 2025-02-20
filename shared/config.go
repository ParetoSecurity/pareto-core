package shared

import (
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/log"
	"github.com/pelletier/go-toml"
)

var Config ParetoConfig
var configPath string

type CheckStatus struct {
	UpdatedAt time.Time
	Passed    bool
	Disabled  bool
}

type ParetoConfig struct {
	TeamID    string
	AuthToken string
	Checks    map[string]CheckStatus
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

	return nil
}
