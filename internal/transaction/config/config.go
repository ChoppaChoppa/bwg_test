package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Server struct {
		Host string `envconfig:"SERVER_HOST" default:"127.0.0.1:9000"`
	}
}

func Parse() (*Config, error) {
	var cfg = &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
