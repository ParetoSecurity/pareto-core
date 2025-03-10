package checks

import (
	"testing"

	"github.com/ParetoSecurity/agent/claims"
	"github.com/samber/lo"
)

func TestClaims(t *testing.T) {
	uuids := []string{}

	for _, claim := range claims.All {
		for _, check := range claim.Checks {
			if lo.Contains(uuids, check.UUID()) {
				t.Errorf("Duplicate check UUID %s", check.UUID())
			}
			uuids = append(uuids, check.UUID())

			if check == nil {
				t.Errorf("Claim %s has a nil check", claim.Title)
			}
			check.RequiresRoot()
			check.Status()

			check.Passed()
			check.Name()
			check.PassedMessage()
			check.FailedMessage()
		}
	}
}
