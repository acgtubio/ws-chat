package config

import "github.com/kelseyhightower/envconfig"

type PostgresConfig struct {
	PostgresHost     string `envconfig:"POSTGRES_HOST"`
	PostgresUsername string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	PostgresDBName   string `envconfig:"POSTGRES_DB" default:"chat"`
}

type ApplicationConfig struct {
	Port int16 `envconfig:"APP_PORT" default:"8080"`
}

type Config struct {
	Postgres    PostgresConfig
	Application ApplicationConfig
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("", &config.Postgres); err != nil {
		return Config{}, err
	}

	if err := envconfig.Process("", &config.Application); err != nil {
		return Config{}, err
	}

	return config, nil
}
