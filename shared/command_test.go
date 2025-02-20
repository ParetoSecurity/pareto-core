package shared

import (
	"strings"
	"testing"
)

func TestRunCommandFixtureSuccess(t *testing.T) {
	// Setup mock output for the command "echo hello"
	RunCommandMocks = make(map[string]string)
	key := "echo hello"
	expectedOutput := "hello\n"
	RunCommandMocks[key] = expectedOutput

	output, err := RunCommand("echo", "hello")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if output != expectedOutput {
		t.Errorf("expected output %q, got %q", expectedOutput, output)
	}
}

func TestRunCommandFixtureNotFound(t *testing.T) {
	// Setup mocks without the fixture for "nonexistent command"
	RunCommandMocks = make(map[string]string)

	output, err := RunCommand("nonexistent", "command")
	if err == nil {
		t.Fatalf("expected error, got nil with output %q", output)
	}
	if !strings.Contains(err.Error(), "RunCommand fixture not found") {
		t.Errorf("expected error to contain %q, got %v", "RunCommand fixture not found", err)
	}
}
