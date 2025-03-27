package systemd

import (
	"errors"
	"testing"

	"github.com/ParetoSecurity/agent/shared"
)

func TestIsTimerEnabled(t *testing.T) {
	tests := []struct {
		name     string
		mocks    []shared.RunCommandMock
		expected bool
	}{
		{
			name: "both services enabled",
			mocks: []shared.RunCommandMock{
				{
					Command: "systemctl",
					Out:     "enabled",
					Args:    []string{"--user", "is-enabled", "paretosecurity-user.timer"},
					Err:     nil,
				},
				{
					Command: "systemctl",
					Out:     "enabled",
					Args:    []string{"--user", "is-enabled", "paretosecurity-user.service"},
					Err:     nil,
				},
			},
			expected: true,
		},
		{
			name: "timer disabled",
			mocks: []shared.RunCommandMock{
				{
					Command: "systemctl",
					Out:     "disabled",
					Args:    []string{"--user", "is-enabled", "paretosecurity-user.timer"},
					Err:     nil,
				},
				{
					Command: "systemctl",
					Out:     "enabled",
					Args:    []string{"--user", "is-enabled", "paretosecurity-user.service"},
					Err:     nil,
				},
			},
			expected: false,
		},
		{
			name: "service disabled",
			mocks: []shared.RunCommandMock{
				{
					Command: "systemctl",
					Out:     "enabled",
					Args:    []string{"--user", "is-enabled", "paretosecurity-user.timer"},
					Err:     nil,
				},
				{
					Command: "systemctl",
					Out:     "disabled",
					Args:    []string{"--user", "is-enabled", "paretosecurity-user.service"},
					Err:     nil,
				},
			},
			expected: false,
		},
		{
			name: "both services disabled",
			mocks: []shared.RunCommandMock{
				{
					Command: "systemctl",
					Out:     "disabled",
					Args:    []string{"--user", "is-enabled", "paretosecurity-user.timer"},
					Err:     nil,
				},
				{
					Command: "systemctl",
					Out:     "disabled",
					Args:    []string{"--user", "is-enabled", "paretosecurity-user.service"},
					Err:     nil,
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			shared.RunCommandMocks = tt.mocks

			// Run test
			result := IsTimerEnabled()

			// Check result
			if result != tt.expected {
				t.Errorf("IsTimerEnabled() = %v, want %v", result, tt.expected)
			}
		})
	}
}
func TestEnableTimer(t *testing.T) {
	tests := []struct {
		name          string
		mocks         []shared.RunCommandMock
		expectedError bool
	}{
		{
			name: "successfully enable both",
			mocks: []shared.RunCommandMock{
				{
					Command: "systemctl",
					Out:     "",
					Args:    []string{"--user", "enable", "paretosecurity-user.timer"},
					Err:     nil,
				},
				{
					Command: "systemctl",
					Out:     "",
					Args:    []string{"--user", "enable", "paretosecurity-user.service"},
					Err:     nil,
				},
			},
			expectedError: false,
		},
		{
			name: "error enabling timer",
			mocks: []shared.RunCommandMock{
				{
					Command: "systemctl",
					Out:     "",
					Args:    []string{"--user", "enable", "paretosecurity-user.timer"},
					Err:     errors.New("failed to enable timer"),
				},
			},
			expectedError: true,
		},
		{
			name: "error enabling service",
			mocks: []shared.RunCommandMock{
				{
					Command: "systemctl",
					Out:     "",
					Args:    []string{"--user", "enable", "paretosecurity-user.timer"},
					Err:     nil,
				},
				{
					Command: "systemctl",
					Out:     "",
					Args:    []string{"--user", "enable", "paretosecurity-user.service"},
					Err:     errors.New("failed to enable service"),
				},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			shared.RunCommandMocks = tt.mocks

			// Run test
			err := EnableTimer()

			// Check result
			if (err != nil) != tt.expectedError {
				t.Errorf("EnableTimer() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestDisableTimer(t *testing.T) {
	tests := []struct {
		name          string
		mocks         []shared.RunCommandMock
		expectedError bool
	}{
		{
			name: "successfully disable both",
			mocks: []shared.RunCommandMock{
				{
					Command: "systemctl",
					Out:     "",
					Args:    []string{"--user", "disable", "paretosecurity-user.timer"},
					Err:     nil,
				},
				{
					Command: "systemctl",
					Out:     "",
					Args:    []string{"--user", "disable", "paretosecurity-user.service"},
					Err:     nil,
				},
			},
			expectedError: false,
		},
		{
			name: "error disabling timer",
			mocks: []shared.RunCommandMock{
				{
					Command: "systemctl",
					Out:     "",
					Args:    []string{"--user", "disable", "paretosecurity-user.timer"},
					Err:     errors.New("failed to disable timer"),
				},
			},
			expectedError: true,
		},
		{
			name: "error disabling service",
			mocks: []shared.RunCommandMock{
				{
					Command: "systemctl",
					Out:     "",
					Args:    []string{"--user", "disable", "paretosecurity-user.timer"},
					Err:     nil,
				},
				{
					Command: "systemctl",
					Out:     "",
					Args:    []string{"--user", "disable", "paretosecurity-user.service"},
					Err:     errors.New("failed to disable service"),
				},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			shared.RunCommandMocks = tt.mocks

			// Run test
			err := DisableTimer()

			// Check result
			if (err != nil) != tt.expectedError {
				t.Errorf("DisableTimer() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}
