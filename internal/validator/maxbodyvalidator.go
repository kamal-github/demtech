package validator

import (
	"context"

	"github.com/kamal-github/demtech/internal/model"
)

type MaxBodySizeValidator struct {
	awsMaxEmailSizeAllowedBytes int64
}

func NewMaxBodySizeValidator(m int64) MaxBodySizeValidator {
	return MaxBodySizeValidator{awsMaxEmailSizeAllowedBytes: m}
}

func (v MaxBodySizeValidator) Validate(ctx context.Context, req model.EmailRequest) error {
	if int64(len(req.Message.Body.Text.Data)) > v.awsMaxEmailSizeAllowedBytes {
		return &model.SESError{Code: "MessageTooLong", Message: "Email body exceeds maximum size"}
	}

	return nil
}
