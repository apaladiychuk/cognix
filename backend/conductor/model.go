package main

import (
	"cognix.ch/api/v2/core/repository"
	"github.com/caarlos0/env/v10"
)

type Config struct {
	DB    *repository.Config
	Debug bool `env:"DEBUG" envDefault:"false"`
}

func ReadConfig() (*Config, error) {
	cfg := &Config{
		DB: &repository.Config{},
	}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
