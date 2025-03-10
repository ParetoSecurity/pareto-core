package checks

import (
	"testing"

	"github.com/ParetoSecurity/agent/shared"
	"github.com/stretchr/testify/assert"
)

func TestCheckSSHConfig(t *testing.T) {

	tests := []struct {
		name           string
		setupMocks     map[string]string
		expectedPassed bool
		expectedDetail string
	}{
		{
			name: "All ok",
			setupMocks: map[string]string{
				"sshd -T": "PasswordAuthentication no\nPermitRootLogin no",
			},
			expectedPassed: true,
			expectedDetail: "",
		},
		{
			name: "PasswordAuthentication is enabled",
			setupMocks: map[string]string{
				"sshd -T": "PasswordAuthentication yes\nPermitRootLogin no",
			},
			expectedPassed: false,
			expectedDetail: "PasswordAuthentication is enabled",
		},
		{
			name: "PermitRootLogin is enabled",
			setupMocks: map[string]string{
				"sshd -T": "PasswordAuthentication no\nPermitRootLogin yes",
			},
			expectedPassed: false,
			expectedDetail: "Root login is enabled",
		},

		{
			name: "PermitEmptyPasswords is enabled",
			setupMocks: map[string]string{
				"sshd -T": "PasswordAuthentication no\nPermitRootLogin no\nPermitEmptyPasswords yes",
			},
			expectedPassed: false,
			expectedDetail: "Empty passwords are allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.RunCommandMocks = tt.setupMocks
			lookPathMock = func(file string) (string, error) {
				return file, nil
			}
			su := &SSHConfigCheck{}

			err := su.Run()
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedPassed, su.passed)
			assert.Equal(t, tt.expectedDetail, su.status)
		})
	}
}
func TestIsRunnable(t *testing.T) {

	tests := []struct {
		name           string
		setupMocks     map[string]string
		expectedResult bool
	}{
		{
			name: "sshd service is active",
			setupMocks: map[string]string{
				"systemctl is-active sshd": "active",
			},
			expectedResult: true,
		},
		{
			name: "ssh service is active",
			setupMocks: map[string]string{
				"systemctl is-active sshd": "inactive",
				"systemctl is-active ssh":  "active",
			},
			expectedResult: true,
		},
		{
			name: "sshd.socket is enabled",
			setupMocks: map[string]string{
				"systemctl is-active sshd":         "inactive",
				"systemctl is-active ssh":          "inactive",
				"systemctl is-enabled sshd.socket": "enabled",
			},
			expectedResult: true,
		},
		{
			name: "ssh.socket is enabled",
			setupMocks: map[string]string{
				"systemctl is-active sshd":         "inactive",
				"systemctl is-active ssh":          "inactive",
				"systemctl is-enabled sshd.socket": "disabled",
				"systemctl is-enabled ssh.socket":  "enabled",
			},
			expectedResult: true,
		},
		{
			name: "all services are inactive or disabled",
			setupMocks: map[string]string{
				"systemctl is-active sshd":         "inactive",
				"systemctl is-active ssh":          "inactive",
				"systemctl is-enabled sshd.socket": "disabled",
				"systemctl is-enabled ssh.socket":  "disabled",
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.RunCommandMocks = tt.setupMocks
			lookPathMock = func(file string) (string, error) {
				return file, nil
			}
			su := &SSHConfigCheck{}

			result := su.IsRunnable()
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestSSHConfigCheck_Name(t *testing.T) {
	su := &SSHConfigCheck{}
	expectedName := "SSH Server Configuration is Secure"
	if su.Name() != expectedName {
		t.Errorf("Expected Name %s, got %s", expectedName, su.Name())
	}
}

func TestSSHConfigCheck_Status(t *testing.T) {
	su := &SSHConfigCheck{}
	expectedStatus := ""
	if su.Status() != expectedStatus {
		t.Errorf("Expected Status %s, got %s", expectedStatus, su.Status())
	}
}

func TestSSHConfigCheck_UUID(t *testing.T) {
	su := &SSHConfigCheck{}
	expectedUUID := "da4edd80-6af0-4fb3-9fc7-f9a0e9d07f3b"
	if su.UUID() != expectedUUID {
		t.Errorf("Expected UUID %s, got %s", expectedUUID, su.UUID())
	}
}

func TestSSHConfigCheck_Passed(t *testing.T) {
	su := &SSHConfigCheck{passed: true}
	expectedPassed := true
	if su.Passed() != expectedPassed {
		t.Errorf("Expected Passed %v, got %v", expectedPassed, su.Passed())
	}
}

func TestSSHConfigCheck_FailedMessage(t *testing.T) {
	su := &SSHConfigCheck{}
	expectedFailedMessage := "SSH configuration is not secure."
	if su.FailedMessage() != expectedFailedMessage {
		t.Errorf("Expected FailedMessage %s, got %s", expectedFailedMessage, su.FailedMessage())
	}
}

func TestSSHConfigCheck_PassedMessage(t *testing.T) {
	su := &SSHConfigCheck{}
	expectedPassedMessage := "SSH configuration is secure."
	if su.PassedMessage() != expectedPassedMessage {
		t.Errorf("Expected PassedMessage %s, got %s", expectedPassedMessage, su.PassedMessage())
	}
}
