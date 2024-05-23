package responder

import (
	"cognix.ch/api/v2/core/ai"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/storage"
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"sync"
	"time"
)

type aiResponder struct {
	aiClient  ai.OpenAIClient
	charRepo  repository.ChatRepository
	embedding *embedding
}

func (r *aiResponder) Send(ctx context.Context, ch chan *Response, wg *sync.WaitGroup, user *model.User, noLLM bool, parentMessage *model.ChatMessage) {
	defer wg.Done()
	message := model.ChatMessage{
		ChatSessionID:   parentMessage.ChatSessionID,
		ParentMessageID: parentMessage.ID,
		MessageType:     model.MessageTypeAssistant,
		TimeSent:        time.Now().UTC(),
		ParentMessage:   parentMessage,
		Message:         "You are using Cognix without an LLM. I can give you the documents retrieved in my knowledge. ",
	}
	if err := r.charRepo.SendMessage(ctx, &message); err != nil {
		ch <- &Response{
			IsValid: err == nil,
			Type:    ResponseMessage,
			Message: &message,
		}
		return
	}

	//docs, err := r.embedding.FindDocuments(ctx, ch, &message, model.CollectionName(true, user.ID, user.TenantID),
	//	model.CollectionName(false, user.ID, user.TenantID))
	//if err != nil {
	//
	//}
	if noLLM {
		return
	}
	message.Message = ""
	//_ = docs
	// docs.Content
	// user chat
	// system_prompt
	// task_prompt
	// default_prompt
	// llm message format : system prompt \n user chat \n task_prompt \n document content1 \n ...\n document content n ( top 5)
	//

	response, err := r.aiClient.Request(ctx, parentMessage.Message)

	if err != nil {
		message.Error = err.Error()
	} else {
		message.Message = response.Message
	}

	if errr := r.charRepo.Update(ctx, &message); errr != nil {
		err = errr
		message.Error = err.Error()
	}
	payload := &Response{
		IsValid: err == nil,
		Type:    ResponseMessage,
		Message: &message,
	}
	if err != nil {
		payload.Type = ResponseError
	}
	ch <- payload
	time.Sleep(10 * time.Millisecond)
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
				MessageID:   message.ID,
			},
		}
	}
}

func NewAIResponder(
	aiClient ai.OpenAIClient,
	charRepo repository.ChatRepository,
	embeddProto proto.EmbedServiceClient,
	milvusClinet storage.MilvusClient,
	docRepo repository.DocumentRepository,
	embeddingModel string,
) ChatResponder {
	return &aiResponder{aiClient: aiClient,
		charRepo:  charRepo,
		embedding: NewEmbeddingResponder(embeddProto, milvusClinet, docRepo, embeddingModel),
	}
}
