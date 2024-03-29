package main

import (
	"cognix.ch/api/v2/api/handler"
	"cognix.ch/api/v2/core/oauth"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/server"
	"github.com/caarlos0/env/v10"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type MountParams struct {
	fx.In
	Router            *gin.Engine
	AuthMiddleware    *server.AuthMiddleware
	AutHandler        *handler.AuthHandler
	SwaggerHandler    *handler.SwaggerHandler
	ConnectorHandler  *handler.ConnectorHandler
	CredentialHandler *handler.CredentialHandler
}

type Config struct {
	DB             *repository.Config
	OAuth          *oauth.Config
	Debug          bool   `env:"DEBUG" envDefault:"false"`
	Port           int    `env:"PORT" envDefault:"8080"`
	RedirectURL    string `env:"REDIRECT_URL"`
	JWTSecret      string `env:"JWT_SECRET" envDefault:"secret"`
	JWTExpiredTime int    `env:"JWT_EXPIRED_TIME" envDefault:"60"`
}

func ReadConfig() (*Config, error) {
	cfg := &Config{
		DB:    &repository.Config{},
		OAuth: &oauth.Config{},
	}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
