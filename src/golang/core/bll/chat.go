package bll

import (
	"cognix.ch/api/v2/core/ai"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/responder"
	"cognix.ch/api/v2/core/storage"
	"cognix.ch/api/v2/core/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"time"
)

type ChatBL interface {
	GetSessions(ctx context.Context, user *model.User) ([]*model.ChatSession, error)
	GetSessionByID(ctx context.Context, user *model.User, id int64) (*model.ChatSession, error)
	CreateSession(ctx context.Context, user *model.User, param *parameters.CreateChatSession) (*model.ChatSession, error)
	SendMessage(ctx *gin.Context, user *model.User, param *parameters.CreateChatMessageRequest) (*responder.Manager, error)
	FeedbackMessage(ctx *gin.Context, user *model.User, id int64, vote bool) (*model.ChatMessageFeedback, error)
}
type chatBL struct {
	chatRepo     repository.ChatRepository
	docRepo      repository.DocumentRepository
	personaRepo  repository.PersonaRepository
	aiBuilder    *ai.Builder
	embedding    proto.EmbedServiceClient
	milvusClinet storage.MilvusClient
}

func (b *chatBL) FeedbackMessage(ctx *gin.Context, user *model.User, id int64, vote bool) (*model.ChatMessageFeedback, error) {
	message, err := b.chatRepo.GetMessageByIDAndUserID(ctx, id, user.ID)
	if err != nil {
		return nil, err
	}
	feedback := message.Feedback
	if feedback == nil {
		feedback = &model.ChatMessageFeedback{
			ChatMessageID: message.ID,
			UserID:        user.ID,
		}
	}
	feedback.UpVotes = vote
	if err = b.chatRepo.MessageFeedback(ctx, feedback); err != nil {
		return nil, err
	}
	return feedback, nil
}

func (b *chatBL) SendMessage(ctx *gin.Context, user *model.User, param *parameters.CreateChatMessageRequest) (*responder.Manager, error) {
	chatSession, err := b.chatRepo.GetSessionByID(ctx.Request.Context(), user.ID, param.ChatSessionID.IntPart())
	if err != nil {
		return nil, err
	}
	message := model.ChatMessage{
		ChatSessionID: chatSession.ID,
		Message:       param.Message,
		MessageType:   model.MessageTypeUser,
		TimeSent:      time.Now().UTC(),
	}
	noLLM := chatSession.Persona == nil
	if err = b.chatRepo.SendMessage(ctx.Request.Context(), &message); err != nil {
		return nil, err
	}
	aiClient := b.aiBuilder.New(chatSession.Persona.LLM)
	resp := responder.NewManager(
		responder.NewAIResponder(aiClient, b.chatRepo,
			b.embedding, b.milvusClinet, b.docRepo, ""),
	)

	go resp.Send(ctx, user, noLLM, &message)
	return resp, nil
}

func (b *chatBL) GetSessions(ctx context.Context, user *model.User) ([]*model.ChatSession, error) {
	return b.chatRepo.GetSessions(ctx, user.ID)
}

func (b *chatBL) GetSessionByID(ctx context.Context, user *model.User, id int64) (*model.ChatSession, error) {
	result, err := b.chatRepo.GetSessionByID(ctx, user.ID, id)
	if err != nil {
		return nil, err
	}
	docs := make([]model.DocumentResponse, 0)
	for i := 0; i < 4; i++ {
		docs = append(docs, model.DocumentResponse{
			ID:          decimal.NewFromInt(int64(i)),
			DocumentID:  "11",
			Link:        fmt.Sprintf("link for document %d", i),
			Content:     fmt.Sprintf("content of document %d", i),
			UpdatedDate: time.Now().UTC().Add(-48 * time.Hour),
		})
	}
	for _, msg := range result.Messages {
		if msg.MessageType == model.MessageTypeAssistant {
			for _, d := range docs {
				md := d
				md.MessageID = msg.ID
				msg.Citations = append(msg.Citations, &md)
			}
		}

	}
	return result, nil
}

func (b *chatBL) CreateSession(ctx context.Context, user *model.User, param *parameters.CreateChatSession) (*model.ChatSession, error) {
	exists, err := b.personaRepo.IsExists(ctx, param.PersonaID.IntPart(), user.TenantID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, utils.ErrorBadRequest.New("persona is not exists")
	}
	session := model.ChatSession{
		UserID:      user.ID,
		Description: param.Description,
		CreatedDate: time.Now().UTC(),
		PersonaID:   param.PersonaID,
		OneShot:     param.OneShot,
	}
	if err = b.chatRepo.CreateSession(ctx, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func NewChatBL(chatRepo repository.ChatRepository,
	personaRepo repository.PersonaRepository,
	docRepo repository.DocumentRepository,
	aiBuilder *ai.Builder,
	embedding proto.EmbedServiceClient,
	milvusClinet storage.MilvusClient,
) ChatBL {
	return &chatBL{chatRepo: chatRepo,
		personaRepo:  personaRepo,
		docRepo:      docRepo,
		aiBuilder:    aiBuilder,
		embedding:    embedding,
		milvusClinet: milvusClinet,
	}
}
