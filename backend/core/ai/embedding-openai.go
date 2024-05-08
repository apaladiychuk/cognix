package ai

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/storage"
	"context"
	_ "github.com/deluan/flowllm/llms/openai"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type (
	embeddingParser struct {
		embeddingModel *model.EmbeddingModel
		client         *openai.Client
		embeddingRepo  repository.EmbeddingModelRepository
		minioClient    storage.MinIOClient
	}
)

func NewEmbeddingParser(client *openai.Client,
	embeddingRepo repository.EmbeddingModelRepository,
	minioClient storage.MinIOClient) EmbeddingParser {

	return &embeddingParser{
		embeddingModel: &model.EmbeddingModel{ModelID: "text-embedding-ada-002"},
		client:         client,
		embeddingRepo:  embeddingRepo,
		minioClient:    minioClient,
	}
}

func (p *embeddingParser) Parse(ctx context.Context, payload *proto.EmbeddingRequest) (*proto.EmbeddingResponse, error) {
	embeddingRequest := openai.EmbeddingRequestStrings{
		Input:          payload.GetContent(),
		Model:          openai.EmbeddingModel(payload.GetEmbeddingModel()),
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
	for _, em := range embeddingResponse.Data {
		zap.S().Infof("%d, %s", em.Index, em.Object)
		response.Payload = append(response.Payload, &proto.EmbeddingResponse_Payload{
			Chunk:   int64(em.Index),
			Content: em.Object,
			Vector:  em.Embedding,
		})

	}
	return response, nil
}
