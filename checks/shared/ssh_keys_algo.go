// Package shared provides SSH key algo utilities.
package shared

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	sharedG "github.com/ParetoSecurity/pareto-core/shared"
)

// KeyType represents the type of SSH key.
type KeyType string

const (
	// Ed25519 is the Ed25519 key algorithm.
	Ed25519 KeyType = "ED25519"
	// Ed25519Sk is the Ed25519-SK key algorithm.
	Ed25519Sk KeyType = "ED25519-SK"
	// Ecdsa is the ECDSA key algorithm.
	Ecdsa KeyType = "ECDSA"
	// EcdsaSk is the ECDSA-SK key algorithm.
	EcdsaSk KeyType = "ECDSA-SK"
	// Dsa is the DSA key algorithm.
	Dsa KeyType = "DSA"
	// Rsa is the RSA key algorithm.
	Rsa KeyType = "RSA"
)

// KeyInfo holds information about a key.
type KeyInfo struct {
	strength  int
	signature string
	keyType   KeyType
}

// SSHKeysAlgo runs the SSH keys algorithm.
type SSHKeysAlgo struct {
	passed  bool
	sshKey  string
	sshPath string
}

// Name returns the name of the check
func (f *SSHKeysAlgo) Name() string {
	return "SSH keys have sufficient algorithm strength"
}

func parseKeyInfo(output string) KeyInfo {
	parts := strings.Fields(strings.TrimSpace(output))
	if len(parts) < 4 {
		return KeyInfo{}
	}

	strength, _ := strconv.Atoi(parts[0])
	keyType := strings.ToUpper(parts[len(parts)-1])
	keyType = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, keyType)

	return KeyInfo{
		strength:  strength,
		signature: parts[1],
		keyType:   KeyType(keyType),
	}
}

func (f *SSHKeysAlgo) isKeyStrong(path string) bool {
	output, err := sharedG.RunCommand("ssh-keygen", "-l", "-f", path)
	if err != nil {
		return false
	}

	info := parseKeyInfo(string(output))
	switch info.keyType {
	case Rsa:
		return info.strength >= 2048
	case Dsa:
		return info.strength >= 8192
	case Ecdsa, EcdsaSk:
		return info.strength >= 521
	case Ed25519, Ed25519Sk:
		return info.strength >= 256
	default:
		return false
	}
}

// Run executes the check
func (f *SSHKeysAlgo) Run() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	f.sshPath = filepath.Join(home, ".ssh")
	entries, err := os.ReadDir(f.sshPath)
	if err != nil {
		return err
	}

	f.passed = true
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".pub") {
			continue
		}

		pubPath := filepath.Join(f.sshPath, entry.Name())
		privPath := strings.TrimSuffix(pubPath, ".pub")

		if _, err := os.Stat(privPath); os.IsNotExist(err) {
			continue
		}

		if !f.isKeyStrong(pubPath) {
			f.passed = false
			f.sshKey = strings.TrimSuffix(entry.Name(), ".pub")
			break
		}
	}

	return nil
}

// Passed returns the status of the check
func (f *SSHKeysAlgo) Passed() bool {
	return f.passed
}

// IsRunnable returns whether SSHKeysAlgo is runnable.
func (f *SSHKeysAlgo) IsRunnable() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	sshPath := filepath.Join(home, ".ssh")
	if _, err := os.Stat(sshPath); os.IsNotExist(err) {
		return false
	}

	//check if there are any private keys in the .ssh directory
	files, err := os.ReadDir(sshPath)
	if err != nil {
		return false
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".pub") {
			privateKeyPath := filepath.Join(sshPath, strings.TrimSuffix(file.Name(), ".pub"))
			if _, err := os.Stat(privateKeyPath); err == nil {
				return true
			}
		}
	}
	return false
}

// UUID returns the UUID of the check
func (f *SSHKeysAlgo) UUID() string {
	return "ef69f752-0e89-46e2-a644-310429ae5f45"
}

// PassedMessage returns the message to return if the check passed
func (f *SSHKeysAlgo) PassedMessage() string {
	return "SSH keys use strong encryption"
}

// FailedMessage returns the message to return if the check failed
func (f *SSHKeysAlgo) FailedMessage() string {
	return "SSH keys are using weak encryption"
}

// RequiresRoot returns whether the check requires root access
func (f *SSHKeysAlgo) RequiresRoot() bool {
	return false
}

// Status returns the status of the check
func (f *SSHKeysAlgo) Status() string {
	if f.Passed() {
		return f.PassedMessage()
	}

	return "SSH key " + f.sshKey + " is using weak encryption"
}
