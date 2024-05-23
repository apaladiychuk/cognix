package main

import (
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"go.uber.org/fx"
)

type Config struct {
	RenewInterval int `env:"RENEW_INTERVAL" envDefault:"30"`
}

var Module = fx.Options(
	repository.DatabaseModule,
	messaging.NatsModule,
	fx.Provide(
		func() (*Config, error) {
			cfg := Config{}
			err := utils.ReadConfig(&cfg)
			return &cfg, err
		},
		repository.NewConnectorRepository,
		NewCronTrigger,
		NewServer,
	),
	fx.Invoke(RunServer),
)

func RunServer(lc fx.Lifecycle, server *Server) error {

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return server.run(ctx)
		},
		OnStop: func(ctx context.Context) error {

			return nil
		},
	})
	return nil
}
