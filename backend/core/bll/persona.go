package bll

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"context"
)

type (
	PersonaBL interface {
		GetAll(ctx context.Context, user *model.User) ([]*model.Persona, error)
		GetByID(ctx context.Context, user *model.User, id int64) (*model.Persona, error)
		//Create(ctx context.Context, user *model.User, param *parameters.CreateConnectorParam) (*model.Persona, error)
		//Update(ctx context.Context, id int64, user *model.User, param *parameters.UpdateConnectorParam) (*model.Persona, error)
	}
	personaBL struct {
		personaRepo repository.PersonaRepository
	}
)

func NewPersonaBL(personaRepo repository.PersonaRepository) PersonaBL {
	return &personaBL{
		personaRepo: personaRepo,
	}
}
func (b *personaBL) GetAll(ctx context.Context, user *model.User) ([]*model.Persona, error) {
	return b.personaRepo.GetAll(ctx, user.TenantID)
}

func (b *personaBL) GetByID(ctx context.Context, user *model.User, id int64) (*model.Persona, error) {
	return b.personaRepo.GetByID(ctx, id, user.TenantID)
}
