package shared

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml"
)

func TestCommitLastState_Success(t *testing.T) {
	// Create a temporary directory for our test file.
	tmpDir, err := os.MkdirTemp("", "commitlaststate_success")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Override statePath to a file in the temporary directory.
	testFile := filepath.Join(tmpDir, "test.state")
	statePath = testFile

	// Prepare a test state.
	testState := LastState{
		UUID:    "test-uuid",
		State:   true,
		Details: "all good",
	}

	// Clear and set the states map for a clean test.
	mutex.Lock()
	states = make(map[string]LastState)
	states[testState.UUID] = testState
	mutex.Unlock()

	// Commit the state to the file.
	if err := CommitLastState(); err != nil {
		t.Fatalf("CommitLastState failed: %v", err)
	}

	// Open the file and decode its contents.
	file, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	var decoded map[string]LastState
	decoder := toml.NewDecoder(file)
	if err := decoder.Decode(&decoded); err != nil {
		t.Fatalf("failed to decode TOML file: %v", err)
	}

	// Validate that the decoded state matches the test state.
	got, exists := decoded[testState.UUID]
	if !exists {
		t.Fatalf("expected state with UUID %s not found", testState.UUID)
	}
	if got != testState {
		t.Fatalf("expected state %+v, got %+v", testState, got)
	}
}

func TestCommitLastState_Error(t *testing.T) {
	// Simulate an error by setting statePath to a directory path.
	tmpDir, err := os.MkdirTemp("", "commitlaststate_error")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	statePath = tmpDir // os.Create on a directory should fail

	// Clear the states map.
	mutex.Lock()
	states = make(map[string]LastState)
	mutex.Unlock()

	// Attempt to commit; it should return an error.
	if err := CommitLastState(); err == nil {
		t.Fatalf("expected error when committing to a directory, got none")
	}
}
