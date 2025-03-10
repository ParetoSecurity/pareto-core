package claims

import "github.com/ParetoSecurity/agent/check"

type Claim struct {
	Title  string
	Checks []check.Check
}
