package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kamal-github/demtech/internal/api"
	"github.com/kamal-github/demtech/internal/api/mocks"
	"github.com/kamal-github/demtech/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestGetEmailStats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockEmailStatsService(ctrl)
	handler := api.NewEmailStatsHandler(mockService)

	tests := []struct {
		name           string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - returns email stats",
			mockSetup: func() {
				mockService.EXPECT().
					GetEmailStats(gomock.Any()).
					Return(model.EmailStats{TotalEmailsSent: 10, TotalErrCount: 0, SuccessCount: 10, Errors: map[string]int{}}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"totalEmailsSent":10,"totalErrCount":0, "successCount": 10}`,
		},
		{
			name: "Failure - service returns error",
			mockSetup: func() {
				mockService.EXPECT().
					GetEmailStats(gomock.Any()).
					Return(model.EmailStats{}, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"service error"}`,
		},
		{
			name: "Failure - service times out",
			mockSetup: func() {
				mockService.EXPECT().
					GetEmailStats(gomock.Any()).
					DoAndReturn(func(ctx context.Context) (model.EmailStats, error) {
						time.Sleep(2 * time.Second) // Simulate delay
						return model.EmailStats{}, context.DeadlineExceeded
					})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"context deadline exceeded"}`,
		},
		{
			name: "Success - empty email stats",
			mockSetup: func() {
				mockService.EXPECT().
					GetEmailStats(gomock.Any()).
					Return(model.EmailStats{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"successCount":0, "totalEmailsSent":0, "totalErrCount":0}`, // Empty JSON
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			router := gin.Default()
			router.GET("/api/v1/email-stats", handler.GetEmailStats)

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/email-stats", nil)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)
			assert.Equal(t, tt.expectedStatus, recorder.Code)

			// Normalize and compare JSON response
			var expectedJSON map[string]interface{}
			var actualJSON map[string]interface{}
			_ = json.Unmarshal([]byte(tt.expectedBody), &expectedJSON)
			_ = json.Unmarshal(recorder.Body.Bytes(), &actualJSON)

			assert.Equal(t, expectedJSON, actualJSON)
		})
	}
}
