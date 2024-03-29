package repository

import (
	"cognix.ch/api/v2/core/model"
	"context"
	"github.com/go-pg/pg/v10"
)

type (
	PersonaRepository interface {
		GetAll(ctx context.Context, tenantID, userID string) ([]*model.Persona, error)
		GetByID(c context.Context, id int64, tenantID, userID string) (*model.Persona, error)
		Create(c context.Context, connector *model.Persona) error
		Update(c context.Context, connector *model.Persona) error
	}
	personaRepository struct {
		db *pg.DB
	}
)

func NewPersonaRepository(db *pg.DB) PersonaRepository {
	return &personaRepository{db: db}
}

func (r *personaRepository) GetAll(ctx context.Context, tenantID, userID string) ([]*model.Persona, error) {
	//TODO implement me
	panic("implement me")
}

func (r *personaRepository) GetByID(c context.Context, id int64, tenantID, userID string) (*model.Persona, error) {
	//TODO implement me
	panic("implement me")
}

func (r *personaRepository) Create(c context.Context, connector *model.Persona) error {
	//TODO implement me
	panic("implement me")
}

func (r *personaRepository) Update(c context.Context, connector *model.Persona) error {
	//TODO implement me
	panic("implement me")
}
