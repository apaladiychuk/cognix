package ai

import (
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/utils"
	"fmt"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var EmbeddingModule = fx.Options(
	fx.Provide(func() (*EmbeddingConfig, error) {
		cfg := EmbeddingConfig{}
		if err := utils.ReadConfig(&cfg); err != nil {
			return nil, err
		}
		return &cfg, nil
	},
		newEmbeddingGRPCClient),
)

func newEmbeddingGRPCClient(cfg *EmbeddingConfig) (proto.EmbedServiceClient, error) {
	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.EmbedderHost, cfg.EmbedderPort), dialOptions...)
	if err != nil {
		return nil, err
	}
	return proto.NewEmbedServiceClient(conn), nil
}
