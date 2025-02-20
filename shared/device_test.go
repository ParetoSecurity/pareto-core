package shared

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeviceAuth(t *testing.T) {
	tests := []struct {
		name      string
		authToken string
		expected  string
	}{
		{
			name:      "Empty AuthToken",
			authToken: "",
			expected:  "",
		},
		{
			name:      "Invalid Base64 AuthToken",
			authToken: "invalid.token",
			expected:  "",
		},
		{
			name: "Valid AuthToken",
			authToken: func() string {
				payload := map[string]string{
					"sub":    "test-sub",
					"teamID": "test-teamID",
					"role":   "test-role",
					"token":  "test-token",
				}
				payloadBytes, _ := json.Marshal(payload)
				encodedPayload := base64.RawURLEncoding.EncodeToString(payloadBytes)
				return "header." + encodedPayload + ".signature"
			}(),
			expected: "test-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Config.AuthToken = tt.authToken
			result := DeviceAuth()
			assert.Equal(t, tt.expected, result)
		})
	}
}
