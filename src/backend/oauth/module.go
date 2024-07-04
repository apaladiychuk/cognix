package main

import (
	"cognix.ch/api/v2/core/oauth"
	"cognix.ch/api/v2/core/server"
	"cognix.ch/api/v2/core/utils"
	"cognix.ch/api/v2/oauth/handler"
	"fmt"
	"github.com/caarlos0/env/v10"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/fx"
	"net/http"
)

var Module = fx.Options(
	fx.Provide(ReadConfig,
		server.NewRouter,
		newOauthHandler,
	),
	fx.Invoke(
		MountRoute,
		RunServer,
	),
)

type Config struct {
	OAuth *oauth.Config
	Port  int `env:"PORT" envDefault:"8080"`
}

type MountParams struct {
	fx.In
	Router       *gin.Engine
	OAuthHandler *handler.OAuthHandler
}

func MountRoute(param MountParams) error {
	param.OAuthHandler.Mount(param.Router)
	return nil
}

func ReadConfig() (*Config, error) {
	cfg := &Config{
		OAuth: &oauth.Config{
			Microsoft: &oauth.MicrosoftConfig{},
			Google:    &oauth.GoogleConfig{},
		},
	}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	utils.InitLogger(false)
	return cfg, nil
}

func newOauthHandler(cfg *Config) *handler.OAuthHandler {
	return handler.NewOAuthHandler(cfg.OAuth)
}

func RunServer(cfg *Config, router *gin.Engine) {
	srv := http.Server{}
	srv.Addr = fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	srv.Handler = router
	otelzap.S().Infof("Start server %s ", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		otelzap.S().Errorf("HTTP server: %s", err.Error())
	}
}
