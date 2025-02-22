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
)

func TestSendEmailAPI(t *testing.T) {
	apiBaseURL := os.Getenv("API_BASE_URL")
	emailReq := model.EmailRequest{
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
	}
	body, err := json.Marshal(emailReq)
	assert.NoError(t, err)
	req, err := http.NewRequest("POST", apiBaseURL+"/api/v1/send-email", bytes.NewBuffer(body))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
