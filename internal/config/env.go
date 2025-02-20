package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	RedisAddr string
	// Configuration for tracking sent emails count for last N hours - say last 24 hours as per AWS.
	// but kept it configurable so as to test it realistically.
	TrackingHoursForEmailsQuota   time.Duration
	AWSMaxEmailSizeAllowedBytes   int64
	AWSMaxDestinations            int
	AWSSandboxAllowedDestinations []string
	AWSVerifiedSourceEmailIDs     []string
	AWSEmailsQuotaForLastNHours   int64
}

func Process() (Env, error) {
	var e Env
	if err := envconfig.Process("", &e); err != nil {
		return Env{}, err
	}

	return e, nil
}
