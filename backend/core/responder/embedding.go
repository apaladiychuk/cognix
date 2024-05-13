package responder

import (
	"cognix.ch/api/v2/core/model"
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"sync"
)

type embedding struct {
}

func (r *embedding) Send(ctx context.Context, ch chan *Response, wg *sync.WaitGroup, parentMessage *model.ChatMessage) {

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
				MessageID:  parentMessage.ID,
			},
		}
	}
	wg.Done()
}

func NewEmbeddingResponder() ChatResponder {
	return &embedding{}
}
