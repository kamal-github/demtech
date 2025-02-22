//go:build integration

package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/kamal-github/demtech/internal/repo"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisEmailTracker(t *testing.T) {
	// Set up a real Redis client (ensure Redis is running)
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Use a short duration for testability
	trackingDuration := 1 * time.Hour
	tracker := repo.NewRedisEmailTracker(client, trackingDuration)

	ctx := context.Background()
	defer client.FlushDB(ctx) // Clean up Redis after test

	// Test TrackSentEmail
	msgID := "test-email-123"
	err := tracker.TrackSentEmail(ctx, msgID)
	assert.NoError(t, err, "TrackSentEmail should not return an error")

	// Test GetLastNHoursCount
	count, err := tracker.GetLastNHoursCount(ctx)
	assert.NoError(t, err, "GetLastNHoursCount should not return an error")
	assert.Equal(t, int64(1), count, "Expected email count to be 1")

	// Test Cleanup (simulate waiting for expiry)
	time.Sleep(2 * time.Second)
	err = tracker.Cleanup(ctx)
	assert.NoError(t, err, "Cleanup should not return an error")

	countAfterCleanup, err := tracker.GetLastNHoursCount(ctx)
	assert.NoError(t, err, "GetLastNHoursCount after cleanup should not return an error")
	assert.Equal(t, int64(1), countAfterCleanup, "Email should still be there within the tracking window")

	// Simulate email expiry (older than tracking duration)
	client.ZAdd(ctx, "test-email-sends", redis.Z{
		Score:  float64(time.Now().Add(-2 * trackingDuration).Unix()), // Add an old email
		Member: "old-email",
	})

	err = tracker.Cleanup(ctx)
	assert.NoError(t, err, "Cleanup should not return an error")

	finalCount, err := tracker.GetLastNHoursCount(ctx)
	assert.NoError(t, err, "GetLastNHoursCount final check should not return an error")
	assert.Equal(t, int64(1), finalCount, "Old emails should have been cleaned up")
}
