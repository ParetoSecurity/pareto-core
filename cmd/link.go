package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	shared "github.com/ParetoSecurity/pareto-core/shared"
	"github.com/ParetoSecurity/pareto-core/team"
	"github.com/caarlos0/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var rsaPublicKey = `
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwGh64DK49GOq1KX+ojyg
Y9JSAZ4cfm5apavetQ42D2gTjfhDu1kivrDRwhjqj7huUWRI2ExMdMHp8CzrJI3P
zpzutEUXTEHloe0vVMZqPoP/r2f1cl4bmDkFZyHr6XTgiYPE4GgMjxUc04J2ksqU
/XbNwOVsBiuy1T2BduLYiYr1UyIx8VqEb+3tunQKlyRKF7a5LoEZatt5F/5vaMMI
4zp1yIc2PMoBdlBH4/tpJmC/PiwjBuwgp5gMIle4Hy7zwW4+rIJzF5P3Tg+Am+Lg
davB8TIZDBlqIWV7zK1kWBPj364a5cnaUP90BnOriMJBh7zPG0FNGTXTiJED2qDM
fajDrji3oAPO24mJsCCzSd8LIREK5c6iAf1X4UI/UFP+UhOBCsANrhNSXRpO2KyM
+60JYzFpMvyhdK9zMo7Tc+KM6R0YRNmBCYK/ePAGk3WU6qxN5+OmSjdTvFrqC4JQ
FyK51WJI80PKvp3B7ZB7XpH5B24wr/OhMRh5YZOcrpuBykfHaMozkDCudgaj/V+x
K79CqMF/BcSxCSBktWQmabYCM164utpmJaCSpZyDtKA4bYVv9iRCGTqFQT7jX+/h
Z37gmg/+TlIdTAeB5TG2ffHxLnRhT4AAhUgYmk+QP3a1hxP5xj2otaSTZ3DxQd6F
ZaoGJg3y8zjrxYBQDC8gF6sCAwEAAQ==
`

type InviteClaims struct {
	TeamAuth string `json:"token"`
	TeamUUID string `json:"teamID"`
	jwt.RegisteredClaims
}

var linkCmd = &cobra.Command{
	Use:   "link <url>",
	Short: "Link team with this device",
	Run: func(cc *cobra.Command, args []string) {
		runLinkCommand(args[0])
	},
}

func runLinkCommand(teamURL string) {
	if shared.IsLinked() {
		log.Warn("Already linked to a team")
		log.Warn("Unlink first with `pareto unlink`")
		log.Infof("Team ID: %s", shared.Config.TeamID)
		os.Exit(1)
	}

	if lo.IsNotEmpty(teamURL) {
		token, err := getTokenFromURL(teamURL)
		if err != nil {
			log.WithError(err).Warn("failed to get token from URL")
			os.Exit(1)
		}

		parsedToken, err := parseJWT(token)
		if err != nil {
			log.WithError(err).Warn("failed to parse JWT")
			os.Exit(1)
		}

		shared.Config.TeamID = parsedToken.TeamUUID
		shared.Config.AuthToken = parsedToken.TeamAuth

		err = team.ReportToTeam(true)
		if err != nil {
			log.WithError(err).Warn("failed to report to team")
			os.Exit(1)
		}

		err = shared.SaveConfig()
		if err != nil {
			log.Errorf("Error saving config: %v", err)
			os.Exit(1)
		}
		log.Infof("Device successfully linked to team: %s", parsedToken.TeamUUID)
	}
}

func getTokenFromURL(teamURL string) (string, error) {

	parsedURL, err := url.Parse(teamURL)
	if err != nil {
		return "", err
	}

	queryParams := parsedURL.Query()
	token := queryParams.Get("token")
	if token == "" {
		return "", fmt.Errorf("token not found in URL")
	}

	return token, nil
}

func parseJWT(token string) (*InviteClaims, error) {
	jwttToken, _ := jwt.ParseWithClaims(token, &InviteClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(strings.ReplaceAll(rsaPublicKey, "\n", "")), nil
	})
	if claims, ok := jwttToken.Claims.(*InviteClaims); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("failed to parse JWT")
}

func init() {
	rootCmd.AddCommand(linkCmd)
}
