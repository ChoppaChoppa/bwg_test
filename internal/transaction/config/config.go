package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Server struct {
		Host string `envconfig:"SERVER_HOST" default:"127.0.0.1:9000"`
	}

	DataBase struct {
		Login    string `envconfig:"DB_LOGIN" default:"admin"`
		Password string `envconfig:"DB_PASSWORD" default:"admin"`
		DBName   string `envconfig:"DB_NAME" default:"postgres"`
		SslMode  string `envconfig:"DB_SSL_NAME" default:"disable"`
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
