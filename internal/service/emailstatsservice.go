package service

import (
	"context"
	"errors"

	"github.com/kamal-github/demtech/internal/model"
)

const (
	unknowErrTypeKey = "unknown"
)

type EmailService interface {
	SendEmail(ctx context.Context, req model.EmailRequest) (*model.SESResponse, error)
}

type EmailsStatsUpdater interface {
	IncrementSuccess(ctx context.Context) error
	IncrementError(ctx context.Context, errorType string) error
}

type EmailsStatsGetter interface {
	GetEmailStats(ctx context.Context) (model.EmailStats, error)
}

type EmailStatsService struct {
	emailService      EmailService
	emailStatsUpdater EmailsStatsUpdater
	emailStatsGetter  EmailsStatsGetter
}

func NewEmailStatsService(s EmailService, u EmailsStatsUpdater, g EmailsStatsGetter) EmailStatsService {
	return EmailStatsService{emailService: s, emailStatsUpdater: u, emailStatsGetter: g}
}

func (es EmailStatsService) SendEmail(ctx context.Context, req model.EmailRequest) (*model.SESResponse, error) {
	var (
		res *model.SESResponse
		err error
	)

	if res, err = es.emailService.SendEmail(ctx, req); err != nil {
		var sesErr *model.SESError
		if errors.As(err, &sesErr) {
			// as the error from SendEmail is expected and should be returned
			// it is ok, if we could not track the error.
			es.emailStatsUpdater.IncrementError(ctx, sesErr.Code)
			return nil, err
		}
		es.emailStatsUpdater.IncrementError(ctx, unknowErrTypeKey)
		return nil, err
	}

	es.emailStatsUpdater.IncrementSuccess(ctx)

	return res, nil
}

func (es EmailStatsService) GetEmailStats(ctx context.Context) (model.EmailStats, error) {
	return es.emailStatsGetter.GetEmailStats(ctx)
}
