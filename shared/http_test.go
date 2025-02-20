package shared

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateSafeHTTPClient_SetUserAgent(t *testing.T) {
	// Set up a test HTTP server that checks the User-Agent header.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := fmt.Sprintf("Pareto Security/%s (Linux; build:%s)", Version, Commit)
		if got := r.Header.Get("User-Agent"); got != expected {
			t.Errorf("Expected User-Agent %q, got %q", expected, got)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create the client.
	client := createSafeHTTPClient()
	// Make a GET request to the test server.
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("client.Get error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, resp.StatusCode)
	}
}

func TestCreateSafeHTTPClient_TLSConfig(t *testing.T) {
	client := createSafeHTTPClient()

	// Assert that the client's transport is a userAgentTransport.
	uat, ok := client.Transport.(*userAgentTransport)
	if !ok {
		t.Fatalf("Expected client.Transport to be *userAgentTransport, got %T", client.Transport)
	}

	// Assert that the embedded transport is an *http.Transport.
	innerTrans, ok := uat.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("Expected userAgentTransport.Transport to be *http.Transport, got %T", uat.Transport)
	}

	// Check that the TLS config is not nil.
	if innerTrans.TLSClientConfig == nil {
		t.Error("TLSClientConfig is nil")
	} else if innerTrans.TLSClientConfig.MinVersion != tls.VersionTLS12 {
		t.Errorf("Expected TLS min version %v, got %v", tls.VersionTLS12, innerTrans.TLSClientConfig.MinVersion)
	}

	// Also test that the client timeout is set as expected.
	if client.Timeout != 30*time.Second {
		t.Errorf("Expected client timeout of %v, got %v", 30*time.Second, client.Timeout)
	}
}
