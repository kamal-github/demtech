package validator

import (
	"context"

	"github.com/kamal-github/demtech/internal/model"
)

type MaxDestinationsValidator struct {
	awsMaxDestinationsAllowed int
}

func NewMaxDestinationsValidator(m int) MaxDestinationsValidator {
	return MaxDestinationsValidator{awsMaxDestinationsAllowed: m}
}

func (v MaxDestinationsValidator) Validate(ctx context.Context, req model.EmailRequest) error {
	if len(req.Destination.All()) > v.awsMaxDestinationsAllowed {
		return &model.SESError{Code: "LimitExceededException", Message: "Too many recipients"}
	}

	return nil
}
