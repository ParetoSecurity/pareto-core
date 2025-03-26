package checks

import (
	"strings"

	"github.com/ParetoSecurity/agent/shared"
)

// PasswordToUnlock represents a check to ensure that a password is required to unlock the screen.
type PasswordToUnlock struct {
	passed bool
}

// Name returns the name of the check
func (f *PasswordToUnlock) Name() string {
	return "Password is required to unlock the screen"
}

func (f *PasswordToUnlock) checkGnome() bool {

	out, err := shared.RunCommand("gsettings", "get", "org.gnome.desktop.screensaver", "lock-enabled")
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

func (f *PasswordToUnlock) checkKDE() bool {
	out, err := shared.RunCommand("kreadconfig5", "--file", "kscreenlockerrc", "--group", "Daemon", "--key", "Autolock")
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

// Run executes the check
func (f *PasswordToUnlock) Run() error {
	anyCheckPerformed := false
	allChecksPassed := true

	// Check if running GNOME
	if _, err := lookPath("gsettings"); err == nil {
		anyCheckPerformed = true
		allChecksPassed = allChecksPassed && f.checkGnome()
	}

	// Check if running KDE
	if _, err := lookPath("kreadconfig5"); err == nil {
		anyCheckPerformed = true
		allChecksPassed = allChecksPassed && f.checkKDE()
	}

	// Performed at least one check and all performed checks passed
	f.passed = anyCheckPerformed && allChecksPassed
	return nil
}

// Passed returns the status of the check
func (f *PasswordToUnlock) Passed() bool {
	return f.passed
}

// IsRunnable returns whether the check can run
func (f *PasswordToUnlock) IsRunnable() bool {
	return true
}

// UUID returns the UUID of the check
func (f *PasswordToUnlock) UUID() string {
	return "37dee029-605b-4aab-96b9-5438e5aa44d8"
}

// PassedMessage returns the message to return if the check passed
func (f *PasswordToUnlock) PassedMessage() string {
	return "Password after sleep or screensaver is on"
}

// FailedMessage returns the message to return if the check failed
func (f *PasswordToUnlock) FailedMessage() string {
	return "Password after sleep or screensaver is off"
}

// RequiresRoot returns whether the check requires root access
func (f *PasswordToUnlock) RequiresRoot() bool {
	return false
}

// Status returns the status of the check
func (f *PasswordToUnlock) Status() string {
	if f.Passed() {
		return f.PassedMessage()
	}
	return f.FailedMessage()
}
