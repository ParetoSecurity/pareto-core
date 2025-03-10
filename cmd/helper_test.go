package cmd

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/ParetoSecurity/agent/claims"
)

func TestRunHelper(t *testing.T) {
	// Create a temporary UNIX socket path
	sockPath := filepath.Join(os.TempDir(), "pareto_helper_test.sock")
	_ = os.Remove(sockPath)
	ln, err := net.Listen("unix", sockPath)
	if err != nil {
		t.Fatalf("failed to listen on socket: %v", err)
	}
	defer ln.Close()
	defer os.Remove(sockPath)

	// Get the underlying file of the listener (dup'ed by the OS)
	unixLn, ok := ln.(*net.UnixListener)
	if !ok {
		t.Fatalf("listener is not a UnixListener")
	}
	lnFile, err := unixLn.File()
	if err != nil {
		t.Fatalf("failed to get file from listener: %v", err)
	}
	defer lnFile.Close()

	// Backup current fd 0
	oldFd, err := syscall.Dup(0)
	if err != nil {
		t.Fatalf("failed to duplicate fd 0: %v", err)
	}
	defer syscall.Dup2(oldFd, 0)
	syscall.Close(oldFd)

	// Replace fd 0 with our listener file descriptor
	if err = syscall.Dup2(int(lnFile.Fd()), 0); err != nil {
		t.Fatalf("failed to duplicate listener fd to 0: %v", err)
	}

	// Override claims.All so no check is run during test
	oldClaims := claims.All
	claims.All = []claims.Claim{}
	defer func() { claims.All = oldClaims }()

	// Run runHelper in a separate goroutine
	done := make(chan struct{})
	go func() {
		runHelper()
		close(done)
	}()

	// Connect to the actual UNIX socket
	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("failed to dial to socket: %v", err)
	}
	defer conn.Close()

	// Send a valid JSON payload with a "uuid" field
	input := map[string]string{"uuid": "test"}
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(input); err != nil {
		t.Fatalf("failed to encode input: %v", err)
	}

	// Read the response from the helper
	var response map[string]bool
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Expect an empty response since no check is run (claims.All is empty)
	if len(response) != 0 {
		t.Errorf("expected empty response, got %v", response)
	}

	// Wait for runHelper to finish (with a timeout)
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("runHelper did not terminate within the timeout")
	}
}
