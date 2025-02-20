package validator

import (
	"context"

	"github.com/kamal-github/demtech/internal/model"
)

type SandboxValidator struct {
	awsSandboxAllowedDestinations []string
}

func NewSandboxValidator(e []string) SandboxValidator {
	return SandboxValidator{awsSandboxAllowedDestinations: e}
}

func (v SandboxValidator) Validate(ctx context.Context, req model.EmailRequest) error {
	sandboxEmails := make(map[string]struct{})
	for _, e := range v.awsSandboxAllowedDestinations {
		sandboxEmails[e] = struct{}{}
	}

	for _, de := range req.Destination.All() {
		if _, ok := sandboxEmails[de]; !ok {
			return &model.SESError{Code: "MessageRejected", Message: "Cannot send emails outside sandbox"}
		}
	}

	return nil
}
