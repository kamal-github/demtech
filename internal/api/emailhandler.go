package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kamal-github/demtech/internal/model"
)

// EmailService defines the interface for sending emails
type EmailService interface {
	SendEmail(ctx context.Context, req model.EmailRequest) (*model.SESResponse, error)
}

type EmailsStatsUpdater interface {
	IncrementSuccess(ctx context.Context) error
	IncrementError(ctx context.Context, errorType string) error
}

// EmailHandler struct
type EmailHandler struct {
	service      EmailService
	statsUpdater EmailsStatsUpdater
}

// NewEmailHandler creates a new EmailHandler
func NewEmailHandler(s EmailService, u EmailsStatsUpdater) *EmailHandler {
	return &EmailHandler{service: s, statsUpdater: u}
}

// SendEmailHandler handles sending emails
func (h *EmailHandler) SendEmailHandler(c *gin.Context) {
	var emailReq model.EmailRequest

	if err := c.ShouldBindJSON(&emailReq); err != nil {
		// Note: AWS retunrs MissingParameter or InvalidParameterValue errortype, but gin
		// returs different error, there will be a  need to map thoses errors manually.
		// Currently, keeping custom key to keep it simple.
		h.statsUpdater.IncrementError(c.Request.Context(), "MissingOrInvalidParameterValue")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.service.SendEmail(ctx, emailReq)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		if awsErr, ok := err.(*model.SESError); ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   awsErr.Code,
				"message": awsErr.Message,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Email sent successfully",
		"messageId": resp.MessageID,
	})
}
