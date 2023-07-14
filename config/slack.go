package config

type SlackConfig struct {
	Webhook string `env:"SLACK_WEBHOOK"`
}
