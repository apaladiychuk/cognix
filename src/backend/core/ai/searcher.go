package ai

import (
	"cognix.ch/api/v2/core/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"time"
)

// VectorSearchConfig is a configuration struct for vector search module.
//
// It contains the configuration options for vector search API,
// including the endpoint for vector search service.
type VectorSearchConfig struct {
	VectorSearchHost string `env:"VECTOR_SEARCHER_GRPC_HOST,required"`
	VectorSearchPort int    `env:"VECTOR_SEARCHER_GRPC_PORT,required"`
	ApiVectorSearch  string `env:"API-VECTOR-SEARCH" envDefault:"INTERNAL"`
}

type GRPCEmbeddingBuilder struct {
	cfg    *VectorSearchConfig
	client proto.SearchServiceClient
	mx     sync.Mutex
}

func NewGRPCEmbeddingBuilder(cfg *VectorSearchConfig) *GRPCEmbeddingBuilder {
	return &GRPCEmbeddingBuilder{
		cfg: cfg,
		mx:  sync.Mutex{},
	}
}
func (e *GRPCEmbeddingBuilder) Client() (proto.SearchServiceClient, error) {
	e.mx.Lock()
	defer e.mx.Unlock()
	if e.client == nil {
		dialOptions := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials())}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", e.cfg.VectorSearchHost, e.cfg.VectorSearchPort), dialOptions...)
		if err != nil {
			return nil, err
		}
		e.client = proto.NewSearchServiceClient(conn)
	}
	return e.client, nil
}
