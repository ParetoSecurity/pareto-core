package checks

import (
	"testing"

	sharedchecks "github.com/ParetoSecurity/agent/checks/shared"
	"github.com/stretchr/testify/assert"
)

func TestSharing_Run(t *testing.T) {
	tests := []struct {
		name     string
		mockFunc func(port int, proto string) bool
		expected bool
	}{
		{
			name: "No ports open",
			mockFunc: func(port int, proto string) bool {
				return false
			},
			expected: true,
		},
		{
			name: "Some ports open",
			mockFunc: func(port int, proto string) bool {
				if port == 445 || port == 2049 {
					return true
				}
				return false
			},
			expected: false,
		},
		{
			name: "All ports open",
			mockFunc: func(port int, proto string) bool {
				return true
			},
			expected: false,
		},
		{
			name: "Only one port open",
			mockFunc: func(port int, proto string) bool {
				return port == 445
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the checkPort function
			sharedchecks.CheckPortMock = tt.mockFunc

			sharing := &Sharing{}
			err := sharing.Run()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, sharing.Passed())
			assert.NotEmpty(t, sharing.UUID())
			assert.False(t, sharing.RequiresRoot())
		})
	}
}

func TestSharing_Name(t *testing.T) {
	sharing := &Sharing{}
	expectedName := "File Sharing is disabled"
	if sharing.Name() != expectedName {
		t.Errorf("Expected Name %s, got %s", expectedName, sharing.Name())
	}
}

func TestSharing_Status(t *testing.T) {
	sharing := &Sharing{}
	expectedStatus := "Sharing services found running on ports:"
	if sharing.Status() != expectedStatus {
		t.Errorf("Expected Status %s, got %s", expectedStatus, sharing.Status())
	}
}

func TestSharing_UUID(t *testing.T) {
	sharing := &Sharing{}
	expectedUUID := "b96524e0-850b-4bb8-abc7-517051b6c14e"
	if sharing.UUID() != expectedUUID {
		t.Errorf("Expected UUID %s, got %s", expectedUUID, sharing.UUID())
	}
}

func TestSharing_Passed(t *testing.T) {
	sharing := &Sharing{passed: true}
	expectedPassed := true
	if sharing.Passed() != expectedPassed {
		t.Errorf("Expected Passed %v, got %v", expectedPassed, sharing.Passed())
	}
}

func TestSharing_FailedMessage(t *testing.T) {
	sharing := &Sharing{}
	expectedFailedMessage := "Sharing services found running "
	if sharing.FailedMessage() != expectedFailedMessage {
		t.Errorf("Expected FailedMessage %s, got %s", expectedFailedMessage, sharing.FailedMessage())
	}
}

func TestSharing_PassedMessage(t *testing.T) {
	sharing := &Sharing{}
	expectedPassedMessage := "No file sharing services found running"
	if sharing.PassedMessage() != expectedPassedMessage {
		t.Errorf("Expected PassedMessage %s, got %s", expectedPassedMessage, sharing.PassedMessage())
	}
}
