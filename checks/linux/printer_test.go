package checks

import (
	"testing"

	sharedchecks "github.com/ParetoSecurity/agent/checks/shared"
	"github.com/stretchr/testify/assert"
)

func TestPrinterRun(t *testing.T) {
	tests := []struct {
		name           string
		mockCheckPort  func(_ int, _ string) bool
		expectedPassed bool
		expectedPorts  map[int]string
	}{
		{
			name: "No ports open",
			mockCheckPort: func(_ int, _ string) bool {
				return false
			},
			expectedPassed: true,
			expectedPorts:  map[int]string{},
		},
		{
			name: "CUPS port open",
			mockCheckPort: func(port int, proto string) bool {
				return port == 631
			},
			expectedPassed: false,
			expectedPorts: map[int]string{
				631: "CUPS",
			},
		},
		{
			name: "Multiple ports open",
			mockCheckPort: func(port int, proto string) bool {
				return port == 631 || port == 515
			},
			expectedPassed: false,
			expectedPorts: map[int]string{
				631: "CUPS",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sharedchecks.CheckPortMock = tt.mockCheckPort
			printer := &Printer{}

			err := printer.Run()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPassed, printer.Passed())
			assert.Equal(t, tt.expectedPorts, printer.ports)
			assert.NotEmpty(t, printer.UUID())
			assert.False(t, printer.RequiresRoot())
		})
	}
}

func TestPrinter_Name(t *testing.T) {
	printer := &Printer{}
	expectedName := "Sharing printers is off"
	if printer.Name() != expectedName {
		t.Errorf("Expected Name %s, got %s", expectedName, printer.Name())
	}
}

func TestPrinter_Status(t *testing.T) {
	printer := &Printer{}
	expectedStatus := "Printer sharing services found running on ports:"
	if printer.Status() != expectedStatus {
		t.Errorf("Expected Status %s, got %s", expectedStatus, printer.Status())
	}
}

func TestPrinter_UUID(t *testing.T) {
	printer := &Printer{}
	expectedUUID := "b96524e0-150b-4bb8-abc7-517051b6c14e"
	if printer.UUID() != expectedUUID {
		t.Errorf("Expected UUID %s, got %s", expectedUUID, printer.UUID())
	}
}

func TestPrinter_Passed(t *testing.T) {
	printer := &Printer{passed: true}
	expectedPassed := true
	if printer.Passed() != expectedPassed {
		t.Errorf("Expected Passed %v, got %v", expectedPassed, printer.Passed())
	}
}

func TestPrinter_FailedMessage(t *testing.T) {
	printer := &Printer{}
	expectedFailedMessage := "Sharing printers is on"
	if printer.FailedMessage() != expectedFailedMessage {
		t.Errorf("Expected FailedMessage %s, got %s", expectedFailedMessage, printer.FailedMessage())
	}
}

func TestPrinter_PassedMessage(t *testing.T) {
	printer := &Printer{}
	expectedPassedMessage := "Sharing printers is off"
	if printer.PassedMessage() != expectedPassedMessage {
		t.Errorf("Expected PassedMessage %s, got %s", expectedPassedMessage, printer.PassedMessage())
	}
}
