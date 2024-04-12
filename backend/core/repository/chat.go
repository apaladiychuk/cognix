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
	GetMessageByIDAndUserID(ctx context.Context, id int64, userID uuid.UUID) (*model.ChatMessage, error)
	MessageFeedback(ctx context.Context, feedback *model.ChatMessageFeedback) error
}

type chatRepository struct {
	db *pg.DB
}

func (r *chatRepository) GetMessageByIDAndUserID(ctx context.Context, id int64, userID uuid.UUID) (*model.ChatMessage, error) {
	var message model.ChatMessage
	if err := r.db.Model(&message).
		Relation("Feedback").
		Join("inner join chat_sessions on chat_sessions.id = chat_messages.session_id and chat_session.user_id = ?", userID).
		Where("chat_messages.id = ?", id).First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "cannot find message by id")
	}
	return &message, nil
}

func (r *chatRepository) MessageFeedback(ctx context.Context, feedback *model.ChatMessageFeedback) error {
	stm := r.db.WithContext(ctx).Model(feedback)
	if feedback.ID == 0 {
		if _, err := stm.Insert(); err != nil {
			return utils.Internal.Wrap(err, "can not add feedback")
		}
		return nil
	}
	if _, err := stm.Where("id = ?", feedback.ID).Update(); err != nil {
		return utils.Internal.Wrap(err, "can not update feedback")
	}
	return nil
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
