package validator_test

import (
	"context"
	"testing"

	"github.com/kamal-github/demtech/internal/model"
	"github.com/kamal-github/demtech/internal/validator"
	"github.com/stretchr/testify/assert"
)

func TestSandboxValidator_Validate(t *testing.T) {
	tests := []struct {
		name         string
		awsIsSandbox bool
		allowed      []string
		dest         model.Destination
		expectErr    bool
	}{
		{
			name:         "All destinations allowed",
			awsIsSandbox: true,
			allowed:      []string{"allowed@example.com", "test@example.com"},
			dest:         model.Destination{ToAddresses: []string{"allowed@example.com"}},
			expectErr:    false,
		},
		{
			name:         "One unallowed destination - should fail",
			awsIsSandbox: true,
			allowed:      []string{"allowed@example.com"},
			dest:         model.Destination{ToAddresses: []string{"allowed@example.com", "notallowed@example.com"}},
			expectErr:    true,
		},
		{
			name:         "No allowed destinations - should fail",
			awsIsSandbox: true,
			allowed:      []string{},
			dest:         model.Destination{ToAddresses: []string{"random@example.com"}},
			expectErr:    true,
		},
		{
			name:         "Not a sandbox",
			awsIsSandbox: false,
			expectErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			v := validator.NewSandboxValidator(tt.awsIsSandbox, tt.allowed)
			req := model.EmailRequest{Destination: tt.dest}

			err := v.Validate(context.Background(), req)

			if tt.expectErr {
				assert.Error(err)
				assert.IsType(&model.SESError{}, err)
				assert.Equal("MessageRejected", err.(*model.SESError).Code)
			} else {
				assert.NoError(err)
			}
		})
	}
}
