package config

import "time"

type CronConfig struct {
	Interval  time.Duration `env:"CRON_INTERVAL" validate:"required"`
	BatchSize int           `env:"CRON_BATCH_SIZE" validate:"required"`
}
