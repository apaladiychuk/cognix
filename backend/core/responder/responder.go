package responder

import (
	"cognix.ch/api/v2/core/model"
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	ResponseMessage  = "message"
	ResponseError    = "error"
	ResponseDocument = "document"
	ResponseEnd      = "end"
)

type Response struct {
	IsValid  bool
	Type     string
	Message  *model.ChatMessage
	Document *model.DocumentResponse
	Err      error
}

type ChatResponder interface {
	Send(cx context.Context, ch chan *Response, wg *sync.WaitGroup, parentMessage *model.ChatMessage)
}
type nopResponder struct {
	ch chan string
}

func (r *nopResponder) Send(cx context.Context, ch chan *Response, wg *sync.WaitGroup, parentMessage *model.ChatMessage) {
	go func() {
		defer wg.Done()
		i := 0
		for i < 3 {
			time.Sleep(20 * time.Second)
			r.ch <- fmt.Sprintf("response %d\n", i+1)
			i++
		}
		close(r.ch)
	}()

	return
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
