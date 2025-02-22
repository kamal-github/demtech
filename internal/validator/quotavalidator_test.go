package validator_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kamal-github/demtech/internal/model"
	"github.com/kamal-github/demtech/internal/validator"
	"github.com/kamal-github/demtech/internal/validator/mocks"
	"github.com/stretchr/testify/assert"
)

func TestQuotaValidator_Validate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		sentEmails int64
		quota      int64
		mockError  error
		expectErr  bool
	}{
		{
			name:       "Within quota limit",
			sentEmails: 5,
			quota:      10,
			mockError:  nil,
			expectErr:  false,
		},
		{
			name:       "Exactly at quota limit",
			sentEmails: 10,
			quota:      10,
			mockError:  nil,
			expectErr:  true,
		},
		{
			name:       "Exceeding quota limit",
			sentEmails: 11,
			quota:      10,
			mockError:  nil,
			expectErr:  true,
		},
		{
			name:       "Error fetching sent email count",
			sentEmails: 0,
			quota:      10,
			mockError:  assert.AnError,
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			mockGetter := mocks.NewMockLastNHoursCountGetter(ctrl)
			mockGetter.EXPECT().GetLastNHoursCount(gomock.Any()).Return(tt.sentEmails, tt.mockError)

			v := validator.NewQuotaValidator(mockGetter, tt.quota)
			req := model.EmailRequest{}

			err := v.Validate(context.Background(), req)

			if tt.expectErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}
