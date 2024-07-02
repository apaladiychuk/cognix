package responder

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/storage"
	"context"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// embedding represents a struct that contains instances of various interfaces and types used for document embedding.
// Fields:
// - embedding: an instance of EmbedServiceClient interface for getting embeddings.
// - milvusClinet: an instance of VectorDBClient interface for interacting with the Milvus storage.
// - docRepo: an instance of DocumentRepository interface for document persistence.
// - embeddingModel: a string representing the embedding model being used.
type embedding struct {
	embedding      proto.EmbedServiceClient
	milvusClinet   storage.VectorDBClient
	docRepo        repository.DocumentRepository
	embeddingModel string
}

// FindDocuments is a method that searches for documents in multiple collections based on a chat message.
//
// Parameters:
//   - ctx: The context.Context object for the request.
//   - ch: The channel to send the response messages to.
//   - message: The chat message containing the content to search.
//   - collectionNames: The names of the collections to search for documents.
//
// Returns:
//   - []*model.DocumentResponse: The list of document responses found.
//   - error: An error if any occurred during the search.
//
// Behavior:
//   - The method calls the GetEmbeding method of the embdding service to get the embedding of the message content.
//   - For each collection specified, the method loads documents using the Load method of the VectorDBClient.
//   - For each loaded document, it creates a DocumentResponse object and populates its fields based on the database data.
//   - The DocumentResponse objects are stored in a map to avoid duplicates, and the valid ones are also added to the result list.
//   - For each valid DocumentResponse, a Response object of type ResponseDocument is sent to the channel.
//
// Errors:
//   - If an error occurs while calling the GetEmbeding method or loading a document, an error is returned and a Response
//     object of type ResponseError is sent to the channel.
//
// Notes:
//   - The MessageID field of the DocumentResponse objects is set to the ID of the input chat message.
//   - The Link field of the DocumentResponse objects is set to the OriginalURL field of the corresponding document from
//     the database. If OriginalURL is empty, the Link field is set to the URL field.
//   - The DocumentID field of the DocumentResponse objects is set to the SourceID field of the corresponding document from
//     the database.
//   - The UpdatedDate field of the DocumentResponse objects is set to LastUpdate field of the corresponding document from
//     the database, if it is not zero. Otherwise, it is set to the CreationDate field.
//   - The DocumentResponse objects are sent to the channel as Response objects with type ResponseDocument.
func (r *embedding) FindDocuments(ctx context.Context,
	ch chan *Response,
	message *model.ChatMessage,
	collectionNames ...string) ([]*model.DocumentResponse, error) {
	response, err := r.embedding.GetEmbeding(ctx, &proto.EmbedRequest{
		Content: message.ParentMessage.Message,
		Model:   r.embeddingModel,
	})
	if err != nil {
		zap.S().Errorf("embeding service %s ", err.Error())
		ch <- &Response{
			IsValid: false,
			Type:    ResponseError,
			Err:     err,
		}
		return nil, err
	}
	var result []*model.DocumentResponse
	mapResult := make(map[string]*model.DocumentResponse)
	for _, collectionName := range collectionNames {
		docs, err := r.milvusClinet.Load(ctx, collectionName, response.GetVector())
		if err != nil {
			zap.S().Errorf("error loading document from vector database :%s", err.Error())
			continue
		}
		for _, doc := range docs {
			resDocument := &model.DocumentResponse{
				ID:        decimal.NewFromInt(doc.DocumentID),
				MessageID: message.ID,
				Content:   doc.Content,
			}
			if dbDoc, err := r.docRepo.FindByID(ctx, doc.DocumentID); err == nil {
				resDocument.Link = dbDoc.OriginalURL
				if resDocument.Link == "" {
					resDocument.Link = dbDoc.URL
				}
				resDocument.DocumentID = dbDoc.SourceID
				if !dbDoc.LastUpdate.IsZero() {
					resDocument.UpdatedDate = dbDoc.LastUpdate.Time
				} else {
					resDocument.UpdatedDate = dbDoc.CreationDate
				}
			}
			if _, ok := mapResult[resDocument.DocumentID]; ok {
				continue
			}
			mapResult[resDocument.DocumentID] = resDocument
			result = append(result, resDocument)
			ch <- &Response{
				IsValid:  true,
				Type:     ResponseDocument,
				Document: resDocument,
			}
		}
	}
	return result, nil
}

// NewEmbeddingResponder creates a new instance of embedding struct
//
// Parameters:
//   - embeddProto : EmbedServiceClient for embedding service API
//   - milvusClient: VectorDBClient for interacting with the Milvus storage
//   - docRepo     : DocumentRepository for interacting with the document data
//   - embeddingModel   : The embedding model string
//
// Returns:
//   - *embedding  : A pointer to the embedding struct
func NewEmbeddingResponder(embeddProto proto.EmbedServiceClient,
	milvusClinet storage.VectorDBClient,
	docRepo repository.DocumentRepository,
	embeddingModel string) *embedding {
	return &embedding{
		embedding:      embeddProto,
		milvusClinet:   milvusClinet,
		embeddingModel: embeddingModel,
		docRepo:        docRepo,
	}
}
