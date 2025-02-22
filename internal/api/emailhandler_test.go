package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kamal-github/demtech/internal/api"
	"github.com/kamal-github/demtech/internal/api/mocks"
	"github.com/kamal-github/demtech/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestEmailHandler_SendEmailHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailService := mocks.NewMockEmailService(ctrl)
	mockStatsUpdater := mocks.NewMockEmailsStatsUpdater(ctrl) // Add mock for Stats Updater

	h := api.NewEmailHandler(mockEmailService, mockStatsUpdater)

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		requestBody   model.EmailRequest
		mockError     error
		expectCode    int
		mockCallsTime int
		mockStatsCall bool
		expectError   string
	}{
		{
			name: "Successful email send",
			requestBody: model.EmailRequest{
				Source: "test@example.com",
				Destination: model.Destination{
					ToAddresses: []string{"recipient@example.com"},
				},
				Message: model.Message{
					Subject: model.Subject{Data: "Test Subject"},
					Body:    model.Body{Text: model.TextBody{Data: "Test Body"}},
				},
				ReturnPath: "bounce@example.com",
			},
			mockCallsTime: 1,
			mockError:     nil,
			expectCode:    http.StatusOK,
			mockStatsCall: true, // Expect IncrementSuccess to be called
		},
		{
			name: "Email service failure",
			requestBody: model.EmailRequest{
				Source: "test@example.com",
				Destination: model.Destination{
					ToAddresses: []string{"recipient@example.com"},
				},
				Message: model.Message{
					Subject: model.Subject{Data: "Test Subject"},
					Body:    model.Body{Text: model.TextBody{Data: "Test Body"}},
				},
				ReturnPath: "bounce@example.com",
			},
			mockCallsTime: 1,
			mockError:     errors.New("send failed"),
			expectCode:    http.StatusInternalServerError,
			mockStatsCall: false,
			expectError:   "Failed to send email",
		},
		{
			name: "Validation failed, returnPath missing",
			requestBody: model.EmailRequest{
				Source: "test@example.com",
				Destination: model.Destination{
					ToAddresses: []string{"recipient@example.com"},
				},
				Message: model.Message{
					Subject: model.Subject{Data: "Test Subject"},
					Body:    model.Body{Text: model.TextBody{Data: "Test Body"}},
				},
			},
			mockCallsTime: 0,
			expectCode:    http.StatusBadRequest,
			mockStatsCall: false,
			expectError:   "EmailRequest.ReturnPath",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			mockEmailService.EXPECT().
				SendEmail(gomock.Any(), gomock.Any()).
				Return(&model.SESResponse{MessageID: "123"}, tt.mockError).
				Times(tt.mockCallsTime)

			if tt.expectCode == http.StatusBadRequest {
				mockStatsUpdater.EXPECT().IncrementError(gomock.Any(), "MissingOrInvalidParameterValue").Times(1)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonBody, _ := json.Marshal(tt.requestBody)
			c.Request, _ = http.NewRequest(http.MethodPost, "api/v1/send-email", bytes.NewBuffer(jsonBody))
			c.Request.Header.Set("Content-Type", "application/json")

			h.SendEmailHandler(c)

			assert.Equal(tt.expectCode, w.Code)

			if tt.expectError != "" {
				assert.Contains(w.Body.String(), tt.expectError)
			}
		})
	}
}
