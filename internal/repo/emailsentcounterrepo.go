package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisEmailTracker uses Redis to track emails in a n-hour window
type RedisEmailTracker struct {
	client                      *redis.Client
	trackingHoursForEmailsQuota time.Duration // 24 hrs for ex.
}

func NewRedisEmailTracker(client *redis.Client, trackingHoursForEmailsQuota time.Duration) *RedisEmailTracker {

	return &RedisEmailTracker{client: client, trackingHoursForEmailsQuota: trackingHoursForEmailsQuota}
}

// TrackSentEmail tracks a new email sent with timestamp in order to be able to get the emails AWS service quota later on.
func (c *RedisEmailTracker) TrackSentEmail(ctx context.Context, msgID string) error {
	// Use sorted set with timestamp as score
	now := float64(time.Now().Unix())
	if err := c.client.ZAdd(ctx, "email_sends", redis.Z{
		Score:  now,
		Member: msgID,
	}).Err(); err != nil {
		return err
	}

	return c.Cleanup(ctx)
}

// GetLastNHoursCount returns count of emails in last N hours
func (c *RedisEmailTracker) GetLastNHoursCount(ctx context.Context) (int64, error) {
	min := float64(time.Now().Add(-c.trackingHoursForEmailsQuota * time.Hour).Unix())
	max := float64(time.Now().Unix())
	return c.client.ZCount(ctx, "email_sends", fmt.Sprintf("%f", min), fmt.Sprintf("%f", max)).Result()
}

// Cleanup removes entries older than N hours periodically
func (c *RedisEmailTracker) Cleanup(ctx context.Context) error {
	min := float64(0)
	max := float64(time.Now().Add(-c.trackingHoursForEmailsQuota * time.Hour).Unix())
	return c.client.ZRemRangeByScore(ctx, "email_sends", fmt.Sprintf("%f", min), fmt.Sprintf("%f", max)).Err()
}
