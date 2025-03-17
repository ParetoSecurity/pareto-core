package cmd

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/ParetoSecurity/agent/claims"
	"github.com/stretchr/testify/assert"
)

func TestHandleConnection(t *testing.T) {

	claims.All = []claims.Claim{}

	// Create a pair of connected sockets
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	// Run handleConnection in a separate goroutine
	go handleConnection(server)

	// Send a valid JSON payload with a "uuid" field
	input := map[string]string{"uuid": "test-uuid"}
	encoder := json.NewEncoder(client)
	if err := encoder.Encode(input); err != nil {
		t.Fatalf("failed to encode input: %v", err)
	}

	// Read the response from the helper
	var response map[string]bool
	decoder := json.NewDecoder(client)
	if err := decoder.Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Validate the response
	expected := map[string]bool{}
	assert.Equal(t, expected, response)
}
