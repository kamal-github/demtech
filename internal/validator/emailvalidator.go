package validator

import (
	"context"
	"regexp"

	"github.com/kamal-github/demtech/internal/model"
)

type EmailValidator struct{}

func NewEmailValidator() EmailValidator {
	return EmailValidator{}
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func (v EmailValidator) Validate(ctx context.Context, req model.EmailRequest) error {
	dest := req.Destination
	for _, email := range dest.All() {
		if !emailRegex.MatchString(email) {
			return &model.SESError{Code: "InvalidParameterValue", Message: "Invalid recipient email address"}
		}
	}

	return nil
}
