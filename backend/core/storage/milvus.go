package storage

import (
	"cognix.ch/api/v2/core/utils"
	"context"
	milvus "github.com/milvus-io/milvus-sdk-go/v2/client"
	"go.uber.org/fx"
)

type (
	Config struct {
		Address string `env:"MILVUS_URL"`
	}
	MilvusClient interface {
	}
	milvusClient struct {
		client milvus.Client
	}
)

var MilvusModule = fx.Options(
	fx.Provide(func() (*Config, error) {
		cfg := Config{}
		err := utils.ReadConfig(&cfg)
		return &cfg, err
	},
		NewMilvusClient,
	),
)

func NewMilvusClient(cfg *Config) (MilvusClient, error) {
	client, err := milvus.NewClient(context.Background(), milvus.Config{
		Address: cfg.Address,
	})
	if err != nil {
		return nil, err
	}
	return milvusClient{
		client: client,
	}, nil
}
