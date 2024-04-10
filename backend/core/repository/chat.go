package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type ChatRepository interface {
	GetSessions(ctx context.Context, userID uuid.UUID) ([]*model.ChatSession, error)
	GetSessionByID(ctx context.Context, userID uuid.UUID, id int64) (*model.ChatSession, error)
	CreateSession(ctx context.Context, session *model.ChatSession) error
	SendMessage(ctx context.Context, message *model.ChatMessage) error
}

type chatRepository struct {
	db *pg.DB
}

func (r *chatRepository) SendMessage(ctx context.Context, message *model.ChatMessage) error {
	if _, err := r.db.WithContext(ctx).Model(message).Insert(); err != nil {
		return utils.Internal.Wrap(err, "can not save message")
	}
	return nil
}

func NewChatRepository(db *pg.DB) ChatRepository {
	return &chatRepository{db: db}
}
func (r *chatRepository) GetSessions(ctx context.Context, userID uuid.UUID) ([]*model.ChatSession, error) {
	sessions := make([]*model.ChatSession, 0)
	if err := r.db.WithContext(ctx).Model(&sessions).
		Where("user_id = ?", userID).
		Order(" created_date desc").Select(); err != nil {
		return nil, utils.NotFound.Wrapf(err, "can not find sessions")
	}
	return sessions, nil
}

func (r *chatRepository) GetSessionByID(ctx context.Context, userID uuid.UUID, id int64) (*model.ChatSession, error) {
	var session model.ChatSession
	if err := r.db.WithContext(ctx).Model(&session).
		Where("chat_session.user_id = ?", userID).
		Where("chat_session.id = ?", id).
		Relation("Persona").
		Relation("Persona.LLM").
		Relation("Messages").First(); err != nil {
		return nil, utils.NotFound.Wrapf(err, "can not find session")
	}
	return &session, nil
}

func (r *chatRepository) CreateSession(ctx context.Context, session *model.ChatSession) error {
	if _, err := r.db.WithContext(ctx).Model(session).Insert(); err != nil {
		return utils.Internal.Wrap(err, "can not create chat session")
	}
	return nil
}
