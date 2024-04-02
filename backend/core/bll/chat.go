package bll

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"time"
)

type ChatBL interface {
	GetSessions(ctx context.Context, user *model.User) ([]*model.ChatSession, error)
	GetSessionByID(ctx context.Context, user *model.User, id int64) (*model.ChatSession, error)
	CreateSession(ctx context.Context, user *model.User, param *parameters.CreateChatSession) (*model.ChatSession, error)
}
type chatBL struct {
	chatRepo    repository.ChatRepository
	personaRepo repository.PersonaRepository
}

func (b *chatBL) GetSessions(ctx context.Context, user *model.User) ([]*model.ChatSession, error) {
	return b.chatRepo.GetSessions(ctx, user.ID)
}

func (b *chatBL) GetSessionByID(ctx context.Context, user *model.User, id int64) (*model.ChatSession, error) {
	return b.chatRepo.GetSessionByID(ctx, user.ID, id)
}

func (b *chatBL) CreateSession(ctx context.Context, user *model.User, param *parameters.CreateChatSession) (*model.ChatSession, error) {
	exists, err := b.personaRepo.IsExists(ctx, param.PersonaID, user.TenantID)
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
	personaRepo repository.PersonaRepository) ChatBL {
	return &chatBL{chatRepo: chatRepo,
		personaRepo: personaRepo}
}
