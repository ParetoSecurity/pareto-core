package checks

import (
	"testing"

	"github.com/ParetoSecurity/agent/shared"
	"github.com/stretchr/testify/assert"
)

func TestCheckUFW(t *testing.T) {
	tests := []struct {
		name           string
		mockOutput     string
		mockError      error
		expectedResult bool
	}{
		{
			name:           "UFW is active",
			mockOutput:     "Status: active",
			mockError:      nil,
			expectedResult: true,
		},
		{
			name:           "UFW is inactive",
			mockOutput:     "Status: inactive",
			mockError:      nil,
			expectedResult: false,
		},
		{
			name:           "UFW command error",
			mockOutput:     "",
			mockError:      assert.AnError,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.RunCommandMocks = convertCommandMapToMocks(map[string]string{
				"ufw status": tt.mockOutput,
			})
			f := &Firewall{}
			result := f.checkUFW()
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestCheckFirewalld(t *testing.T) {
	tests := []struct {
		name           string
		mockOutput     string
		mockError      error
		expectedResult bool
	}{
		{
			name:           "Firewalld is active",
			mockOutput:     "active",
			mockError:      nil,
			expectedResult: true,
		},
		{
			name:           "Firewalld is inactive",
			mockOutput:     "inactive",
			mockError:      nil,
			expectedResult: false,
		},
		{
			name:           "Firewalld command error",
			mockOutput:     "",
			mockError:      assert.AnError,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.RunCommandMocks = convertCommandMapToMocks(map[string]string{
				"systemctl is-active firewalld": tt.mockOutput,
			})

			f := &Firewall{}
			result := f.checkFirewalld()
			assert.Equal(t, tt.expectedResult, result)
			assert.NotEmpty(t, f.UUID())
			assert.True(t, f.RequiresRoot())
		})
	}
}

func TestFirewall_Run(t *testing.T) {
	tests := []struct {
		name                string
		mockUFWOutput       string
		mockFirewalldOutput string
		expectedPassed      bool
		expectedStatus      string
	}{
		{
			name:                "UFW is active",
			mockUFWOutput:       "Status: active",
			mockFirewalldOutput: "",
			expectedPassed:      true,
			expectedStatus:      "Firewall is on",
		},
		{
			name:                "Firewalld is active",
			mockUFWOutput:       "Status: inactive",
			mockFirewalldOutput: "active",
			expectedPassed:      true,
			expectedStatus:      "Firewall is on",
		},
		{
			name:                "Both UFW and Firewalld are inactive",
			mockUFWOutput:       "Status: inactive",
			mockFirewalldOutput: "inactive",
			expectedPassed:      false,
			expectedStatus:      "Firewall is off",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.RunCommandMocks = convertCommandMapToMocks(map[string]string{
				"ufw status":                    tt.mockUFWOutput,
				"systemctl is-active firewalld": tt.mockFirewalldOutput,
			})

			f := &Firewall{}
			err := f.Run()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPassed, f.Passed())
			assert.Equal(t, tt.expectedStatus, f.Status())
		})
	}
}

func TestFirewall_Name(t *testing.T) {
	f := &Firewall{}
	expectedName := "Firewall is on"
	if f.Name() != expectedName {
		t.Errorf("Expected Name %s, got %s", expectedName, f.Name())
	}
}

func TestFirewall_Status(t *testing.T) {
	f := &Firewall{}
	expectedStatus := "Firewall is off"
	if f.Status() != expectedStatus {
		t.Errorf("Expected Status %s, got %s", expectedStatus, f.Status())
	}
}

func TestFirewall_UUID(t *testing.T) {
	f := &Firewall{}
	expectedUUID := "2e46c89a-5461-4865-a92e-3b799c12034a"
	if f.UUID() != expectedUUID {
		t.Errorf("Expected UUID %s, got %s", expectedUUID, f.UUID())
	}
}

func TestFirewall_Passed(t *testing.T) {
	f := &Firewall{passed: true}
	expectedPassed := true
	if f.Passed() != expectedPassed {
		t.Errorf("Expected Passed %v, got %v", expectedPassed, f.Passed())
	}
}

func TestFirewall_FailedMessage(t *testing.T) {
	f := &Firewall{}
	expectedFailedMessage := "Firewall is off"
	if f.FailedMessage() != expectedFailedMessage {
		t.Errorf("Expected FailedMessage %s, got %s", expectedFailedMessage, f.FailedMessage())
	}
}

func TestFirewall_PassedMessage(t *testing.T) {
	f := &Firewall{}
	expectedPassedMessage := "Firewall is on"
	if f.PassedMessage() != expectedPassedMessage {
		t.Errorf("Expected PassedMessage %s, got %s", expectedPassedMessage, f.PassedMessage())
	}
}

func TestCheckIptables(t *testing.T) {
	tests := []struct {
		name           string
		mockOutput     string
		mockError      error
		expectedResult bool
	}{
		{
			name: "Iptables has rules",
			mockOutput: `Chain INPUT (policy ACCEPT)
num  target     prot opt source               destination         
1    ACCEPT     tcp  --  0.0.0.0/0            0.0.0.0/0           tcp dpt:22
2    DROP       all  --  10.0.0.0/8           0.0.0.0/0           
`,
			mockError:      nil,
			expectedResult: true,
		},
		{
			name: "Iptables has no rules",
			mockOutput: `Chain INPUT (policy ACCEPT)
num  target     prot opt source               destination         
`,
			mockError:      nil,
			expectedResult: false,
		},
		{
			name:           "Iptables command error",
			mockOutput:     "",
			mockError:      assert.AnError,
			expectedResult: false,
		},
		{
			name: "Malformed rule line",
			mockOutput: `Chain INPUT (policy ACCEPT)
num  target     prot opt source               destination         
invalid line
`,
			mockError:      nil,
			expectedResult: false,
		},
		{
			name: "Non-numeric rule number",
			mockOutput: `Chain INPUT (policy ACCEPT)
num  target     prot opt source               destination         
abc  ACCEPT     tcp  --  0.0.0.0/0            0.0.0.0/0           
`,
			mockError:      nil,
			expectedResult: false,
		},
		{
			name: "NixOS style custom chain",
			mockOutput: `Chain INPUT (policy ACCEPT)
num  target     prot opt source               destination         
1    nixos-fw   all  --  anywhere             anywhere            
`,
			mockError:      nil,
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.RunCommandMocks = convertCommandMapToMocks(map[string]string{
				"iptables -L INPUT --line-numbers": tt.mockOutput,
			})
			f := &Firewall{}
			result := f.checkIptables()
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestFirewall_fwCmdsAreAvailable(t *testing.T) {

	tests := []struct {
		name           string
		mockLookPath   func(string) (string, error)
		expectedResult bool
		expectedStatus string
	}{
		{
			name: "All firewall commands are available",
			mockLookPath: func(cmd string) (string, error) {
				return "/usr/bin/" + cmd, nil
			},
			expectedResult: true,
			expectedStatus: "",
		},
		{
			name: "Only UFW is available",
			mockLookPath: func(cmd string) (string, error) {
				if cmd == "ufw" {
					return "/usr/bin/ufw", nil
				}
				return "", assert.AnError
			},
			expectedResult: true,
			expectedStatus: "",
		},
		{
			name: "Only firewalld is available",
			mockLookPath: func(cmd string) (string, error) {
				if cmd == "firewalld" {
					return "/usr/bin/firewalld", nil
				}
				return "", assert.AnError
			},
			expectedResult: true,
			expectedStatus: "",
		},
		{
			name: "Only iptables is available",
			mockLookPath: func(cmd string) (string, error) {
				if cmd == "iptables" {
					return "/usr/bin/iptables", nil
				}
				return "", assert.AnError
			},
			expectedResult: true,
			expectedStatus: "",
		},
		{
			name: "No firewall commands are available",
			mockLookPath: func(cmd string) (string, error) {
				return "", assert.AnError
			},
			expectedResult: false,
			expectedStatus: "Neither ufw, firewalld nor iptables are present, check cannot run",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lookPathMock = tt.mockLookPath
			f := &Firewall{}
			result := f.fwCmdsAreAvailable()
			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedStatus, f.status)
		})
	}
}

func TestFirewall_Run_NoFirewallCommands(t *testing.T) {
	f := &Firewall{
		status: "Neither ufw, firewalld nor iptables are present, check cannot run",
		passed: false,
	}

	err := f.Run()
	assert.NoError(t, err)
	assert.False(t, f.Passed())
	assert.Equal(t, "Neither ufw, firewalld nor iptables are present, check cannot run", f.status)
}
