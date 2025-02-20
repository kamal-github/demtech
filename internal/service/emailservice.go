package service

import (
	"context"
	"math/rand"

	"github.com/google/uuid"
	"github.com/kamal-github/demtech/internal/model"
)

type Validator interface {
	Validate(ctx context.Context, req model.EmailRequest) error
}

type SentEmailTracker interface {
	TrackSentEmail(ctx context.Context, msgID string) error
}

type EmailService struct {
	validators       []Validator
	sentEmailTracker SentEmailTracker
}

func NewEmailService(validators []Validator, sentEmailTracker SentEmailTracker) EmailService {
	return EmailService{validators: validators, sentEmailTracker: sentEmailTracker}
}

func (es EmailService) SendEmail(ctx context.Context, req model.EmailRequest) (*model.SESResponse, error) {
	for _, v := range es.validators {
		if err := v.Validate(ctx, req); err != nil {
			return nil, err
		}
	}

	if failure := randomFailure(); failure != nil {
		return nil, failure
	}

	msgID := generateMessageID()
	if err := es.sentEmailTracker.TrackSentEmail(ctx, msgID); err != nil {
		return nil, &model.SESError{Code: "InternalFailure", Message: "Unexpected internal error occurred."}
	}

	return &model.SESResponse{MessageID: msgID}, nil
}

func randomFailure() *model.SESError {
	failures := []model.SESError{
		{Code: "InternalFailure", Message: "Unexpected internal error occurred."},
		{Code: "ThrottlingException", Message: "Rate limit exceeded."},
		{Code: "AccountSendingPaused", Message: "Email sending is disabled for your account."},
	}
	if rand.Intn(100) < 30 { // 30% failure chance
		return &failures[rand.Intn(len(failures))]
	}
	return nil
}

func generateMessageID() string {
	return uuid.NewString()
}
