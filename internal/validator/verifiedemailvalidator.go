package validator

import (
	"context"

	"github.com/kamal-github/demtech/internal/model"
)

/*
The message must be sent from a verified email address or
domain. If you attempt to send email using a non-verified
address or domain, the operation results in an
"Email address not verified" error.
*/
type VerifiedEmailValidator struct {
	awsVerifiedSourceEmailIDs []string
}

func NewVerifiedEmailValidator(v []string) VerifiedEmailValidator {
	return VerifiedEmailValidator{awsVerifiedSourceEmailIDs: v}
}

func (v VerifiedEmailValidator) Validate(ctx context.Context, req model.EmailRequest) error {
	for _, e := range v.awsVerifiedSourceEmailIDs {
		if e == req.Source {
			return nil
		}
	}

	return &model.SESError{Code: "EmailAddressNotVerified", Message: "Email address not verified"}
}
