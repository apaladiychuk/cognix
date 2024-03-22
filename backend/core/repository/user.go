package repository

import (
	"cognix.ch/api/v2/core/model"
	"context"
	"github.com/go-pg/pg/v10"
)

type (
	UserRepository interface {
		GetByUserName(c context.Context, username string) (*model.User, error)
	}
	// UserRepository provides database operations with User model
	userRepository struct {
		db *pg.DB
	}
)

func NewUserRepository(db *pg.DB) UserRepository {
	return &userRepository{db: db}
}

func (u userRepository) GetByUserName(c context.Context, username string) (*model.User, error) {
	return &model.User{}, nil
}
