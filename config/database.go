package config

import "fmt"

type DatabaseConfig struct {
	Host     string `env:"DATABASE_HOST" validate:"required"`
	Port     int    `env:"DATABASE_PORT" validate:"required"`
	User     string `env:"DATABASE_USER" validate:"required"`
	Password string `env:"DATABASE_PASSWORD" validate:"required"`
	Database string `env:"DATABASE_NAME" validate:"required"`
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Database,
	)
}
