package model

import (
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"time"
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
		OneShot     bool           `json:"one_shot,omitempty"`
		Messages    []*ChatMessage `json:"messages,omitempty" pg:"rel:has-many"`
	}

	ChatMessage struct {
		tableName          struct{}  `pg:"chat_messages"`
		ID                 int64     `json:"id,omitempty"`
		ChatSessionID      int64     `json:"chat_session_id,omitempty"`
		Message            string    `json:"message,omitempty"`
		MessageType        string    `json:"message_type,omitempty"`
		TimeSent           time.Time `json:"time_sent,omitempty"`
		TokenCount         int       `json:"token_count,omitempty"`
		ParentMessage      int       `json:"parent_message,omitempty"`
		LatestChildMessage int       `json:"latest_child_message,omitempty"`
		RephrasedQuery     string    `json:"rephrased_query,omitempty"`
		Citations          JSON      `json:"citations,omitempty"`
		Error              string    `json:"error,omitempty"`
	}
)
