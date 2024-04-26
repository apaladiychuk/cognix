package parameters

import (
	"cognix.ch/api/v2/core/model"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/shopspring/decimal"
	"time"
)

const (
	MessageFeedbackUpvote   = "upvote"
	MessageFeedbackDownvote = "downvote"
)

type CreateChatSession struct {
	Description string          `json:"description"`
	PersonaID   decimal.Decimal `json:"persona_id"`
	OneShot     bool            `json:"one_shot"`
}

func (v CreateChatSession) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.PersonaID, validation.Required),
	)
}

type CreateChatMessageRequest struct {
	ChatSessionId    decimal.Decimal   `json:"chat_session_id,omitempty"`
	ParentMessageId  decimal.Decimal   `json:"parent_message_id,omitempty"`
	Message          string            `json:"message,omitempty"`
	PromptId         decimal.Decimal   `json:"prompt_id,omitempty"`
	SearchDocIds     []decimal.Decimal `json:"search_doc_ids,omitempty"`
	RetrievalOptions RetrievalDetails  `json:"retrieval_options,omitempty"`
	QueryOverride    string            `json:"query_override,omitempty"`
	NoAiAnswer       bool              `json:"no_ai_answer,omitempty"`
}

type RetrievalDetails struct {
	RunSearch               string      `json:"run_search,omitempty"`
	RealTime                bool        `json:"real_time,omitempty"`
	Filters                 BaseFilters `json:"filters,omitempty"`
	EnableAutoDetectFilters bool        `json:"enable_auto_detect_filters,omitempty"`
	Offset                  int         `json:"offset,omitempty"`
	Limit                   int         `json:"limit,omitempty"`
}
type BaseFilters struct {
	SourceType  []model.SourceType `json:"source_type,omitempty"`
	DocumentSet []string           `json:"document_set,omitempty"`
	TimeCutoff  time.Time          `json:"time_cutoff,omitempty"`
	Tags        []string           `json:"tags,omitempty"`
}

type MessageFeedbackParam struct {
	ID   decimal.Decimal `json:"id"`
	Vote string          `json:"vote"`
}

func (v MessageFeedbackParam) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.ID, validation.Required),
		validation.Field(&v.Vote, validation.Required, validation.In(MessageFeedbackDownvote, MessageFeedbackUpvote)))
}
