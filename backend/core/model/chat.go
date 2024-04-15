package model

import (
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"time"
)

const (
	MessageTypeUser      = "user"
	MessageTypeAssistant = "assistant"
	MessageTypeSystem    = "system"
)

type (
	ChatSession struct {
		tableName   struct{}       `pg:"chat_sessions"`
		ID          int64          `json:"id,omitempty"`
		UserID      uuid.UUID      `json:"user_id,omitempty"`
		Description string         `json:"description,omitempty"`
		CreatedDate time.Time      `json:"created_date,omitempty"`
		DeletedDate pg.NullTime    `json:"deleted_date,omitempty"`
		PersonaID   int64          `json:"persona_id,omitempty"`
		OneShot     bool           `json:"one_shot,omitempty" pg:",use_zero"`
		Messages    []*ChatMessage `json:"messages,omitempty" pg:"rel:has-many"`
		Persona     *Persona       `json:"persona,omitempty" pg:"rel:has-one"`
	}

	ChatMessage struct {
		tableName          struct{}             `pg:"chat_messages"`
		ID                 int64                `json:"id,omitempty"`
		ChatSessionID      int64                `json:"chat_session_id,omitempty"`
		Message            string               `json:"message,omitempty"`
		MessageType        string               `json:"message_type,omitempty"`
		TimeSent           time.Time            `json:"time_sent,omitempty"`
		TokenCount         int                  `json:"token_count,omitempty" pg:",use_zero"`
		ParentMessage      int64                `json:"parent_message,omitempty" pg:",use_zero"`
		LatestChildMessage int                  `json:"latest_child_message,omitempty" pg:",use_zero"`
		RephrasedQuery     string               `json:"rephrased_query,omitempty" pg:",use_zero"`
		Citations          JSON                 `json:"citations,omitempty"`
		Error              string               `json:"error,omitempty" pg:",use_zero"`
		Feedback           *ChatMessageFeedback `json:"feedback,omitempty" pg:"rel:has-one,fk:id,join_fk:chat_message_id"`
	}

	ChatMessageFeedback struct {
		tableName     struct{}  `pg:"chat_message_feedbacks"`
		ID            int64     `json:"id,omitempty"`
		ChatMessageID int64     `json:"chat_message_id,omitempty"`
		UserID        uuid.UUID `json:"user_id,omitempty"`
		UpVotes       bool      `json:"up_votes,omitempty" pg:",use_zero"`
		Feedback      string    `json:"feedback,omitempty" pg:",use_zero"`
	}
)
