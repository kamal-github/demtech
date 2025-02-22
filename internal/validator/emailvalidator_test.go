package validator_test

import (
	"context"
	"testing"

	"github.com/kamal-github/demtech/internal/model"
	"github.com/kamal-github/demtech/internal/validator"
	"github.com/stretchr/testify/assert"
)

func TestEmailValidator_Validate(t *testing.T) {
	tests := []struct {
		name      string
		dest      model.Destination
		expectErr bool
	}{
		{
			name: "All valid email addresses",
			dest: model.Destination{
				ToAddresses:  []string{"valid@example.com"},
				CcAddresses:  []string{"cc@example.com"},
				BccAddresses: []string{"bcc@example.com"},
			},
			expectErr: false,
		},
		{
			name: "One invalid email address",
			dest: model.Destination{
				ToAddresses: []string{"invalid-email"},
			},
			expectErr: true,
		},
		{
			name: "Multiple valid and invalid email addresses",
			dest: model.Destination{
				ToAddresses:  []string{"valid@example.com", "invalid-email"},
				CcAddresses:  []string{"cc@example.com"},
				BccAddresses: []string{"bcc@example.com"},
			},
			expectErr: true,
		},
		{
			name:      "Empty destination - should pass",
			dest:      model.Destination{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			v := validator.NewEmailValidator()
			req := model.EmailRequest{Destination: tt.dest}

			err := v.Validate(context.Background(), req)

			if tt.expectErr {
				assert.Error(err)
				assert.IsType(&model.SESError{}, err)
				assert.Equal("InvalidParameterValue", err.(*model.SESError).Code)
			} else {
				assert.NoError(err)
			}
		})
	}
}
