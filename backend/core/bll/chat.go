package bll

import (
	"cognix.ch/api/v2/core/ai"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/responder"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/gin-gonic/gin"
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
	chatRepo    repository.ChatRepository
	personaRepo repository.PersonaRepository
	aiBuilder   *ai.Builder
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
	chatSession, err := b.chatRepo.GetSessionByID(ctx.Request.Context(), user.ID, param.ChatSessionId.IntPart())
	if err != nil {
		return nil, err
	}
	message := model.ChatMessage{
		ChatSessionID: chatSession.ID,
		Message:       param.Message,
		MessageType:   model.MessageTypeUser,
		TimeSent:      time.Now().UTC(),
	}
	if err = b.chatRepo.SendMessage(ctx.Request.Context(), &message); err != nil {
		return nil, err
	}
	aiClient := b.aiBuilder.New(chatSession.Persona.LLM)
	resp := responder.NewManager(
		responder.NewAIResponder(aiClient, b.chatRepo),
		responder.NewEmbeddingResponder())

	go resp.Send(ctx, &message)
	return resp, nil
}

func (b *chatBL) GetSessions(ctx context.Context, user *model.User) ([]*model.ChatSession, error) {
	return b.chatRepo.GetSessions(ctx, user.ID)
}

func (b *chatBL) GetSessionByID(ctx context.Context, user *model.User, id int64) (*model.ChatSession, error) {
	return b.chatRepo.GetSessionByID(ctx, user.ID, id)
}

func (b *chatBL) CreateSession(ctx context.Context, user *model.User, param *parameters.CreateChatSession) (*model.ChatSession, error) {
	exists, err := b.personaRepo.IsExists(ctx, param.PersonaID.IntPart(), user.TenantID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, utils.InvalidInput.New("persona is not exists")
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
	aiBuilder *ai.Builder,
) ChatBL {
	return &chatBL{chatRepo: chatRepo,
		personaRepo: personaRepo,
		aiBuilder:   aiBuilder,
	}
}
