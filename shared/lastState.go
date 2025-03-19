package shared

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/caarlos0/log"
	"github.com/pelletier/go-toml"
)

type LastState struct {
	UUID    string `json:"uuid"`
	State   bool   `json:"state"`
	Details string `json:"details"`
}

var (
	mutex       sync.RWMutex
	states      = make(map[string]LastState)
	lastModTime time.Time
	statePath   string
)

func init() {
	states = make(map[string]LastState)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Warn("failed to get user home directory, using current directory instead")
		homeDir = "."
	}
	statePath = filepath.Join(homeDir, ".cache", "paretosecurity.state")
}

// Commit writes the current state map to the TOML file.
func CommitLastState() error {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.Create(statePath)
	if err != nil {
		return err
	}
	defer file.Close()
	lastModTime = time.Now()
	encoder := toml.NewEncoder(file)
	return encoder.Encode(states)
}

// AllChecksPassed returns true if all checks have passed.
func AllChecksPassed() bool {
	mutex.RLock()
	defer mutex.RUnlock()

	for _, state := range states {
		if !state.State {
			return false
		}
	}
	return true
}

// PrintStates loads and prints all stored states with their UUIDs, state values, and details.
func PrintStates() {
	loadStates()

	for uuid, state := range states {
		log.Infof("UUID: %s, State: %v, Details: %s", uuid, state.State, state.Details)
	}
}

// UpdateState updates the LastState struct in the in-memory map and commits to the TOML file.
func UpdateLastState(newState LastState) {
	mutex.Lock()
	defer mutex.Unlock()

	states[newState.UUID] = newState
}

// GetState retrieves the LastState struct by UUID.
func GetLastState(uuid string) (LastState, bool, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	loadStates()

	state, exists := states[uuid]
	return state, exists, nil
}

func GetLastStates() map[string]LastState {
	mutex.RLock()
	defer mutex.RUnlock()
	loadStates()

	return states
}

func GetModifiedTime() time.Time {
	mutex.RLock()
	defer mutex.RUnlock()
	loadStates()

	return lastModTime
}

func loadStates() {
	fileInfo, err := os.Stat(statePath)
	if err != nil {
		return
	}

	if fileInfo.ModTime().After(lastModTime) {
		file, err := os.Open(statePath)
		if err != nil {
			return
		}
		defer file.Close()

		decoder := toml.NewDecoder(file)
		if err := decoder.Decode(&states); err != nil {
			return
		}
		lastModTime = fileInfo.ModTime()
	}
}
