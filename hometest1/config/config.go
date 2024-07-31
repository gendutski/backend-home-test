package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HttpPort string `envconfig:"HTTP_PORT" default:"8080"`
}

func Get() Config {
	var cfg Config
	envconfig.MustProcess("", &cfg)
	return cfg
}
