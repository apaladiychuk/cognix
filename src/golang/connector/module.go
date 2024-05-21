package main

import (
	"cognix.ch/api/v2/core/ai"
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/storage"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-resty/resty/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

type Config struct {
	OAuthURL string `env:"OAUTH_URL,required"`
}

var Module = fx.Options(
	repository.DatabaseModule,
	messaging.NatsModule,
	storage.MilvusModule,
	storage.MinioModule,
	ai.ChunkingModule,
	fx.Provide(
		func() (*Config, error) {
			cfg := Config{}
			err := utils.ReadConfig(&cfg)
			if err != nil {
				zap.S().Errorf(err.Error())
				return nil, err
			}
			return &cfg, nil
		},
		newOauthClient,
		repository.NewConnectorRepository,
		repository.NewCredentialRepository,
		repository.NewDocumentRepository,
		repository.NewEmbeddingModelRepository,
		NewExecutor,
	),
	fx.Invoke(RunServer),
)

func newOauthClient(cfg *Config) *resty.Client {
	return resty.New().
		SetTimeout(time.Minute).
		SetBaseURL(cfg.OAuthURL)
}
func RunServer(lc fx.Lifecycle, executor *Executor) error {

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go executor.run(context.Background(), model.TopicExecutor, model.SubscriptionExecutor, executor.runConnector)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			executor.msgClient.Close()
			return nil
		},
	})
	return nil
}
