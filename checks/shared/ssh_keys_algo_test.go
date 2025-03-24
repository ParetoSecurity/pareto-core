package shared

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"os"
	"testing"

	"golang.org/x/crypto/ssh"
)

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
func TestSSHKeysAlgo_isKeyStrong(t *testing.T) {
	// Override osReadFile for testing

	tests := []struct {
		name     string
		keyData  string
		expected bool
	}{
		{
			name:     "RSA 1024 bit key (weak)",
			keyData:  generateRealKey(t, "rsa", 1024),
			expected: false,
		},
		{
			name:     "DSS key (weak)",
			keyData:  "ssh-dss AAAAB3NzaC1kc3MAAACBAJxTK4HhLxW+v3uYv+RBS1L3seXnbU1alYGXjCJ4dmpyi1IzZ2pHGnNLXhjb/JMXwpZ0Fp+fZaRfLonnLq1xwULZLvlL0bhbmaj7VdwTV7yD5JC4CsdLMbrIXl9Q93OWJIDRlcDFW/CmpWz2QMIymUP7GOtGWB+aZxQ/0+Um3SLHAAAAFQCt56ak5Xt1llWeORQ9EFZJjwpUpQAAAIB4YkEr7r/6P9zzDmSYVMgH9PRwsRVcCZWGYpJi9xUdmghxqI7qZtSXBYOOB7QTtm9H3xvj0E9DHCkJINXr3C1MZNdpm6Jq0Wjg8XOJJCDWvhJS1ITTEZCrXgIPYHXXn9Q0yO0c5i0zLIUTHCMgPEb+uK5lSQcqWEjQ0Z7SeIleSQAAAIBGQJ10MDXi8+TIAUODz/4lGT5J3A0H5jzTb3ez8vlYQ0zRuxe3uFUeqt6d7CqBYmjSVuFKu0tMFQpFsGP2JlQ8t1YNm4/FKdRp7MllCH6CJGPk+IXeysWNz+a9bQF7A5+OiL3xSttIOy4kD8//F+B02nHaP3mTFxBvAICjdpyhQA== test@example.com",
			expected: false,
		},
		{
			name:     "Unknown key type",
			keyData:  "ssh-unknown AAAAC3NzaC1lZDI1NTE5AAAAIFoZWbQD4XFRnOGZt8YuoN26+OMb4YKHbMDm0/lAWJxz test@example.com",
			expected: false,
		},
		{
			name:     "Invalid key format",
			keyData:  "invalid-key-format",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			osReadFileMock = func(path string) ([]byte, error) {
				return []byte(tt.keyData), nil
			}

			sshCheck := &SSHKeysAlgo{}
			result := sshCheck.isKeyStrong("dummy/path")
			if result != tt.expected {
				t.Errorf("isKeyStrong() = %v, want %v", result, tt.expected)
			}
		})
	}

	// Test file read error
	t.Run("File read error", func(t *testing.T) {
		osReadFileMock = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		sshCheck := &SSHKeysAlgo{}
		result := sshCheck.isKeyStrong("dummy/path")
		if result != false {
			t.Errorf("isKeyStrong() = %v, want %v", result, false)
		}
	})
}

// generateRealKey generates real SSH public keys for testing
func generateRealKey(t *testing.T, keyType string, bits int) string {
	t.Helper()

	var pubKey ssh.PublicKey

	switch keyType {
	case "rsa":
		privateKey, err := rsa.GenerateKey(rand.Reader, bits)
		if err != nil {
			t.Fatalf("Failed to generate RSA key: %v", err)
		}

		pubKey, err = ssh.NewPublicKey(&privateKey.PublicKey)
		if err != nil {
			t.Fatalf("Failed to convert RSA key to SSH format: %v", err)
		}

	case "ecdsa":
		var curve elliptic.Curve
		switch bits {
		case 256:
			curve = elliptic.P256()
		case 384:
			curve = elliptic.P384()
		case 521:
			curve = elliptic.P521()
		default:
			curve = elliptic.P256()
		}

		privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			t.Fatalf("Failed to generate ECDSA key: %v", err)
		}

		pubKey, err = ssh.NewPublicKey(&privateKey.PublicKey)
		if err != nil {
			t.Fatalf("Failed to convert ECDSA key to SSH format: %v", err)
		}

	case "ed25519":
		publicKey, _, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			t.Fatalf("Failed to generate Ed25519 key: %v", err)
		}

		pubKey, err = ssh.NewPublicKey(publicKey)
		if err != nil {
			t.Fatalf("Failed to convert Ed25519 key to SSH format: %v", err)
		}
	}

	return string(ssh.MarshalAuthorizedKey(pubKey))
}
