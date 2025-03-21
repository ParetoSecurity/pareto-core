package cmd

import (
	"testing"

	"github.com/ParetoSecurity/agent/shared"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestParseJWT(t *testing.T) {
	validToken := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJqdEBuaXRlby5jbyIsInRlYW1JRCI6IjI0MjljNDllLTM3YmItNDFiYi05MDc3LTZiYjYyMDJlMjU1YiIsInJvbGUiOiJ0ZWFtIiwiaWF0IjoxNzM2NDE3NDEwLCJ0b2tlbiI6ImV5SmhiR2NpT2lKSVV6VXhNaUlzSW5SNWNDSTZJa3BYVkNKOS5leUp5YjJ4bGN5STZXeUpzYVc1clgyUmxkbWxqWlNJc0luVndaR0YwWlY5a1pYWnBZMlVpWFN3aWRHVmhiVjlwWkNJNklqSTBNamxqTkRsbExUTTNZbUl0TkRGaVlpMDVNRGMzTFRaaVlqWXlNREpsTWpVMVlpSXNJbWx1ZG1sMFpWOXBaQ0k2SWpSaFpqZzFNVFJoTFdObE5qTXRORGMwTnkwNE1EZG1MVFZqTXpnek9XUTNPRE0wTVNJc0luTjFZaUk2SW1wMFFHNXBkR1Z2TG1Odklpd2lhV0YwSWpveE56TTJOREUzTkRFd2ZRLnZTSFU0Nm5yZWo2aVRvZHdIaEJLRUxTREUxN3dxQzJoc0VRX1RYUmN4Z0ZaN3dsSHptUkNnTFlEOGtwenpSTzdvME85dTdnYXppTXJQTF9vUHVPdXJRIn0.kZqUzuRO7R9Bd6U8krlRj18CmmRMX1uwUNToYwVn-OYsCViP0ae--Mbo4E4brWrtXUm0PXVQLhR0Ml0xeTNLJx7JNVPFPCCOugNLAvL42g3RL7nk3kjYZ2ugbvK_uGrQTtFZojRTkYpDv3YgKpeNpoMpmT3GTK9PRG3YXkfXkPgZyrIrLwaXn57Tr88MOcFbyq1VD5M1UPizGHJDfkmldP4ROmKSEfc8iNcIrYV7uIcqBWoTzqKnLxjG6FQ9Ylsrw_-kpfzfa-8tbaWrhY-UgjSllY4WUUG95tkLVxlKHcKDZHsYWXWZO-nMdZF7JlFN8MpPEJDCq_E9tOVqbWcEh1DCWrXa33Sm5ZfvdSBBkhzUnvTwDTDjCDCMhA9gNcdMfEoKCh11lDD8r3FRvIlioBVKZ3GNm25AtfbcypH8jobdnUIBrgtrPxyadv63o0IEshtTX4kswUkGqvwMlDD-r-J2oPrEkN_JRJshTpYezUagIEvYvXAPjNU2kVWOJFnS9MCLuJa4Di99omEnS9oRemgJP0tR6Z84sbTiXJJIsa1sEY8MZDAqXD1U1OHtfAo7vL5z5SyPjQPnaKMacttNx0gfHFA1rP2Vdsj5m6nYQtBqZpFUVvOgKa3bZQRYOWho1IF22dhdqOBgjUG2CPe8fZ0mpAaQMJ5SHUh613ShOyI"
	invalidToken := "invalid.token.string"

	t.Run("valid token", func(t *testing.T) {
		claims, err := parseJWT(validToken)
		assert.NoError(t, err)
		assert.Equal(t, "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJyb2xlcyI6WyJsaW5rX2RldmljZSIsInVwZGF0ZV9kZXZpY2UiXSwidGVhbV9pZCI6IjI0MjljNDllLTM3YmItNDFiYi05MDc3LTZiYjYyMDJlMjU1YiIsImludml0ZV9pZCI6IjRhZjg1MTRhLWNlNjMtNDc0Ny04MDdmLTVjMzgzOWQ3ODM0MSIsInN1YiI6Imp0QG5pdGVvLmNvIiwiaWF0IjoxNzM2NDE3NDEwfQ.vSHU46nrej6iTodwHhBKELSDE17wqC2hsEQ_TXRcxgFZ7wlHzmRCgLYD8kpzzRO7o0O9u7gaziMrPL_oPuOurQ", claims.TeamAuth)
		assert.Equal(t, "2429c49e-37bb-41bb-9077-6bb6202e255b", claims.TeamUUID)
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := parseJWT(invalidToken)
		assert.Error(t, err)
	})
}

func TestRunLinkCommand_Success(t *testing.T) {
	defer gock.Off()

	// Set default ReportURL for test
	team.ReportURL = "https://dash.paretosecurity.com"

	gock.New("https://dash.paretosecurity.com").
		Reply(200).
		JSON([]map[string]string{{"status": "ok"}})

	// Reset config
	shared.Config.TeamID = ""
	shared.Config.AuthToken = ""

	defer func() {
		shared.Config.TeamID = ""
		shared.Config.AuthToken = ""
		shared.SaveConfig()
	}()

	// Use the same valid token as in TestParseJWT.
	validToken := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJqdEBuaXRlby5jbyIsInRlYW1JRCI6IjI0MjljNDllLTM3YmItNDFiYi05MDc3LTZiYjYyMDJlMjU1YiIsInJvbGUiOiJ0ZWFtIiwiaWF0IjoxNzM2NDE3NDEwLCJ0b2tlbiI6ImV5SmhiR2NpT2lKSVV6VXhNaUlzSW5SNWNDSTZJa3BYVkNKOS5leUp5YjJ4bGN5STZXeUpzYVc1clgyUmxkbWxqWlNJc0luVndaR0YwWlY5a1pYWnBZMlVpWFN3aWRHVmhiVjlwWkNJNklqSTBNamxqTkRsbExUTTNZbUl0TkRGaVlpMDVNRGMzTFRaaVlqWXlNREpsTWpVMVlpSXNJbWx1ZG1sMFpWOXBaQ0k2SWpSaFpqZzFNVFJoTFdObE5qTXRORGMwTnkwNE1EZG1MVFZqTXpnek9XUTNPRE0wTVNJc0luTjFZaUk2SW1wMFFHNXBkR1Z2TG1Odklpd2lhV0YwSWpveE56TTJOREUzTkRFd2ZRLnZTSFU0Nm5yZWo2aVRvZHdIaEJLRUxTREUxN3dxQzJoc0VRX1RYUmN4Z0ZaN3dsSHptUkNnTFlEOGtwenpSTzdvME85dTdnYXppTXJQTF9vUHVPdXJRIn0.kZqUzuRO7R9Bd6U8krlRj18CmmRMX1uwUNToYwVn-OYsCViP0ae--Mbo4E4brWrtXUm0PXVQLhR0Ml0xeTNLJx7JNVPFPCCOugNLAvL42g3RL7nk3kjYZ2ugbvK_uGrQTtFZojRTkYpDv3YgKpeNpoMpmT3GTK9PRG3YXkfXkPgZyrIrLwaXn57Tr88MOcFbyq1VD5M1UPizGHJDfkmldP4ROmKSEfc8iNcIrYV7uIcqBWoTzqKnLxjG6FQ9Ylsrw_-kpfzfa-8tbaWrhY-UgjSllY4WUUG95tkLVxlKHcKDZHsYWXWZO-nMdZF7JlFN8MpPEJDCq_E9tOVqbWcEh1DCWrXa33Sm5ZfvdSBBkhzUnvTwDTDjCDCMhA9gNcdMfEoKCh11lDD8r3FRvIlioBVKZ3GNm25AtfbcypH8jobdnUIBrgtrPxyadv63o0IEshtTX4kswUkGqvwMlDD-r-J2oPrEkN_JRJshTpYezUagIEvYvXAPjNU2kVWOJFnS9MCLuJa4Di99omEnS9oRemgJP0tR6Z84sbTiXJJIsa1sEY8MZDAqXD1U1OHtfAo7vL5z5SyPjQPnaKMacttNx0gfHFA1rP2Vdsj5m6nYQtBqZpFUVvOgKa3bZQRYOWho1IF22dhdqOBgjUG2CPe8fZ0mpAaQMJ5SHUh613ShOyI"
	expectedTeamUUID := "2429c49e-37bb-41bb-9077-6bb6202e255b"
	expectedTeamAuth := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJyb2xlcyI6WyJsaW5rX2RldmljZSIsInVwZGF0ZV9kZXZpY2UiXSwidGVhbV9pZCI6IjI0MjljNDllLTM3YmItNDFiYi05MDc3LTZiYjYyMDJlMjU1YiIsImludml0ZV9pZCI6IjRhZjg1MTRhLWNlNjMtNDc0Ny04MDdmLTVjMzgzOWQ3ODM0MSIsInN1YiI6Imp0QG5pdGVvLmNvIiwiaWF0IjoxNzM2NDE3NDEwfQ.vSHU46nrej6iTodwHhBKELSDE17wqC2hsEQ_TXRcxgFZ7wlHzmRCgLYD8kpzzRO7o0O9u7gaziMrPL_oPuOurQ"

	// Construct the URL with the valid token.
	url := "http://example.com?token=" + validToken

	// Call the function under test.
	runLinkCommand(url)

	// Assert that shared.Config was updated.
	assert.Equal(t, expectedTeamUUID, shared.Config.TeamID)
	assert.Equal(t, expectedTeamAuth, shared.Config.AuthToken)
}

func TestRunLinkCommand_CustomServer(t *testing.T) {
	defer gock.Off()

	customServer := "https://custom.example.com"
	team.ReportURL = customServer

	gock.New(customServer).
		Reply(200).
		JSON([]map[string]string{{"status": "ok"}})

	// Reset config
	shared.Config.TeamID = ""
	shared.Config.AuthToken = ""

	defer func() {
		shared.Config.TeamID = ""
		shared.Config.AuthToken = ""
		shared.SaveConfig()
		// Reset to default
		team.ReportURL = "https://dash.paretosecurity.com"
	}()

	// Use the same valid token as in TestParseJWT.
	validToken := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJqdEBuaXRlby5jbyIsInRlYW1JRCI6IjI0MjljNDllLTM3YmItNDFiYi05MDc3LTZiYjYyMDJlMjU1YiIsInJvbGUiOiJ0ZWFtIiwiaWF0IjoxNzM2NDE3NDEwLCJ0b2tlbiI6ImV5SmhiR2NpT2lKSVV6VXhNaUlzSW5SNWNDSTZJa3BYVkNKOS5leUp5YjJ4bGN5STZXeUpzYVc1clgyUmxkbWxqWlNJc0luVndaR0YwWlY5a1pYWnBZMlVpWFN3aWRHVmhiVjlwWkNJNklqSTBNamxqTkRsbExUTTNZbUl0TkRGaVlpMDVNRGMzTFRaaVlqWXlNREpsTWpVMVlpSXNJbWx1ZG1sMFpWOXBaQ0k2SWpSaFpqZzFNVFJoTFdObE5qTXRORGMwTnkwNE1EZG1MVFZqTXpnek9XUTNPRE0wTVNJc0luTjFZaUk2SW1wMFFHNXBkR1Z2TG1Odklpd2lhV0YwSWpveE56TTJOREUzTkRFd2ZRLnZTSFU0Nm5yZWo2aVRvZHdIaEJLRUxTREUxN3dxQzJoc0VRX1RYUmN4Z0ZaN3dsSHptUkNnTFlEOGtwenpSTzdvME85dTdnYXppTXJQTF9vUHVPdXJRIn0.kZqUzuRO7R9Bd6U8krlRj18CmmRMX1uwUNToYwVn-OYsCViP0ae--Mbo4E4brWrtXUm0PXVQLhR0Ml0xeTNLJx7JNVPFPCCOugNLAvL42g3RL7nk3kjYZ2ugbvK_uGrQTtFZojRTkYpDv3YgKpeNpoMpmT3GTK9PRG3YXkfXkPgZyrIrLwaXn57Tr88MOcFbyq1VD5M1UPizGHJDfkmldP4ROmKSEfc8iNcIrYV7uIcqBWoTzqKnLxjG6FQ9Ylsrw_-kpfzfa-8tbaWrhY-UgjSllY4WUUG95tkLVxlKHcKDZHsYWXWZO-nMdZF7JlFN8MpPEJDCq_E9tOVqbWcEh1DCWrXa33Sm5ZfvdSBBkhzUnvTwDTDjCDCMhA9gNcdMfEoKCh11lDD8r3FRvIlioBVKZ3GNm25AtfbcypH8jobdnUIBrgtrPxyadv63o0IEshtTX4kswUkGqvwMlDD-r-J2oPrEkN_JRJshTpYezUagIEvYvXAPjNU2kVWOJFnS9MCLuJa4Di99omEnS9oRemgJP0tR6Z84sbTiXJJIsa1sEY8MZDAqXD1U1OHtfAo7vL5z5SyPjQPnaKMacttNx0gfHFA1rP2Vdsj5m6nYQtBqZpFUVvOgKa3bZQRYOWho1IF22dhdqOBgjUG2CPe8fZ0mpAaQMJ5SHUh613ShOyI"
	expectedTeamUUID := "2429c49e-37bb-41bb-9077-6bb6202e255b"
	expectedTeamAuth := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJyb2xlcyI6WyJsaW5rX2RldmljZSIsInVwZGF0ZV9kZXZpY2UiXSwidGVhbV9pZCI6IjI0MjljNDllLTM3YmItNDFiYi05MDc3LTZiYjYyMDJlMjU1YiIsImludml0ZV9pZCI6IjRhZjg1MTRhLWNlNjMtNDc0Ny04MDdmLTVjMzgzOWQ3ODM0MSIsInN1YiI6Imp0QG5pdGVvLmNvIiwiaWF0IjoxNzM2NDE3NDEwfQ.vSHU46nrej6iTodwHhBKELSDE17wqC2hsEQ_TXRcxgFZ7wlHzmRCgLYD8kpzzRO7o0O9u7gaziMrPL_oPuOurQ"

	// Construct the URL with the valid token.
	url := "http://example.com?token=" + validToken

	// Call the function under test.
	runLinkCommand(url)

	// Assert that shared.Config was updated.
	assert.Equal(t, expectedTeamUUID, shared.Config.TeamID)
	assert.Equal(t, expectedTeamAuth, shared.Config.AuthToken)

	// Verify that custom server URL was used
	assert.Equal(t, customServer, team.ReportURL)
}
