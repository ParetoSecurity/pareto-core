package check

import (
	"testing"
	"time"

	"github.com/ParetoSecurity/agent/shared"
)

type MockCheck struct {
	uuid       string
	passed     bool
	isRunnable bool
}

func (m *MockCheck) Name() string          { return "MockCheck" }
func (m *MockCheck) PassedMessage() string { return "Passed" }
func (m *MockCheck) FailedMessage() string { return "Failed" }
func (m *MockCheck) Run() error            { return nil }
func (m *MockCheck) Passed() bool          { return m.passed }
func (m *MockCheck) IsRunnable() bool      { return m.isRunnable }
func (m *MockCheck) UUID() string          { return m.uuid }
func (m *MockCheck) Status() string        { return "Status" }
func (m *MockCheck) RequiresRoot() bool    { return false }

func TestRegister(t *testing.T) {

	mockCheck := &MockCheck{
		uuid:       "1234",
		passed:     true,
		isRunnable: true,
	}

	registeredCheck := Register(mockCheck)

	if registeredCheck.UUID() != mockCheck.UUID() {
		t.Errorf("Expected UUID %s, got %s", mockCheck.UUID(), registeredCheck.UUID())
	}

	if len(shared.Config.Checks) != 1 {
		t.Errorf("Expected 1 check in the map, got %d", len(shared.Config.Checks))
	}

	checkStatus, exists := shared.Config.Checks[mockCheck.UUID()]
	if !exists {
		t.Errorf("Check with UUID %s not found in the map", mockCheck.UUID())
	}

	if checkStatus.Passed != mockCheck.Passed() {
		t.Errorf("Expected Passed %v, got %v", mockCheck.Passed(), checkStatus.Passed)
	}

	if checkStatus.Disabled != !mockCheck.IsRunnable() {
		t.Errorf("Expected Disabled %v, got %v", !mockCheck.IsRunnable(), checkStatus.Disabled)
	}

	if time.Since(checkStatus.UpdatedAt) > time.Second {
		t.Errorf("Expected UpdatedAt to be recent, got %v", checkStatus.UpdatedAt)
	}
}

func TestMockCheck_Name(t *testing.T) {
	mockCheck := &MockCheck{}
	expectedName := "MockCheck"
	if mockCheck.Name() != expectedName {
		t.Errorf("Expected Name %s, got %s", expectedName, mockCheck.Name())
	}
}

func TestMockCheck_Status(t *testing.T) {
	mockCheck := &MockCheck{}
	expectedStatus := "Status"
	if mockCheck.Status() != expectedStatus {
		t.Errorf("Expected Status %s, got %s", expectedStatus, mockCheck.Status())
	}
}

func TestMockCheck_UUID(t *testing.T) {
	mockCheck := &MockCheck{uuid: "1234"}
	expectedUUID := "1234"
	if mockCheck.UUID() != expectedUUID {
		t.Errorf("Expected UUID %s, got %s", expectedUUID, mockCheck.UUID())
	}
}

func TestMockCheck_Passed(t *testing.T) {
	mockCheck := &MockCheck{passed: true}
	expectedPassed := true
	if mockCheck.Passed() != expectedPassed {
		t.Errorf("Expected Passed %v, got %v", expectedPassed, mockCheck.Passed())
	}
}

func TestMockCheck_FailedMessage(t *testing.T) {
	mockCheck := &MockCheck{}
	expectedFailedMessage := "Failed"
	if mockCheck.FailedMessage() != expectedFailedMessage {
		t.Errorf("Expected FailedMessage %s, got %s", expectedFailedMessage, mockCheck.FailedMessage())
	}
}

func TestMockCheck_PassedMessage(t *testing.T) {
	mockCheck := &MockCheck{}
	expectedPassedMessage := "Passed"
	if mockCheck.PassedMessage() != expectedPassedMessage {
		t.Errorf("Expected PassedMessage %s, got %s", expectedPassedMessage, mockCheck.PassedMessage())
	}
}
