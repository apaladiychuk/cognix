package responder

import (
	"cognix.ch/api/v2/core/ai"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"context"
	"sync"
	"time"
)

type aiResponder struct {
	aiClient ai.OpenAIClient
	charRepo repository.ChatRepository
}

func (r *aiResponder) Send(ctx context.Context, ch chan *Response, wg *sync.WaitGroup, parentMessage *model.ChatMessage) {

	response, err := r.aiClient.Request(ctx, parentMessage.Message)
	message := model.ChatMessage{
		ChatSessionID: parentMessage.ChatSessionID,
		ParentMessage: parentMessage.ID,
		MessageType:   model.MessageTypeAssistant,
		TimeSent:      time.Now().UTC(),
	}
	if err != nil {
		message.Error = err.Error()
	} else {
		message.Message = response.Message
	}

	if errr := r.charRepo.SendMessage(ctx, &message); errr != nil {
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
	wg.Done()
}

func NewAIResponder(
	aiClient ai.OpenAIClient,
	charRepo repository.ChatRepository,
) ChatResponder {
	return &aiResponder{aiClient: aiClient,
		charRepo: charRepo,
	}
}
