package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
)

type (
	LLMRepository interface {
		GetAll(ctx context.Context) ([]*model.LLM, error)
	}
	llmRepository struct {
		db *pg.DB
	}
)

func NewLLMRepository(db *pg.DB) LLMRepository {
	return &llmRepository{db: db}
}

func (l *llmRepository) GetAll(ctx context.Context) ([]*model.LLM, error) {
	llms := make([]*model.LLM, 0)
	if err := l.db.WithContext(ctx).Model(&llms).Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find llm")
	}
	return llms, nil
}
