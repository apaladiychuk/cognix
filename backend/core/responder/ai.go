package responder

import (
	"cognix.ch/api/v2/core/ai"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"context"
	"fmt"
	"github.com/shopspring/decimal"
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
	time.Sleep(10 * time.Millisecond)
	for i := 0; i < 4; i++ {
		ch <- &Response{
			IsValid: true,
			Type:    ResponseDocument,
			Message: nil,
			Document: &model.DocumentResponse{
				ID:         decimal.NewFromInt(int64(i)),
				DocumentID: "11",
				Link:       fmt.Sprintf("link for document %d", i),
				Content:    fmt.Sprintf("content of document %d", i),
				MessageID:  message.ID,
			},
		}
	}

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
