package shared

import (
	"os"
	"testing"
)

// checkPortMock is a mock function used for testing purposes. It simulates
// checking the availability of a port for a given protocol. The function
// takes an integer port number and a string representing the protocol
// (e.g., "tcp", "udp") as arguments, and returns a boolean indicating
// whether the port is available (true) or not (false).
var CheckPortMock func(port int, proto string) bool

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
