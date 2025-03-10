package claims

import (
	"github.com/ParetoSecurity/agent/check"
	checks "github.com/ParetoSecurity/agent/checks/darwin"
	shared "github.com/ParetoSecurity/agent/checks/shared"
)

var All = []Claim{
	{"Access Security", []check.Check{
		check.Register(&shared.SSHKeys{}),
		check.Register(&shared.SSHKeysAlgo{}),
		check.Register(&checks.PasswordManagerCheck{}),
	}},
	{"Application Updates", []check.Check{
		check.Register(&shared.ParetoUpdated{}),
	}},
	{"Firewall & Sharing", []check.Check{
		check.Register(&shared.RemoteLogin{}),
	}},
	{"System Integrity", []check.Check{}},
}
