package shared

import (
	"errors"

	"os/exec"
	"strings"
	"testing"

	"github.com/caarlos0/log"
)

// RunCommandMock represents a mock command with its arguments, output, and error
type RunCommandMock struct {
	Command string
	Args    []string
	Out     string
	Err     error
}

// RunCommandMocks is a slice that stores mock command outputs.
var RunCommandMocks []RunCommandMock

// RunCommand executes a command with the given name and arguments, and returns
// the combined standard output and standard error as a string. If testing is
// enabled, it returns a predefined fixture instead of executing the command.
func RunCommand(name string, arg ...string) (string, error) {

	// Check if testing is enabled and enable harnessing
	if testing.Testing() {
		for _, mock := range RunCommandMocks {
			if mock.Command == name && strings.Join(mock.Args, " ") == strings.Join(arg, " ") {
				return mock.Out, mock.Err
			}
		}
		return "", errors.New("RunCommand fixture not found: " + name + " " + strings.Join(arg, " "))
	}

	cmd := exec.Command(name, arg...)
	output, err := cmd.CombinedOutput()
	log.WithField("cmd", string(name+" "+strings.Join(arg, " "))).WithError(err).Debug(string(output))
	return string(output), err
}
