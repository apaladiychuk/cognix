package responder

import (
	"cognix.ch/api/v2/core/ai"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"context"
	"time"
)

type aiResponder struct {
	aiClient ai.OpenAIClient
	charRepo repository.ChatRepository
	ch       chan *Response
}

func (r *aiResponder) Send(ctx context.Context, parentMessage *model.ChatMessage) error {
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
	go func() {
		r.ch <- &Response{
			IsValid: err == nil,
			Message: &message,
		}
	}()
	return nil
}

func (r *aiResponder) Receive() (*Response, bool) {
	response := <-r.ch
	return response, false
}

func NewAIResponder(
	aiClient ai.OpenAIClient,
	charRepo repository.ChatRepository,
) ChatResponder {
	return &aiResponder{aiClient: aiClient,
		charRepo: charRepo,
		ch:       make(chan *Response),
	}
}
