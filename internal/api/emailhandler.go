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

// EmailHandler struct
type EmailHandler struct {
	service EmailService
}

// NewEmailHandler creates a new EmailHandler
func NewEmailHandler(service EmailService) *EmailHandler {
	return &EmailHandler{service: service}
}

// SendEmailHandler handles sending emails
func (h *EmailHandler) SendEmailHandler(c *gin.Context) {
	var emailReq model.EmailRequest

	if err := c.ShouldBindJSON(&emailReq); err != nil {
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
