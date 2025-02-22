//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/kamal-github/demtech/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmailStatsAPI(t *testing.T) {
	apiBaseURL := os.Getenv("API_BASE_URL")
	emailReqs := []model.EmailRequest{
		{
			Source: "sender@example.com",
			Destination: model.Destination{
				ToAddresses:  []string{"recipient@example.com"},
				CcAddresses:  []string{"cc@example.com"},
				BccAddresses: []string{"bcc@example.com"},
			},
			Message: model.Message{
				Subject: model.Subject{Data: "Test Email Subject"},
				Body:    model.Body{Text: model.TextBody{Data: "This is the email body."}},
			},
			ConfigurationSetName: "default-config",
			ReplyToAddresses:     []string{"reply@example.com"},
			ReturnPath:           "bounce@example.com",
			ReturnPathArn:        "arn:aws:ses:us-east-1:123456789012:identity/bounce@example.com",
			SourceArn:            "arn:aws:ses:us-east-1:123456789012:identity/sender@example.com",
			Tags: []model.Tag{
				{Name: "campaign", Value: "welcome-email"},
				{Name: "userId", Value: "12345"},
			},
		},
		{
			// Source missing
			Destination: model.Destination{
				ToAddresses:  []string{"recipient@example.com"},
				CcAddresses:  []string{"cc@example.com"},
				BccAddresses: []string{"bcc@example.com"},
			},
			Message: model.Message{
				Subject: model.Subject{Data: "Test Email Subject"},
				Body:    model.Body{Text: model.TextBody{Data: "This is the email body."}},
			},
			ConfigurationSetName: "default-config",
			ReplyToAddresses:     []string{"reply@example.com"},
			ReturnPath:           "bounce@example.com",
			ReturnPathArn:        "arn:aws:ses:us-east-1:123456789012:identity/bounce@example.com",
			SourceArn:            "arn:aws:ses:us-east-1:123456789012:identity/sender@example.com",
			Tags: []model.Tag{
				{Name: "campaign", Value: "welcome-email"},
				{Name: "userId", Value: "12345"},
			},
		},
	}

	client := &http.Client{Timeout: 15 * time.Second}
	for _, emailReq := range emailReqs {
		body, err := json.Marshal(emailReq)
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, apiBaseURL+"/api/v1/send-email", bytes.NewBuffer(body))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		_, err = client.Do(req)
		require.NoError(t, err)
	}

	req, err := http.NewRequest(http.MethodGet, apiBaseURL+"/api/v1/email-stats", nil)
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var stats model.EmailStats
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&stats))
	exp := model.EmailStats{TotalEmailsSent: 2, SuccessCount: 1, TotalErrCount: 1, Errors: map[string]int{"MissingOrInvalidParameterValue": 1}}
	assert.Equal(t, exp, stats)
}
