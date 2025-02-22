package validator

import (
	"context"

	"github.com/kamal-github/demtech/internal/model"
)

type SandboxValidator struct {
	awsIsSandbox                  bool
	awsSandboxAllowedDestinations []string
}

func NewSandboxValidator(s bool, e []string) SandboxValidator {
	return SandboxValidator{awsIsSandbox: s, awsSandboxAllowedDestinations: e}
}

func (v SandboxValidator) Validate(ctx context.Context, req model.EmailRequest) error {
	if !v.awsIsSandbox {
		return nil
	}

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
