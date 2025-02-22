package service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kamal-github/demtech/internal/model"
	"github.com/kamal-github/demtech/internal/service"
	"github.com/kamal-github/demtech/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestEmailServiceImpl_SendEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		validators []func(*mocks.MockValidator)
		trackerErr error
		expectErr  bool
	}{
		{
			name: "Successful email send",
			validators: []func(*mocks.MockValidator){
				func(mv *mocks.MockValidator) {
					mv.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				},
			},
			trackerErr: nil,
			expectErr:  false,
		},
		{
			name: "Validation failure",
			validators: []func(*mocks.MockValidator){
				func(mv *mocks.MockValidator) {
					mv.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(assert.AnError).Times(1)
				},
			},
			trackerErr: nil,
			expectErr:  true,
		},
		{
			name: "Tracking failure",
			validators: []func(*mocks.MockValidator){
				func(mv *mocks.MockValidator) {
					mv.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				},
			},
			trackerErr: assert.AnError,
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			mockValidator := mocks.NewMockValidator(ctrl)
			mockTracker := mocks.NewMockSentEmailTracker(ctrl)

			for _, v := range tt.validators {
				v(mockValidator)
			}

			mockTracker.EXPECT().TrackSentEmail(gomock.Any(), gomock.Any()).Return(tt.trackerErr).AnyTimes()

			es := service.NewEmailService([]service.Validator{mockValidator}, mockTracker, service.FailureConfig{FailRandomly: false, FailPercentage: 0})
			req := model.EmailRequest{
				Destination: model.Destination{ToAddresses: []string{"test@example.com"}},
			}

			_, err := es.SendEmail(context.Background(), req)

			if tt.expectErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}
