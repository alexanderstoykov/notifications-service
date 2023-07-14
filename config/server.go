package config

import (
	"fmt"
)

type ServerConfig struct {
	Host            string `env:"SERVER_HOST"`
	Port            int    `env:"SERVER_PORT"`
	ShutdownTimeout int    `env:"SERVER_SHUTDOWN_TIMEOUT"`
}

func (s ServerConfig) PortString() string {
	return fmt.Sprintf(":%d", s.Port)
}

func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
