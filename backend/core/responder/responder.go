package responder

import (
	"cognix.ch/api/v2/core/model"
	"context"
	"fmt"
	"time"
)

const (
	ResponseMessage = "message"
	ResponseError   = "error"
)

type Response struct {
	IsValid bool
	Message *model.ChatMessage
}

type ChatResponder interface {
	Send(cx context.Context, parentMessage *model.ChatMessage) error
	Receive() (*Response, bool)
}

type nopResponder struct {
	ch chan string
}

func (r *nopResponder) Send(cx context.Context, parentMessage *model.ChatMessage) error {
	go func() {
		i := 0
		for i < 3 {
			time.Sleep(20 * time.Second)
			r.ch <- fmt.Sprintf("response %d\n", i+1)
			i++
		}
		close(r.ch)
	}()

	return nil
}

func (r *nopResponder) Receive() (*Response, bool) {
	message, ok := <-r.ch
	if !ok {
		return nil, false
	}
	return &Response{Message: &model.ChatMessage{Message: message}}, true
}

func NewNopResponder() ChatResponder {
	ch := make(chan string)
	return &nopResponder{
		ch: ch,
	}
}
