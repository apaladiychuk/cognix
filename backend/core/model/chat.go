package model

import (
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"time"
)

type (
	ChatSession struct {
		tableName   struct{}       `pg:"chat_sessions"`
		ID          int            `json:"id,omitempty"`
		UserID      uuid.UUID      `json:"user_id,omitempty"`
		Description string         `json:"description,omitempty"`
		CreatedDate time.Time      `json:"created_date,omitempty"`
		DeletedDate null.Time      `json:"deleted_date,omitempty"`
		PersonaID   int            `json:"persona_id,omitempty"`
		OneShot     bool           `json:"one_shot,omitempty"`
		Messages    []*ChatMessage `json:"messages,omitempty" pg:"rel:has-many"`
	}

	ChatMessage struct {
		tableName          struct{}  `pg:"chat_messages"`
		ID                 int       `json:"id,omitempty"`
		ChatSessionID      int       `json:"chat_session_id,omitempty"`
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
