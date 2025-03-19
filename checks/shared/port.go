package shared

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/caarlos0/log"
)

// checkPort tests if a port is open
func CheckPort(port int, proto string) bool {

	if testing.Testing() {
		return CheckPortMock(port, proto)
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false
	}

	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			continue
		}

		// Filter out 127.0.0.1
		if ip.IsLoopback() {
			continue
		}

		address := net.JoinHostPort(ip.String(), fmt.Sprintf("%d", port))
		conn, err := net.DialTimeout(proto, address, 1*time.Second)
		if err == nil {
			defer conn.Close()
			log.WithField("address", address).WithField("state", true).Debug("Checking port")
			return true
		}
	}

	return false
}
