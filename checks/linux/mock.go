package checks

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ParetoSecurity/agent/shared"
)

// lookPathMock is a mock function that simulates the behavior of
// the os/exec.LookPath function. It takes a file name as input
// and returns the path to the executable file along with an error
// if the file is not found or any other issue occurs.
var lookPathMock func(file string) (string, error)

func lookPath(file string) (string, error) {
	if testing.Testing() {
		return lookPathMock(file)
	}
	return exec.LookPath(file)
}

var osStatMock func(file string) (os.FileInfo, error)

// osStat checks if a file exists by attempting to get its file info.
// During testing, it uses a mock implementation via osStatMock.
// It returns the file path if the file exists, otherwise returns an empty string and error.
func osStat(file string) (os.FileInfo, error) {
	if testing.Testing() {
		return osStatMock(file)
	}
	return os.Stat(file)
}

var filepathGlobMock func(pattern string) ([]string, error)

// filepathGlob retrieves file paths that match the provided glob pattern.
//
// In a testing environment (when testing.Testing() returns true), it delegates
// the matching to filepathGlobMock to simulate the behavior. Otherwise, it uses
// the standard library's filepath.Glob to perform glob pattern matching.
func filepathGlob(pattern string) ([]string, error) {
	if testing.Testing() {
		return filepathGlobMock(pattern)
	}
	return filepath.Glob(pattern)
}

var osReadFileMock func(file string) ([]byte, error)

// osReadFile reads the contents of the specified file.
//
// If the testing mode is enabled, it delegates the file reading to a mock function.
// Otherwise, it reads the file from disk using the standard os.ReadFile function.
func osReadFile(file string) ([]byte, error) {
	if testing.Testing() {
		return osReadFileMock(file)
	}
	return os.ReadFile(file)
}

var osReadDirMock func(dirname string) ([]os.DirEntry, error)

// osReadDir reads the directory specified by dirname and returns a slice of os.DirEntry.
// In testing mode, it delegates to osReadDirMock for controlled behavior; otherwise,
// it uses os.ReadDir from the standard library.
func osReadDir(dirname string) ([]os.DirEntry, error) {
	if testing.Testing() {
		return osReadDirMock(dirname)
	}
	return os.ReadDir(dirname)
}

// mockDirEntry is a simple implementation of os.DirEntry for testing.
type mockDirEntry struct {
	name  string
	isDir bool
	mode  fs.FileMode
	info  os.FileInfo // optional, may be nil if you don’t need it
}

// Name returns the file name.
func (m mockDirEntry) Name() string {
	return m.name
}

// IsDir returns true if the entry represents a directory.
func (m mockDirEntry) IsDir() bool {
	return m.isDir
}

// Type returns the file mode bits that describe the file type.
func (m mockDirEntry) Type() fs.FileMode {
	return m.mode
}

// Info returns the os.FileInfo for the entry.
// In this simple mock, if m.info is nil, we return an error.
func (m mockDirEntry) Info() (os.FileInfo, error) {
	if m.info != nil {
		return m.info, nil
	}
	return nil, errors.New("file info not available")
}

// convertCommandMapToMocks converts a map of command strings to outputs
// into a slice of RunCommandMock structs
func convertCommandMapToMocks(commandMap map[string]string) []shared.RunCommandMock {
	mocks := []shared.RunCommandMock{}
	for cmd, output := range commandMap {
		parts := strings.Split(cmd, " ")
		command := parts[0]
		args := []string{}
		if len(parts) > 1 {
			args = parts[1:]
		}
		mocks = append(mocks, shared.RunCommandMock{
			Command: command,
			Args:    args,
			Out:     output,
			Err:     nil,
		})
	}
	return mocks
}
