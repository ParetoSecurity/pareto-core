package claims

import (
	"github.com/ParetoSecurity/agent/check"
	checks "github.com/ParetoSecurity/agent/checks/linux"
	shared "github.com/ParetoSecurity/agent/checks/shared"
)

var All = []Claim{
	{"Access Security", []check.Check{
		check.Register(&checks.Autologin{}),
		check.Register(&checks.DockerAccess{}),
		check.Register(&checks.PasswordToUnlock{}),
		check.Register(&shared.SSHKeys{}),
		check.Register(&shared.SSHKeysAlgo{}),
		check.Register(&checks.SSHConfigCheck{}),
		check.Register(&checks.PasswordManagerCheck{}),
	}},
	{"Application Updates", []check.Check{
		check.Register(&checks.ApplicationUpdates{}),
		check.Register(&shared.ParetoUpdated{}),
	}},
	{"Firewall & Sharing", []check.Check{
		check.Register(&checks.Firewall{}),
		check.Register(&checks.Printer{}),
		check.Register(&shared.RemoteLogin{}),
		check.Register(&checks.Sharing{}),
	}},
	{"System Integrity", []check.Check{
		check.Register(&checks.SecureBoot{}),
		check.Register(&checks.EncryptingFS{}),
	}},
}
