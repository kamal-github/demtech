package validator_test

import (
	"context"
	"testing"

	"github.com/kamal-github/demtech/internal/model"
	"github.com/kamal-github/demtech/internal/validator"
	"github.com/stretchr/testify/assert"
)

func TestMaxDestinationsValidator_Validate(t *testing.T) {
	tests := []struct {
		name       string
		destEmails []string
		maxDest    int
		expectErr  bool
	}{
		{
			name:       "Destinations within limit",
			destEmails: []string{"a@example.com", "b@example.com"},
			maxDest:    3,
			expectErr:  false,
		},
		{
			name:       "Destinations exactly at limit",
			destEmails: []string{"a@example.com", "b@example.com", "c@example.com"},
			maxDest:    3,
			expectErr:  false,
		},
		{
			name:       "Destinations exceeding limit",
			destEmails: []string{"a@example.com", "b@example.com", "c@example.com", "d@example.com"},
			maxDest:    3,
			expectErr:  true,
		},
		{
			name:       "No destinations",
			destEmails: []string{},
			maxDest:    3,
			expectErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			v := validator.NewMaxDestinationsValidator(tt.maxDest)
			req := model.EmailRequest{Destination: model.Destination{ToAddresses: tt.destEmails}}

			err := v.Validate(context.Background(), req)

			if tt.expectErr {
				assert.Error(err)
				assert.IsType(&model.SESError{}, err)
				assert.Equal("LimitExceededException", err.(*model.SESError).Code)
			} else {
				assert.NoError(err)
			}
		})
	}
}
