package ai

import (
	_ "github.com/deluan/flowllm/llms/openai"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	EmbeddingConfig struct {
		EmbeddingURL string `env:"EMBEDDING_GRPC_URL"`
	}
)

func (v EmbeddingConfig) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.EmbeddingURL, validation.Required))
}
