package ai

import (
	_ "github.com/deluan/flowllm/llms/openai"
)

// EmbeddingConfig is a configuration struct for embedding module.
//
// It contains the configuration options for connecting to the embedding server
// over gRPC.
type EmbeddingConfig struct {
	EmbedderHost string `env:"EMBEDDER_GRPC_HOST,required"`
	EmbedderPort int    `env:"EMBEDDER_GRPC_PORT,required"`
}
