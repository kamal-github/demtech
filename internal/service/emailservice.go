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

type FailureConfig struct {
	FailRandomly   bool
	FailPercentage int
}

type EmailServiceImpl struct {
	validators       []Validator
	sentEmailTracker SentEmailTracker
	failureConfig    FailureConfig
}

func NewEmailService(validators []Validator, sentEmailTracker SentEmailTracker, cfg FailureConfig) EmailServiceImpl {
	return EmailServiceImpl{validators: validators, sentEmailTracker: sentEmailTracker, failureConfig: cfg}
}

func (es EmailServiceImpl) SendEmail(ctx context.Context, req model.EmailRequest) (*model.SESResponse, error) {
	for _, v := range es.validators {
		if err := v.Validate(ctx, req); err != nil {
			return nil, err
		}
	}

	if failure := es.randomFailure(); failure != nil {
		return nil, failure
	}

	msgID := generateMessageID()

	// For every message that you send, the total number of recipients
	// (including each recipient in the To:, CC: and BCC: fields) is counted
	// against the maximum number of emails you can send in a 24-hour period
	// (your sending quota), therefore we have to track each email being sent
	// to all the receipients.
	// In short - "One email is multiplexed to many recipients".
	for _, dest := range req.Destination.All() {
		if err := es.sentEmailTracker.TrackSentEmail(ctx, msgID+"-"+dest); err != nil {
			return nil, &model.SESError{Code: "InternalFailure", Message: "Unexpected internal error occurred."}
		}
	}

	return &model.SESResponse{MessageID: msgID}, nil
}

func (es EmailServiceImpl) randomFailure() *model.SESError {
	if !es.failureConfig.FailRandomly {
		return nil
	}

	failures := []model.SESError{
		{Code: "InternalFailure", Message: "Unexpected internal error occurred."},
		{Code: "ThrottlingException", Message: "Rate limit exceeded."},
		{Code: "AccountSendingPaused", Message: "Email sending is disabled for your account."},
		{Code: "AccessDeniedException", Message: "Access denied."},
		{Code: "InvalidClientTokenId", Message: "Invalid client token."},
		{Code: "SignatureDoesNotMatch", Message: "Signature does not match."},
		{Code: "RequestExpired", Message: "Request Expired."},
		{Code: "ThrottlingException", Message: "Request was throttled due to exceeding AWS SES limits"},
		{Code: "TooManyRequestsException", Message: "Too many requests."},
		{Code: "MessageRejected", Message: "Message rejected."},
	}

	if es.failureConfig.FailPercentage < 0 || es.failureConfig.FailPercentage > 100 {
		es.failureConfig.FailPercentage = 20
	}
	if rand.Intn(100) < es.failureConfig.FailPercentage {
		return &failures[rand.Intn(len(failures))]
	}
	return nil
}

func generateMessageID() string {
	return uuid.NewString()
}
