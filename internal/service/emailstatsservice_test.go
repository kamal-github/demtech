package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kamal-github/demtech/internal/model"
	"github.com/kamal-github/demtech/internal/service"
	"github.com/kamal-github/demtech/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestEmailStatsService_SendEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailService := mocks.NewMockEmailService(ctrl)
	mockStatsUpdater := mocks.NewMockEmailsStatsUpdater(ctrl)

	tests := []struct {
		name      string
		emailErr  error
		expectErr bool
		statErr   error
	}{
		{
			name:      "Successful email send",
			emailErr:  nil,
			statErr:   nil,
			expectErr: false,
		},
		{
			name:      "Email service failure",
			emailErr:  errors.New("send failed"),
			statErr:   nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			mockEmailService.EXPECT().SendEmail(gomock.Any(), gomock.Any()).Return(&model.SESResponse{MessageID: "123"}, tt.emailErr).Times(1)

			if tt.emailErr == nil {
				mockStatsUpdater.EXPECT().IncrementSuccess(gomock.Any()).Return(tt.statErr).Times(1)
			} else {
				mockStatsUpdater.EXPECT().IncrementError(gomock.Any(), gomock.Any()).Return(tt.statErr).Times(1)
			}

			es := service.NewEmailStatsService(mockEmailService, mockStatsUpdater, nil)

			_, err := es.SendEmail(context.Background(), model.EmailRequest{})

			if tt.expectErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestGetEmailStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStatsGetter := mocks.NewMockEmailsStatsGetter(ctrl)

	// Create service with the mocked EmailsStatsGetter
	emailStatsService := service.NewEmailStatsService(nil, nil, mockStatsGetter)

	tests := []struct {
		name           string
		mockSetup      func()
		expectedResult model.EmailStats
		expectedError  error
	}{
		{
			name: "Success - valid email stats",
			mockSetup: func() {
				mockStatsGetter.EXPECT().
					GetEmailStats(gomock.Any()).
					Return(model.EmailStats{TotalEmailsSent: 10, TotalErrCount: 0, SuccessCount: 10, Errors: map[string]int{}}, nil)
			},
			expectedResult: model.EmailStats{TotalEmailsSent: 10, TotalErrCount: 0, SuccessCount: 10, Errors: map[string]int{}},
			expectedError:  nil,
		},
		{
			name: "Failure - storage error",
			mockSetup: func() {
				mockStatsGetter.EXPECT().
					GetEmailStats(gomock.Any()).
					Return(model.EmailStats{}, errors.New("database error"))
			},
			expectedResult: model.EmailStats{},
			expectedError:  errors.New("database error"),
		},
		{
			name: "Success - empty email stats",
			mockSetup: func() {
				mockStatsGetter.EXPECT().
					GetEmailStats(gomock.Any()).
					Return(model.EmailStats{}, nil)
			},
			expectedResult: model.EmailStats{},
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			// Call the method
			result, err := emailStatsService.GetEmailStats(context.Background())

			// Validate results
			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
