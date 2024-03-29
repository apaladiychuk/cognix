package repository

import (
	"cognix.ch/api/v2/core/utils"
	"go.uber.org/fx"
)

var DatabaseModule = fx.Options(
	fx.Provide(func() (*Config, error) {
		cfg := Config{}
		err := utils.ReadConfig(&cfg)
		return &cfg, err
	},
		NewDatabase,
		NewUserRepository,
		NewCredentialRepository,
		NewConnectorRepository,
		NewLLMRepository,
		NewPersonaRepository,
	),
)
