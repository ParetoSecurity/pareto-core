package shared

import (
	"testing"

	"github.com/ParetoSecurity/agent/shared"
	"github.com/h2non/gock"
)

func TestParetoUpdated_Run(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	shared.Config.TeamID = "test-team-id"
	shared.Config.AuthToken = "test-auth-token"
	defer func() {
		shared.Config.TeamID = ""
		shared.Config.AuthToken = ""
	}()

	gock.New("https://paretosecurity.com/api/updates").
		Reply(200).
		JSON([]map[string]string{{"tag_name": "1.7.91"}})
	check := &ParetoUpdated{}
	err := check.Run()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

}

func TestParetoUpdated_RunPublic(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	shared.Config.TeamID = ""
	shared.Config.AuthToken = ""

	gock.New("https://api.github.com").
		Reply(200).
		JSON([]map[string]string{{"tag_name": "1.7.91"}})
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
	expectedStatus := "Pareto Security is outdated "
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
	expectedFailedMessage := "Pareto Security is outdated "
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
