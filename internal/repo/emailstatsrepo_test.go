//go:build integration

package repo_test

import (
	"context"
	"testing"

	"github.com/kamal-github/demtech/internal/repo"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func TestEmailStatsRepoImpl_Integration(t *testing.T) {
	redisClient := setupRedisClient()
	defer redisClient.Close()

	ctx := context.Background()
	repo := repo.NewEmailStatsRepo(redisClient)

	// Test IncrementSuccess
	err := repo.IncrementSuccess(ctx)
	assert.NoError(t, err)

	// Test IncrementError with a specific error type
	err = repo.IncrementError(ctx, "Timeout")
	assert.NoError(t, err)

	// Fetch the stats and validate
	stats, err := repo.GetEmailStats(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, stats.SuccessCount)
	assert.Equal(t, 1, stats.TotalErrCount)
	assert.Equal(t, 2, stats.TotalEmailsSent)
	assert.Equal(t, 1, stats.Errors["Timeout"])
}
