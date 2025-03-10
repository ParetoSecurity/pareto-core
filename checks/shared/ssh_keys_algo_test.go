package shared

import (
	"errors"
	"testing"

	sharedG "github.com/ParetoSecurity/agent/shared"
)

func TestIsKeyStrong(t *testing.T) {
	// Save original RunCommand and restore at end.

	tests := []struct {
		name       string
		output     string
		err        error
		expectPass bool
	}{
		{
			name:       "RSA meets requirement",
			output:     "2048 abc dummy RSA",
			err:        nil,
			expectPass: true,
		},
		{
			name:       "RSA meets requirement in parentheses",
			output:     "2048 SHA256:redacted .ssh/id_rsa.pub (RSA)",
			err:        nil,
			expectPass: true,
		},
		{
			name:       "RSA below requirement",
			output:     "2047 abc dummy RSA",
			err:        nil,
			expectPass: false,
		},
		{
			name:       "DSA meets requirement",
			output:     "8192 abc dummy DSA",
			err:        nil,
			expectPass: true,
		},
		{
			name:       "DSA below requirement",
			output:     "8191 abc dummy DSA",
			err:        nil,
			expectPass: false,
		},
		{
			name:       "ECDSA meets requirement",
			output:     "521 abc dummy ECDSA",
			err:        nil,
			expectPass: true,
		},
		{
			name:       "ECDSA below requirement",
			output:     "520 abc dummy ECDSA",
			err:        nil,
			expectPass: false,
		},
		{
			name:       "Ed25519 meets requirement",
			output:     "256 abc dummy ED25519",
			err:        nil,
			expectPass: true,
		},
		{
			name:       "Ed25519 meets requirement in parentheses",
			output:     "256 SHA256:redacted example@example.org (ED25519)",
			err:        nil,
			expectPass: true,
		},
		{
			name:       "Ed25519 below requirement",
			output:     "255 abc dummy ED25519",
			err:        nil,
			expectPass: false,
		},
		{
			name:       "Unknown key type",
			output:     "1024 abc dummy UNKNOWN",
			err:        nil,
			expectPass: false,
		},
		{
			name:       "RunCommand returns error",
			output:     "",
			err:        errors.New("command failed"),
			expectPass: false,
		},
		{
			name:       "Malformed output (less than 4 fields)",
			output:     "2048 abc",
			err:        nil,
			expectPass: false,
		},
	}

	algo := SSHKeysAlgo{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Override RunCommand to simulate different outputs.
			sharedG.RunCommandMocks = map[string]string{
				"ssh-keygen -l -f dummy/path": tc.output,
			}
			result := algo.isKeyStrong("dummy/path")
			if result != tc.expectPass {
				t.Errorf("expected %v, got %v for output %q", tc.expectPass, result, tc.output)
			}
		})
	}
	t.Run("Run command", func(t *testing.T) {
		_ = algo.IsRunnable()
		_ = algo.Run()
	})
}

func TestSSHKeysAlgo_Name(t *testing.T) {
	dockerAccess := &SSHKeysAlgo{}
	expectedName := "SSH keys have sufficient algorithm strength"
	if dockerAccess.Name() != expectedName {
		t.Errorf("Expected Name %s, got %s", expectedName, dockerAccess.Name())
	}
}

func TestSSHKeysAlgo_Status(t *testing.T) {
	dockerAccess := &SSHKeysAlgo{}
	expectedStatus := "SSH key  is using weak encryption"
	if dockerAccess.Status() != expectedStatus {
		t.Errorf("Expected Status %s, got %s", expectedStatus, dockerAccess.Status())
	}
}

func TestSSHKeysAlgo_UUID(t *testing.T) {
	dockerAccess := &SSHKeysAlgo{}
	expectedUUID := "ef69f752-0e89-46e2-a644-310429ae5f45"
	if dockerAccess.UUID() != expectedUUID {
		t.Errorf("Expected UUID %s, got %s", expectedUUID, dockerAccess.UUID())
	}
}

func TestSSHKeysAlgo_Passed(t *testing.T) {
	dockerAccess := &SSHKeysAlgo{passed: true}
	expectedPassed := true
	if dockerAccess.Passed() != expectedPassed {
		t.Errorf("Expected Passed %v, got %v", expectedPassed, dockerAccess.Passed())
	}
}

func TestSSHKeysAlgo_FailedMessage(t *testing.T) {
	dockerAccess := &SSHKeysAlgo{}
	expectedFailedMessage := "SSH keys are using weak encryption"
	if dockerAccess.FailedMessage() != expectedFailedMessage {
		t.Errorf("Expected FailedMessage %s, got %s", expectedFailedMessage, dockerAccess.FailedMessage())
	}
}

func TestSSHKeysAlgo_PassedMessage(t *testing.T) {
	dockerAccess := &SSHKeysAlgo{}
	expectedPassedMessage := "SSH keys use strong encryption"
	if dockerAccess.PassedMessage() != expectedPassedMessage {
		t.Errorf("Expected PassedMessage %s, got %s", expectedPassedMessage, dockerAccess.PassedMessage())
	}
}
