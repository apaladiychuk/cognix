package ai

import "cognix.ch/api/v2/core/model"

type Builder struct {
	clients map[int64]OpenAIClient
}

func NewBuilder() *Builder {
	return &Builder{clients: make(map[int64]OpenAIClient)}
}

func (b *Builder) New(llm *model.LLM) OpenAIClient {
	if client, ok := b.clients[llm.ID]; ok {
		return client
	}
	client := NewOpenAIClient(llm)
	b.clients[llm.ID] = client
	return client
}
