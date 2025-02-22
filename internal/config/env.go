package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	RedisAddr string `envconfig:"REDIS_ADDR"`
	// Configuration for tracking sent emails count for last N hours - say last 24 hours as per AWS.
	// but kept it configurable so as to test it realistically.
	TrackingHoursForEmailsQuota   time.Duration `envconfig:"TRACKING_HOURS_FOR_EMAILS_QUOTA"`
	AWSMaxEmailSizeAllowedBytes   int64         `envconfig:"AWS_MAX_EMAIL_SIZE_ALLOWED_BYTES"`
	AWSMaxDestinations            int           `envconfig:"AWS_MAX_DESTINATIONS"`
	AWSIsSandBox                  bool          `envconfig:"AWS_IS_SANDBOX"`
	AWSSandboxAllowedDestinations []string      `envconfig:"AWS_SANDBOX_ALLOWED_DESTINATIONS"`
	AWSVerifiedSourceEmailIDs     []string      `envconfig:"AWS_VERIFIED_SOURCE_EMAIL_IDS"`
	AWSEmailsQuotaForLastNHours   int64         `envconfig:"AWS_EMAILS_QUOTA_FOR_LAST_N_HOURS"`
	FailRandomly                  bool          `envconfig:"FAIL_RANDOMLY"`
	FailPercentage                int           `envconfig:"FAIL_PERCENTAGE"`
}

func Process() (Env, error) {
	var e Env
	if err := envconfig.Process("", &e); err != nil {
		return Env{}, err

	}

	return e, nil
}
