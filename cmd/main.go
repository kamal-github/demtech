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
	env := loadConfig()

	router := setupRouter()
	redisCli := setupRedis(env)

	emailStatsRepo := repo.NewEmailStatsRepo(redisCli)
	emailStatsService := setupEmailService(env, redisCli, emailStatsRepo)

	registerRoutes(router, emailStatsService, emailStatsRepo)

	server := startServer(router)
	gracefulShutdown(server)
}

// loadConfig initializes environment configuration
func loadConfig() config.Env {
	env, err := config.Process()
	if err != nil {
		log.Fatal("Failed to load environment configuration")
	}
	return env
}

// setupRouter initializes the Gin router with middleware
func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	return router
}

// setupRedis initializes Redis client and verifies connection
func setupRedis(env config.Env) *redis.Client {
	redisCli := redis.NewClient(&redis.Options{
		Addr: env.RedisAddr,
	})

	ctx := context.Background()
	if err := redisCli.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return redisCli
}

// setupEmailService initializes email service and its dependencies
func setupEmailService(env config.Env, redisCli *redis.Client, emailStatsRepo repo.EmailStatsRepoImpl) service.EmailStatsService {
	sentEmailTracker := repo.NewRedisEmailTracker(redisCli, env.TrackingHoursForEmailsQuota)

	validators := []service.Validator{
		validator.NewEmailValidator(),
		validator.NewMaxBodySizeValidator(env.AWSMaxEmailSizeAllowedBytes),
		validator.NewMaxDestinationsValidator(env.AWSMaxDestinations),
		validator.NewSandboxValidator(env.AWSIsSandBox, env.AWSSandboxAllowedDestinations),
		validator.NewVerifiedEmailValidator(env.AWSVerifiedSourceEmailIDs),
		validator.NewQuotaValidator(sentEmailTracker, env.AWSEmailsQuotaForLastNHours),
	}

	emailService := service.NewEmailService(validators, sentEmailTracker, service.FailureConfig{
		FailRandomly:   env.FailRandomly,
		FailPercentage: env.FailPercentage,
	})

	// Wrap email service with stats tracking
	return service.NewEmailStatsService(emailService, emailStatsRepo, emailStatsRepo)
}

// registerRoutes sets up API routes
func registerRoutes(router *gin.Engine, emailStatsService service.EmailStatsService, emailStatsRepo repo.EmailStatsRepoImpl) {
	apiGroup := router.Group("/api/v1")
	emailHandler := api.NewEmailHandler(emailStatsService, emailStatsRepo)
	emailStatsHandler := api.NewEmailStatsHandler(emailStatsService)

	apiGroup.POST("/send-email", emailHandler.SendEmailHandler)
	apiGroup.GET("/email-stats", emailStatsHandler.GetEmailStats)
}

// startServer initializes and starts the HTTP server
func startServer(router *gin.Engine) *http.Server {
	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	go func() {
		log.Println("Starting server on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	return server
}

// gracefulShutdown handles cleanup and graceful termination
func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
