package ai

import (
	"cognix.ch/api/v2/core/proto"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SearcherGrpc struct {
}

type SearcherGRPC struct {
	embeddBuilder *GRPCEmbeddingBuilder
}

func (i *SearcherGRPC) FindDocuments(ctx context.Context, userID, tenantID uuid.UUID,
	embeddingModel string,
	message string, collectionNames ...string) ([]*SearcherResponse, error) {
	embedding, err := i.embeddBuilder.Client()
	if err != nil {
		return nil, err
	}
	response, err := embedding.VectorSearch(ctx, &proto.SearchRequest{
		Content: message,
		Model:   embeddingModel,
	})
	if err != nil {
		zap.S().Errorf("embeding service %s ", err.Error())
		return nil, err
	}
	var result []*SearcherResponse

	for _, docID := range response.GetVector() {
		resDocument := &SearcherResponse{
			DocumentID: docID,
		}
		result = append(result, resDocument)
	}
	return result, nil
}
