package main

import (
	"cognix.ch/api/v2/core/repository"
	"context"
	"go.uber.org/fx"
)

var Module = fx.Options(fx.Provide(
	repository.DatabaseModule,
	repository.NewConnectorRepository,
	NewExecutor,
),
	fx.Invoke(RunServer),
)

func RunServer(lc fx.Lifecycle, executor *executor) error {

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return executor.run(ctx)
		},
		OnStop: func(ctx context.Context) error {

			return nil
		},
	})
	return nil
}
