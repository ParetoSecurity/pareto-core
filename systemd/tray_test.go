package systemd

import (
	"testing"

	"github.com/ParetoSecurity/agent/shared"
	"github.com/stretchr/testify/assert"
)

func TestIsTrayIconEnabled(t *testing.T) {
	tests := []struct {
		name     string
		mock     shared.RunCommandMock
		expected bool
	}{
		{
			name: "service is enabled",
			mock: shared.RunCommandMock{
				Command: "systemctl",
				Args:    []string{"--user", "is-enabled", "paretosecurity-trayicon.service"},
				Out:     "enabled\n",
				Err:     nil,
			},
			expected: true,
		},
		{
			name: "service is disabled",
			mock: shared.RunCommandMock{
				Command: "systemctl",
				Args:    []string{"--user", "is-enabled", "paretosecurity-trayicon.service"},
				Out:     "disabled\n",
				Err:     nil,
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			shared.RunCommandMocks = []shared.RunCommandMock{tc.mock}
			defer func() { shared.RunCommandMocks = nil }()

			result := IsTrayIconEnabled()
			assert.Equal(t, tc.expected, result)
		})
	}
}
func TestEnableTrayIcon(t *testing.T) {
	tests := []struct {
		name    string
		mock    shared.RunCommandMock
		wantErr bool
	}{
		{
			name: "enable succeeds",
			mock: shared.RunCommandMock{
				Command: "systemctl",
				Args:    []string{"--user", "enable", "paretosecurity-trayicon.service"},
				Out:     "",
				Err:     nil,
			},
			wantErr: false,
		},
		{
			name: "enable fails",
			mock: shared.RunCommandMock{
				Command: "systemctl",
				Args:    []string{"--user", "enable", "paretosecurity-trayicon.service"},
				Out:     "",
				Err:     assert.AnError,
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			shared.RunCommandMocks = []shared.RunCommandMock{tc.mock}
			defer func() { shared.RunCommandMocks = nil }()

			err := EnableTrayIcon()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestDisableTrayIcon(t *testing.T) {
	tests := []struct {
		name    string
		mock    shared.RunCommandMock
		wantErr bool
	}{
		{
			name: "disable succeeds",
			mock: shared.RunCommandMock{
				Command: "systemctl",
				Args:    []string{"--user", "disable", "paretosecurity-trayicon.service"},
				Out:     "",
				Err:     nil,
			},
			wantErr: false,
		},
		{
			name: "disable fails",
			mock: shared.RunCommandMock{
				Command: "systemctl",
				Args:    []string{"--user", "disable", "paretosecurity-trayicon.service"},
				Out:     "",
				Err:     assert.AnError,
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			shared.RunCommandMocks = []shared.RunCommandMock{tc.mock}
			defer func() { shared.RunCommandMocks = nil }()

			err := DisableTrayIcon()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
