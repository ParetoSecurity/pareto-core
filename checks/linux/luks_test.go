package checks

import (
	"testing"

	"github.com/ParetoSecurity/agent/shared"
	"github.com/stretchr/testify/assert"
)

func TestMaybeCryptoViaKernel(t *testing.T) {
	tests := []struct {
		name     string
		cmdline  string
		expected bool
	}{
		{
			name:     "cryptdevice present and root",
			cmdline:  "BOOT_IMAGE=/vmlinuz-linux cryptdevice=UUID=1234-5678-90AB-CDEF:cryptroot:root root=/dev/mapper/cryptroot",
			expected: true,
		},
		{
			name:     "cryptdevice present but not root",
			cmdline:  "BOOT_IMAGE=/vmlinuz-linux cryptdevice=UUID=1234-5678-90AB-CDEF:cryptroot:other root=/dev/mapper/cryptroot",
			expected: false,
		},
		{
			name:     "cryptdevice not present",
			cmdline:  "BOOT_IMAGE=/vmlinuz-linux root=/dev/sda1",
			expected: false,
		},
		{
			name:     "empty cmdline",
			cmdline:  "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the shared.ReadFile function
			shared.ReadFileMocks = map[string]string{
				"/proc/cmdline": tt.cmdline,
			}

			result := maybeCryptoViaKernel()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEncryptingFS_Name(t *testing.T) {
	e := &EncryptingFS{}
	expectedName := "Filesystem encryption is enabled"
	if e.Name() != expectedName {
		t.Errorf("Expected Name %s, got %s", expectedName, e.Name())
	}
}

func TestEncryptingFS_Status(t *testing.T) {
	e := &EncryptingFS{}
	expectedStatus := ""
	if e.Status() != expectedStatus {
		t.Errorf("Expected Status %s, got %s", expectedStatus, e.Status())
	}
}

func TestEncryptingFS_UUID(t *testing.T) {
	e := &EncryptingFS{}
	expectedUUID := "21830a4e-84f1-48fe-9c5b-beab436b2cdb"
	if e.UUID() != expectedUUID {
		t.Errorf("Expected UUID %s, got %s", expectedUUID, e.UUID())
	}
}

func TestEncryptingFS_Passed(t *testing.T) {
	e := &EncryptingFS{passed: true}
	expectedPassed := true
	if e.Passed() != expectedPassed {
		t.Errorf("Expected Passed %v, got %v", expectedPassed, e.Passed())
	}
}

func TestEncryptingFS_FailedMessage(t *testing.T) {
	e := &EncryptingFS{}
	expectedFailedMessage := "Block device encryption is disabled"
	if e.FailedMessage() != expectedFailedMessage {
		t.Errorf("Expected FailedMessage %s, got %s", expectedFailedMessage, e.FailedMessage())
	}
}

func TestEncryptingFS_PassedMessage(t *testing.T) {
	e := &EncryptingFS{}
	expectedPassedMessage := "Block device encryption is enabled"
	if e.PassedMessage() != expectedPassedMessage {
		t.Errorf("Expected PassedMessage %s, got %s", expectedPassedMessage, e.PassedMessage())
	}
}
