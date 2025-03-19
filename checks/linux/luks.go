package checks

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/ParetoSecurity/agent/shared"
	"github.com/caarlos0/log"
)

type EncryptingFS struct {
	passed bool
	status string
}

// Name returns the name of the check
func (f *EncryptingFS) Name() string {
	return "Filesystem encryption is enabled"
}

// Passed returns the status of the check
func (f *EncryptingFS) Passed() bool {
	return f.passed
}

// CanRun returns whether the check can run
func (f *EncryptingFS) IsRunnable() bool {
	can := shared.IsSocketServicePresent()
	if !can {
		f.status = "Root helper is not available, check cannot run. See https://paretosecurity.com/docs/linux/root-helper for more information."
	}
	return can
}

// UUID returns the UUID of the check
func (f *EncryptingFS) UUID() string {
	return "21830a4e-84f1-48fe-9c5b-beab436b2cdb"
}

// PassedMessage returns the message to return if the check passed
func (f *EncryptingFS) PassedMessage() string {
	return "Block device encryption is enabled"
}

// FailedMessage returns the message to return if the check failed
func (f *EncryptingFS) FailedMessage() string {
	return "Block device encryption is disabled"
}

// RequiresRoot returns whether the check requires root access
func (f *EncryptingFS) RequiresRoot() bool {
	return true
}

// Run executes the check
func (f *EncryptingFS) Run() error {

	if f.RequiresRoot() && !shared.IsRoot() {
		log.Debug("Running check via helper")
		// Run as root
		passed, err := shared.RunCheckViaHelper(f.UUID())
		if err != nil {
			log.WithError(err).Warn("Failed to run check via helper")
			return err
		}
		f.passed = passed
		return nil
	}
	log.Debug("Running check directly")
	encryptedDevices := make(map[string]string)

	// Read crypttab to get encrypted devices
	crypttab, err := os.Open("/etc/crypttab")
	if err == nil {
		scanner := bufio.NewScanner(crypttab)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				encryptedDevices[fields[0]] = strings.TrimPrefix(strings.Trim(fields[1], `"`), "UUID=")
			}
		}
		crypttab.Close()
	}
	log.WithField("encryptedDevices", encryptedDevices).Debug("Found encrypted devices")
	cmd := exec.Command("blkid")
	output, err := cmd.Output()
	if err != nil {
		log.WithError(err).Warn("Failed to run blkid")
		return err
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, `TYPE="crypto_LUKS"`) {
			log.WithField("line", line).Debug("Found encrypted device")
			for _, uuid := range encryptedDevices {
				if strings.Contains(line, uuid) {
					f.passed = true
					f.status = f.PassedMessage()
					return nil
				}
			}
		}
	}

	f.passed = maybeCryptoViaKernel()
	f.status = f.FailedMessage()

	return nil
}

// Status returns the status of the check
func (f *EncryptingFS) Status() string {
	if f.Passed() {
		return f.PassedMessage()
	}
	return f.status
}

func maybeCryptoViaKernel() bool {
	// Read kernel parameters to check if root is booted via crypt
	cmdline, err := shared.ReadFile("/proc/cmdline")
	if err != nil {
		log.WithError(err).Warn("Failed to read /proc/cmdline")
	}

	params := strings.Fields(string(cmdline))
	for _, param := range params {
		if strings.HasPrefix(param, "cryptdevice=") {
			parts := strings.Split(param, ":")
			if len(parts) == 3 && parts[2] == "root" {
				return true
			}
		}
	}
	return false
}
