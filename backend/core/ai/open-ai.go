package ai

import (
	"cognix.ch/api/v2/core/model"
	"context"
	openai "github.com/sashabaranov/go-openai"
)

type (
	Response struct {
		Message string
	}
	OpenAIClient interface {
		Request(ctx context.Context, message string) (*Response, error)
	}

	openAIClient struct {
		client  *openai.Client
		modelID string
	}
)

func (o *openAIClient) Request(ctx context.Context, message string) (*Response, error) {
	userMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	}
	resp, err := o.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    o.modelID,
			Messages: []openai.ChatCompletionMessage{userMessage},
		},
	)
	o.client.CreateEmbeddings(ctx, &openai.EmbeddingRequest{
		Input:          nil,
		Model:          "",
		User:           "",
		EncodingFormat: "",
		Dimensions:     0,
	})
	if err != nil {
		return nil, err
	}
	response := &Response{Message: resp.Choices[0].Message.Content}
	return response, nil
}

func NewOpenAIClient(llm *model.LLM) OpenAIClient {

	return &openAIClient{
		client:  openai.NewClient(llm.ApiKey),
		modelID: llm.ModelID,
	}
}
