package validator_test

import (
	"context"
	"testing"

	"github.com/kamal-github/demtech/internal/model"
	"github.com/kamal-github/demtech/internal/validator"
	"github.com/stretchr/testify/assert"
)

func TestMaxBodySizeValidator_Validate(t *testing.T) {
	tests := []struct {
		name      string
		bodyText  string
		maxSize   int64
		expectErr bool
	}{
		{
			name:      "Body within limit",
			bodyText:  "Short email body",
			maxSize:   100,
			expectErr: false,
		},
		{
			name:      "Body exactly at limit",
			bodyText:  "1234567890",
			maxSize:   10,
			expectErr: false,
		},
		{
			name:      "Body exceeding limit",
			bodyText:  "This email body is too long for the allowed size",
			maxSize:   20,
			expectErr: true,
		},
		{
			name:      "Empty body",
			bodyText:  "",
			maxSize:   50,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			v := validator.NewMaxBodySizeValidator(tt.maxSize)
			req := model.EmailRequest{Message: model.Message{Body: model.Body{Text: model.TextBody{Data: tt.bodyText}}}}

			err := v.Validate(context.Background(), req)

			if tt.expectErr {
				assert.Error(err)
				assert.IsType(&model.SESError{}, err)
				assert.Equal("MessageTooLong", err.(*model.SESError).Code)
			} else {
				assert.NoError(err)
			}
		})
	}
}
