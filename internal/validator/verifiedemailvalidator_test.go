package validator_test

import (
	"context"
	"testing"

	"github.com/kamal-github/demtech/internal/model"
	"github.com/kamal-github/demtech/internal/validator"
	"github.com/stretchr/testify/assert"
)

func TestVerifiedEmailValidator_Validate(t *testing.T) {
	tests := []struct {
		name      string
		verified  []string
		email     string
		expectErr bool
	}{
		{
			name:      "Valid email - found in verified list",
			verified:  []string{"verified@example.com", "test@example.com"},
			email:     "verified@example.com",
			expectErr: false,
		},
		{
			name:      "Invalid email - not in verified list",
			verified:  []string{"verified@example.com", "test@example.com"},
			email:     "unverified@example.com",
			expectErr: true,
		},
		{
			name:      "Empty verified list - should fail",
			verified:  []string{},
			email:     "some@example.com",
			expectErr: true,
		},
		{
			name:      "Case-sensitive match - should fail if case differs",
			verified:  []string{"verified@example.com"},
			email:     "Verified@example.com",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			v := validator.NewVerifiedEmailValidator(tt.verified)
			req := model.EmailRequest{Source: tt.email}

			err := v.Validate(context.Background(), req)

			if tt.expectErr {
				assert.Error(err)
				assert.IsType(&model.SESError{}, err)
				assert.Equal("MailFromDomainNotVerifiedException", err.(*model.SESError).Code)
			} else {
				assert.NoError(err)
			}
		})
	}
}
