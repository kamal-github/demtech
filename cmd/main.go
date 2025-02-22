package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kamal-github/demtech/internal/api"
	"github.com/kamal-github/demtech/internal/config"
	"github.com/kamal-github/demtech/internal/repo"
	"github.com/kamal-github/demtech/internal/service"
	"github.com/kamal-github/demtech/internal/validator"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load environment variables from .env file (if exists)
	env, err := config.Process()
	if err != nil {
		log.Fatal("failed to load env config")
	}

	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Initialize router
	router := gin.New()

	// Middleware stack
	router.Use(
		gin.Logger(),   // Logs requests
		gin.Recovery(), // Recovers from panics
	)

	// Initialize services
	redisCli := redis.NewClient(&redis.Options{
		Addr: env.RedisAddr,
	})
	// Check Redis connection
	ctx := context.Background()
	if err := redisCli.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	sentEmailTracker := repo.NewRedisEmailTracker(redisCli, env.TrackingHoursForEmailsQuota)

	validators := []service.Validator{
		validator.NewEmailValidator(),
		validator.NewMaxBodySizeValidator(env.AWSMaxEmailSizeAllowedBytes),
		validator.NewMaxDestinationsValidator(env.AWSMaxDestinations),
		validator.NewSandboxValidator(env.AWSIsSandBox, env.AWSSandboxAllowedDestinations),
		validator.NewVerifiedEmailValidator(env.AWSVerifiedSourceEmailIDs),
		validator.NewQuotaValidator(sentEmailTracker, env.AWSEmailsQuotaForLastNHours),
	}

	emailService := service.NewEmailService(validators, sentEmailTracker, service.FailureConfig{FailRandomly: env.FailRandomly, FailPercentage: env.FailPercentage})
	emailStatsRepo := repo.NewEmailStatsRepo(redisCli)

	// Decorate the core email service with stats therefore taking `emailService` as downstream.
	emailStatsService := service.NewEmailStatsService(emailService, emailStatsRepo, emailStatsRepo)
	emailHandler := api.NewEmailHandler(emailStatsService, emailStatsRepo)
	emailStatsHandler := api.NewEmailStatsHandler(emailStatsService)

	// Register routes
	api := router.Group("/api/v1")
	api.POST("/send-email", emailHandler.SendEmailHandler)
	api.GET("/email-stats", emailStatsHandler.GetEmailStats)

	// Server settings
	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// Graceful shutdown handling
	go func() {
		log.Println("Starting server on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for termination signal (SIGINT, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create context with timeout for cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
