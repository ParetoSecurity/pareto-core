package shared

// checkPortMock is a mock function used for testing purposes. It simulates
// checking the availability of a port for a given protocol. The function
// takes an integer port number and a string representing the protocol
// (e.g., "tcp", "udp") as arguments, and returns a boolean indicating
// whether the port is available (true) or not (false).
var CheckPortMock func(port int, proto string) bool
