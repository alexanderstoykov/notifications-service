package config

type SMSConfig struct {
	SID    string `env:"SMS_SID" validate:"required"`
	Token  string `env:"SMS_TOKEN" validate:"required"`
	Number string `env:"SMS_NUMBER" validate:"required"`
}
