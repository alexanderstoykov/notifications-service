package config

// MailConfig holds configuration for Mail config.
type MailConfig struct {
	EmailSender  string `env:"EMAIL_SENDER" validate:"required"`
	SMTPHost     string `env:"EMAIL_SMTP_HOST" validate:"required"`
	SMTPPort     int    `env:"EMAIL_SMTP_PORT,default=456" validate:"required"`
	SMTPUsername string `env:"EMAIL_SMTP_USERNAME" validate:"required"`
	SMTPPassword string `env:"EMAIL_SMTP_PASSWORD" validate:"required"`
}
