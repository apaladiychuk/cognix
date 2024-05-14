package responder

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/storage"
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"sync"
	"time"
)

type embedding struct {
	embedding    proto.EmbeddServiceClient
	milvusClinet storage.MilvusClient
	embeddingModel string
}

func (r *embedding) Send(ctx context.Context, ch chan *Response, wg *sync.WaitGroup, parentMessage *model.ChatMessage) {
	response, err := r.embedding.GetEmbedd(ctx, &proto.EmbeddRequest{
		Content: parentMessage.Message,
		Model:   r.embeddingModel,
	})
	if err != nil {
		ch <- &Response{
			IsValid:  false,
			Type:     ResponseError,
			Err: err,
		}
		return
	}
	docs, err := r.milvusClinet.

	for i := 0; i < 4; i++ {
		ch <- &Response{
			IsValid: true,
			Type:    ResponseDocument,
			Message: nil,
			Document: &model.DocumentResponse{
				ID:          decimal.NewFromInt(int64(i)),
				DocumentID:  "11",
				Link:        fmt.Sprintf("link for document %d", i),
				Content:     fmt.Sprintf("content of document %d", i),
				UpdatedDate: time.Now().UTC().Add(-48 * time.Hour),
				MessageID:   parentMessage.ID,
			},
		}
	}
	wg.Done()
}


func (r *embedding) FindDocuments(ctx context.Context,
		ch chan *Response, wg *sync.WaitGroup,
		parentMessage *model.ChatMessage,
		collectionNames ...string) {
	response, err := r.embedding.GetEmbedd(ctx, &proto.EmbeddRequest{
		Content: parentMessage.Message,
		Model:   r.embeddingModel,
	})
	if err != nil {
		ch <- &Response{
			IsValid:  false,
			Type:     ResponseError,
			Err: err,
		}
		return
	}
	for _, collectionName := range collectionNames {
		docs, err := r.milvusClinet.Load(ctx, collectionName, response.GetVector())
		if err != nil {
			ch <- &Response{
				IsValid:  false,
				Type:     ResponseError,
				Err: err,
			}
			continue
		}
		for _, doc := range docs {
			ch <- &Response{
				IsValid:  true,
				Type:     ResponseDocument,
				Document: &model.DocumentResponse{
					ID:         decimal.NewFromInt(doc.DocumentID),
					MessageID:  parentMessage.ID,
					Link:       doc.,
					DocumentID: decimal.NewFromInt(doc.DocumentID),
					Content:    "",
				},
			}
		}
	}


	for i := 0; i < 4; i++ {
		ch <- &Response{
			IsValid: true,
			Type:    ResponseDocument,
			Message: nil,
			Document: &model.DocumentResponse{
				ID:         decimal.NewFromInt(int64(i)),
				DocumentID: "11",
				Link:       fmt.Sprintf("link for document %d", i),
				Content:    fmt.Sprintf("content of document %d", i),
				MessageID:  parentMessage.ID,
			},
		}
	}
	wg.Done()
}


func NewEmbeddingResponder(embedding proto.EmbeddServiceClient,
	milvusClinet storage.MilvusClient,
	embeddingModel string ) ChatResponder {
	return &embedding{
		embedding:    embedding,
		milvusClinet: milvusClinet,
		embeddingModel: embeddingModel,
	}
}
