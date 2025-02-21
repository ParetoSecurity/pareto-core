package shared

import (
	"testing"
)

func TestParetoUpdated_Run(t *testing.T) {

	check := &ParetoUpdated{}
	err := check.Run()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

}

func TestParetoUpdated_Name(t *testing.T) {
	dockerAccess := &ParetoUpdated{}
	expectedName := "Pareto Security is up to date"
	if dockerAccess.Name() != expectedName {
		t.Errorf("Expected Name %s, got %s", expectedName, dockerAccess.Name())
	}
}

func TestParetoUpdated_Status(t *testing.T) {
	dockerAccess := &ParetoUpdated{}
	expectedStatus := "Pareto Security is oudated"
	if dockerAccess.Status() != expectedStatus {
		t.Errorf("Expected Status %s, got %s", expectedStatus, dockerAccess.Status())
	}
}

func TestParetoUpdated_UUID(t *testing.T) {
	dockerAccess := &ParetoUpdated{}
	expectedUUID := "44e4754a-0b42-4964-9cc2-b88b2023cb1e"
	if dockerAccess.UUID() != expectedUUID {
		t.Errorf("Expected UUID %s, got %s", expectedUUID, dockerAccess.UUID())
	}
}

func TestParetoUpdated_Passed(t *testing.T) {
	dockerAccess := &ParetoUpdated{passed: true}
	expectedPassed := true
	if dockerAccess.Passed() != expectedPassed {
		t.Errorf("Expected Passed %v, got %v", expectedPassed, dockerAccess.Passed())
	}
}

func TestParetoUpdated_FailedMessage(t *testing.T) {
	dockerAccess := &ParetoUpdated{}
	expectedFailedMessage := "Pareto Security is oudated"
	if dockerAccess.FailedMessage() != expectedFailedMessage {
		t.Errorf("Expected FailedMessage %s, got %s", expectedFailedMessage, dockerAccess.FailedMessage())
	}
}

func TestParetoUpdated_PassedMessage(t *testing.T) {
	dockerAccess := &ParetoUpdated{}
	expectedPassedMessage := "Pareto Security is up to date"
	if dockerAccess.PassedMessage() != expectedPassedMessage {
		t.Errorf("Expected PassedMessage %s, got %s", expectedPassedMessage, dockerAccess.PassedMessage())
	}
}
