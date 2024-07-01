package responder

import (
	"cognix.ch/api/v2/core/ai"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/storage"
	"context"
	"github.com/google/uuid"
	"strings"
	"sync"
	"time"
)

// aiResponder represents a type that handles AI responses in a chat application.
//
// Fields:
// - aiClient: an implementation of the OpenAIClient interface for making requests to the OpenAI chat API.
// - charRepo: an implementation of the ChatRepository interface for interacting with the chat repository.
// - embedding: an instance of the embedding type for handling document embeddings.
type aiResponder struct {
	aiClient  ai.OpenAIClient
	charRepo  repository.ChatRepository
	embedding *embedding
}

// Send sends a chat message to the AI Responder.
// It creates a new ChatMessage with the provided parameters, sets its attributes,
// and sends it to the AI Responder using the charRepo.SendMessage method.
//
// Parameters:
// - ctx: the context.Context object for the request.
// - ch: the channel to send the response to.
// - wg: the sync.WaitGroup to wait for the response.
// - user: the model.User object representing the user.
// - noLLM: a boolean value indicating whether to skip LLM processing.
// - parentMessage: the parent ChatMessage of the new message.
// - persona: the model.Persona object representing the persona.
//
// Returns: none.
func (r *aiResponder) Send(ctx context.Context,
	ch chan *Response,
	wg *sync.WaitGroup,
	user *model.User, noLLM bool, parentMessage *model.ChatMessage, persona *model.Persona) {
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

	docs, err := r.embedding.FindDocuments(ctx, ch, &message, model.CollectionName(user.ID, uuid.NullUUID{Valid: true, UUID: user.TenantID}),
		model.CollectionName(user.ID, uuid.NullUUID{Valid: false}))
	if err != nil {

	}
	if noLLM {
		return
	}
	messageParts := []string{
		persona.Prompt.SystemPrompt,
		parentMessage.Message,
		persona.Prompt.TaskPrompt,
	}

	for _, doc := range docs {
		messageParts = append(messageParts, doc.Content)
		if doc.ID.IntPart() != 0 {
			message.DocumentPairs = append(message.DocumentPairs, &model.ChatMessageDocumentPair{
				ChatMessageID: message.ID,
				DocumentID:    doc.ID,
			})
		}
	}
	message.Citations = docs
	message.Message = ""
	//_ = docs
	// docs.Content
	// user chat
	// system_prompt
	// task_prompt
	// default_prompt
	// llm message format : system prompt \n user chat \n task_prompt \n document content1 \n ...\n document content n ( top 5)
	//
	response, err := r.aiClient.Request(ctx, strings.Join(messageParts, "\n"))

	if err != nil {
		message.Error = err.Error()
	} else {
		message.Message = response.Message
	}

	if errr := r.charRepo.UpdateMessage(ctx, &message); errr != nil {
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
}

// NewAIResponder creates a new instance of AIResponder.
//
// Parameters:
//   - aiClient: The AI client for making requests to the OpenAI chat API.
//   - charRepo: The chat repository for interacting with the chat data.
//   - embeddProto: The EmbedService client for embedding service API.
//   - milvusClinet: The MilvusClient for interacting with the Milvus storage.
//   - docRepo: The document repository for interacting with the document data.
//   - embeddingModel: The embedding model string.
//
// Returns:
//   - ChatResponder: The ChatResponder object for sending chat responses.
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
