package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamal-github/demtech/internal/model"
)

type EmailStatsService interface {
	GetEmailStats(context.Context) (model.EmailStats, error)
}

type EmailStatsHandler struct {
	emailStatsService EmailStatsService
}

func NewEmailStatsHandler(s EmailStatsService) EmailStatsHandler {
	return EmailStatsHandler{emailStatsService: s}
}

func (h EmailStatsHandler) GetEmailStats(c *gin.Context) {
	ctx := c.Request.Context()
	stats, err := h.emailStatsService.GetEmailStats(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
