package shared

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/caarlos0/log"
	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/requests/reqtest"
)

var (
	once              sync.Once
	safeHTTPTransport http.RoundTripper
)

type userAgentTransport struct {
	Transport http.RoundTripper
	agent     string
}

func (uat *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", uat.agent)
	return uat.Transport.RoundTrip(req)
}

func createSafeHTTPClient() *http.Client {
	// Define a custom transport with TLS configurations
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			// Enforce secure TLS settings
			MinVersion: tls.VersionTLS12,
			// Reject insecure certificates (set to true for development)
			InsecureSkipVerify: false,
		},
	}

	// Create a custom HTTP client
	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // Set a timeout to prevent hanging requests
	}

	// Add the custom User-Agent to all requests
	customTransport := &http.Transport{
		TLSClientConfig: transport.TLSClientConfig,
		Proxy:           transport.Proxy,
		DialContext:     transport.DialContext,
	}

	// Wrap the transport to insert the User-Agent header
	client.Transport = &userAgentTransport{
		Transport: customTransport,
		agent:     fmt.Sprintf("Pareto Security/%s (build:%s)", Version, Commit),
	}

	return client
}

func transport() http.RoundTripper {
	baseTrans := createSafeHTTPClient().Transport
	if testing.Testing() {
		return reqtest.Caching(baseTrans, "fixtures")
	}
	logger := func(req *http.Request, res *http.Response, err error, d time.Duration) {
		log.Debugf("method=%q url=%q err=%v status=%q duration=%v\n",
			req.Method, req.URL, err, res.Status, d.Round(1*time.Second))
	}

	return requests.LogTransport(baseTrans, logger)
}

// SafeHTTPTransport is a custom HTTP transport with secure TLS settings and a custom User-Agent.
func HTTPTransport() http.RoundTripper {
	once.Do(func() {
		safeHTTPTransport = transport()
	})
	return safeHTTPTransport
}
