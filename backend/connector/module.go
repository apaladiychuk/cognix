package main

import (
	"cognix.ch/api/v2/core/ai"
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/storage"
	"context"
	"go.uber.org/fx"
)

var Module = fx.Options(
	repository.DatabaseModule,
	messaging.NatsModule,
	storage.MilvusModule,
	storage.MinioModule,
	fx.Provide(
		repository.NewConnectorRepository,
		repository.NewDocumentRepository,
		repository.NewEmbeddingModelRepository,
		ai.NewEmbeddingParser,
		NewExecutor,
	),
	fx.Invoke(RunServer),
)

func RunServer(lc fx.Lifecycle, executor *Executor) error {

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := executor.run(context.Background(), model.TopicEmbedding, model.SubscriptionEmbedding, executor.runEmbedding); err != nil {
				return err
			}
			return executor.run(context.Background(), model.TopicExecutor, model.SubscriptionExecutor, executor.runConnector)
		},
		OnStop: func(ctx context.Context) error {
			executor.msgClient.Close()
			return nil
		},
	})
	return nil
}
