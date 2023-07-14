package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/joeshaw/envdecode"
)

type Config struct {
	Server   ServerConfig   `env:""`
	Database DatabaseConfig `env:""`
	Slack    SlackConfig    `env:""`
	SMS      SMSConfig      `env:""`
	Email    MailConfig     `env:""`
	Cron     CronConfig     `env:""`
}

func NewConfig() (*Config, error) {
	var config Config

	err := envdecode.Decode(&config)
	if err != nil {
		if err != envdecode.ErrNoTargetFieldsAreSet {
			return nil, err
		}
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) validate() error {
	return validator.New().Struct(c)
}
