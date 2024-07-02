package ai

import (
	"cognix.ch/api/v2/core/storage"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type SearcherResponse struct {
	DocumentID int64  `json:"document_id,omitempty"`
	Content    string `json:"content,omitempty"`
}
type Searcher interface {
	FindDocuments(ctx context.Context, userID, tenantID uuid.UUID,
		embeddingModel string,
		message string,
		collectionNames ...string) ([]*SearcherResponse, error)
}

func NewSearcher(
	searcherType string,
	embeddBuilder *EmbeddingBuilder,
	vectorDB storage.VectorDBClient,
	embeddGRPCBuilder *GRPCEmbeddingBuilder,
) (Searcher, error) {
	switch searcherType {
	case VectorSearchInternal:
		return &InternalSearcher{
			embeddBuilder: embeddBuilder,
			vectorDB:      vectorDB,
		}, nil
	case VectorSearchGRPCService:
		return &SearcherGRPC{
			embeddBuilder: embeddGRPCBuilder,
		}, nil
	}
	return nil, fmt.Errorf("vector searcher %s not implemented", searcherType)
}
