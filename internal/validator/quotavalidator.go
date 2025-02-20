package validator

import (
	"context"

	"github.com/kamal-github/demtech/internal/model"
)

type LastNHoursCountGetter interface {
	GetLastNHoursCount(ctx context.Context) (int64, error)
}

type QuotaValidator struct {
	lastNHoursCountGetter       LastNHoursCountGetter
	awsEmailsQuotaForLastNHours int64
}

func NewQuotaValidator(g LastNHoursCountGetter, eq int64) QuotaValidator {
	return QuotaValidator{lastNHoursCountGetter: g, awsEmailsQuotaForLastNHours: eq}
}

func (v QuotaValidator) Validate(ctx context.Context, req model.EmailRequest) error {
	emailsSent, err := v.lastNHoursCountGetter.GetLastNHoursCount(ctx)
	if err != nil {
		return err
	}

	if emailsSent >= v.awsEmailsQuotaForLastNHours {
		return &model.SESError{Code: "ThrottlingException", Message: "Sending quota exceeded"}
	}

	return nil
}
