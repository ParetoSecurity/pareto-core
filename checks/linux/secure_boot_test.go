package checks

import (
	"os"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestSecureBoot_Run(t *testing.T) {
	tests := []struct {
		name           string
		mockFiles      map[string][]byte
		expectedPassed bool
		expectedStatus string
	}{
		{
			name: "SecureBoot enabled",
			mockFiles: map[string][]byte{
				"/sys/firmware/efi/efivars/SecureBoot-1234": {0, 0, 0, 0, 1},
			},
			expectedPassed: true,
			expectedStatus: "SecureBoot is enabled",
		},
		{
			name: "SecureBoot disabled",
			mockFiles: map[string][]byte{
				"/sys/firmware/efi/efivars/SecureBoot-1234": {0, 0, 0, 0, 0},
			},
			expectedPassed: false,
			expectedStatus: "SecureBoot is disabled",
		},
		{
			name:           "SecureBoot EFI variable not found",
			mockFiles:      map[string][]byte{},
			expectedPassed: false,
			expectedStatus: "Could not find SecureBoot EFI variable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock os.ReadFile
			filepathGlobMock = func(_ string) ([]string, error) {
				return lo.Keys(tt.mockFiles), nil
			}
			osReadFileMock = func(file string) ([]byte, error) {
				return tt.mockFiles[file], nil
			}
			sb := &SecureBoot{}
			err := sb.Run()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPassed, sb.Passed())
			assert.Equal(t, tt.expectedStatus, sb.Status())
		})
	}
}

func TestSecureBoot_IsRunnable(t *testing.T) {
	tests := []struct {
		name           string
		mockStatError  error
		expectedResult bool
		expectedStatus string
	}{
		{
			name:           "System running in UEFI mode",
			mockStatError:  nil,
			expectedResult: true,
			expectedStatus: "System is not running in UEFI mode",
		},
		{
			name:           "System not running in UEFI mode",
			mockStatError:  os.ErrNotExist,
			expectedResult: false,
			expectedStatus: "System is not running in UEFI mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock os.Stat
			osStatMock = func(_ string) (os.FileInfo, error) {
				if tt.mockStatError != nil {
					return nil, tt.mockStatError
				}
				return nil, nil
			}
			sb := &SecureBoot{}
			result := sb.IsRunnable()
			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedStatus, sb.Status())
		})
	}
}

func TestSecureBoot_Name(t *testing.T) {
	sb := &SecureBoot{}
	expectedName := "SecureBoot is enabled"
	if sb.Name() != expectedName {
		t.Errorf("Expected Name %s, got %s", expectedName, sb.Name())
	}
}

func TestSecureBoot_Status(t *testing.T) {
	sb := &SecureBoot{}
	expectedStatus := "SecureBoot is disabled"
	if sb.Status() != expectedStatus {
		t.Errorf("Expected Status %s, got %s", expectedStatus, sb.Status())
	}
}

func TestSecureBoot_UUID(t *testing.T) {
	sb := &SecureBoot{}
	expectedUUID := "c96524f2-850b-4bb9-abc7-517051b6c14e"
	if sb.UUID() != expectedUUID {
		t.Errorf("Expected UUID %s, got %s", expectedUUID, sb.UUID())
	}
}

func TestSecureBoot_Passed(t *testing.T) {
	sb := &SecureBoot{passed: true}
	expectedPassed := true
	if sb.Passed() != expectedPassed {
		t.Errorf("Expected Passed %v, got %v", expectedPassed, sb.Passed())
	}
}

func TestSecureBoot_FailedMessage(t *testing.T) {
	sb := &SecureBoot{}
	expectedFailedMessage := "SecureBoot is disabled"
	if sb.FailedMessage() != expectedFailedMessage {
		t.Errorf("Expected FailedMessage %s, got %s", expectedFailedMessage, sb.FailedMessage())
	}
}

func TestSecureBoot_PassedMessage(t *testing.T) {
	sb := &SecureBoot{}
	expectedPassedMessage := "SecureBoot is enabled"
	if sb.PassedMessage() != expectedPassedMessage {
		t.Errorf("Expected PassedMessage %s, got %s", expectedPassedMessage, sb.PassedMessage())
	}
}
