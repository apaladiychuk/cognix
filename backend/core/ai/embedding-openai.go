package ai

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"context"
	_ "github.com/deluan/flowllm/llms/openai"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"os"
)

type (
	embeddingParser struct {
		embeddingModel *model.EmbeddingModel
		client         *openai.Client
	}
)

func NewEmbeddingParser(embeddingModel *model.EmbeddingModel) EmbeddingParser {
	//remove-it
	apiKey := os.Getenv("OPENAI_API_KEY")
	return &embeddingParser{
		embeddingModel: embeddingModel,
		client:         openai.NewClient(apiKey),
	}
}

func (p *embeddingParser) Parse(ctx context.Context, payload *proto.EmbeddingRequest) (*proto.EmbeddingResponse, error) {
	embeddingModel := payload.GetEmbeddingModel()
	if embeddingModel == "" {
		embeddingModel = p.embeddingModel.ModelID
	}
	embeddingRequest := openai.EmbeddingRequestStrings{
		Input:          payload.GetContent(),
		Model:          openai.EmbeddingModel(embeddingModel),
		EncodingFormat: openai.EmbeddingEncodingFormatFloat,
		Dimensions:     0,
	}

	embeddingResponse, err := p.client.CreateEmbeddings(ctx, embeddingRequest)
	if err != nil {
		return nil, err
	}
	response := &proto.EmbeddingResponse{
		DocumentId: payload.GetDocumentId(),
		Key:        payload.Key,
	}
	for i, em := range embeddingResponse.Data {
		zap.S().Infof("%d, %s", em.Index, em.Object)
		response.Payload = append(response.Payload, &proto.EmbeddingResponse_Payload{
			Chunk:   int64(em.Index),
			Content: payload.Content[i],
			Vector:  em.Embedding,
		})

	}
	return response, nil
}
