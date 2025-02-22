package repo

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kamal-github/demtech/internal/model"
	"github.com/redis/go-redis/v9"
)

const emailStatsStorageKey = "email-stats"

type EmailStatsRepoImpl struct {
	redisClient *redis.Client
}

func NewEmailStatsRepo(c *redis.Client) EmailStatsRepoImpl {
	return EmailStatsRepoImpl{redisClient: c}
}

// Increment success count
func (r EmailStatsRepoImpl) IncrementSuccess(ctx context.Context) error {
	_, err := r.redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, emailStatsStorageKey, "totalEmailsSent", 1)
		pipe.HIncrBy(ctx, emailStatsStorageKey, "successCount", 1)
		return nil
	})
	return err
}

// Increment error count
func (r EmailStatsRepoImpl) IncrementError(ctx context.Context, errorType string) error {
	_, err := r.redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, emailStatsStorageKey, "totalEmailsSent", 1)
		pipe.HIncrBy(ctx, emailStatsStorageKey, "totalErrCount", 1)
		pipe.HIncrBy(ctx, emailStatsStorageKey, "errors:"+errorType, 1)
		return nil
	})
	return err
}

// Retrieve EmailStats from Redis
func (r EmailStatsRepoImpl) GetEmailStats(ctx context.Context) (model.EmailStats, error) {
	data, err := r.redisClient.HGetAll(ctx, emailStatsStorageKey).Result()
	if err != nil {
		return model.EmailStats{}, err
	}
	if len(data) == 0 {
		return model.EmailStats{}, fmt.Errorf("no data found for key: %s", emailStatsStorageKey)
	}

	stats := model.EmailStats{
		Errors: make(map[string]int),
	}
	for field, value := range data {
		num, err := strconv.Atoi(value)
		if err != nil {
			return model.EmailStats{}, err
		}

		switch field {
		case "totalEmailsSent":
			stats.TotalEmailsSent = num
		case "successCount":
			stats.SuccessCount = num
		case "totalErrCount":
			stats.TotalErrCount = num
		default:
			// Handle error counts (fields like "errors:Timeout", "errors:Invalid")
			if len(field) > 7 && field[:7] == "errors:" {
				errorType := field[7:] // Extract error type
				stats.Errors[errorType] = num
			}
		}
	}
	return stats, nil
}
