package main

import (
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/repository"
	"context"
	"go.uber.org/fx"
)

var Module = fx.Options(
	repository.DatabaseModule,
	messaging.NatsModule,
	fx.Provide(
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
