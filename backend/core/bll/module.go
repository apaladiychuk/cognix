package bll

import (
	"cognix.ch/api/v2/core/utils"
	"go.uber.org/fx"
)

type Config struct {
	RedirectURL string `env:"REDIRECT_URL"`
}

var BLLModule = fx.Options(
	fx.Provide(func() (*Config, error) {
		cfg := Config{}
		err := utils.ReadConfig(&cfg)
		return &cfg, err
	}),
	fx.Provide(
		NewConnectorBL,
		NewCredentialBL,
		NewAuthBL,
		NewPersonaBL,
		NewChatBL,
		NewDocumentBL,
		NewDocumentSetBL,
		NewEmbeddingModelBL,
	),
)
