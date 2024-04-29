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
		repository.NewDocumentRepository,
		NewExecutor,
	),
	fx.Invoke(RunServer),
)

func RunServer(lc fx.Lifecycle, executor *executor) error {

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return executor.run(context.Background())
		},
		OnStop: func(ctx context.Context) error {

			return nil
		},
	})
	return nil
}
