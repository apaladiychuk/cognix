package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type (
	PersonaRepository interface {
		GetAll(ctx context.Context, tenantID uuid.UUID) ([]*model.Persona, error)
		GetByID(ctx context.Context, id int64, tenantID uuid.UUID) (*model.Persona, error)
		Create(ctx context.Context, persona *model.Persona) error
		Update(ctx context.Context, persona *model.Persona) error
		IsExists(ctx context.Context, id int64, tenantID uuid.UUID) (bool, error)
	}
	personaRepository struct {
		db *pg.DB
	}
)

func (r *personaRepository) IsExists(ctx context.Context, id int64, tenantID uuid.UUID) (bool, error) {
	exist, err := r.db.WithContext(ctx).Model(&model.Persona{}).
		Where("id = ?", id).Where("tenant_id = ?", tenantID).
		Exists()
	if err != nil {
		return false, utils.NotFound.Wrap(err, "can not find persona")
	}
	return exist, nil
}

func NewPersonaRepository(db *pg.DB) PersonaRepository {
	return &personaRepository{db: db}
}

func (r *personaRepository) GetAll(ctx context.Context, tenantID uuid.UUID) ([]*model.Persona, error) {
	personas := make([]*model.Persona, 0)
	if err := r.db.WithContext(ctx).Model(&personas).
		Where("tenant_id = ?", tenantID).Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "personas not found")
	}
	return personas, nil
}

func (r *personaRepository) GetByID(ctx context.Context, id int64, tenantID uuid.UUID) (*model.Persona, error) {
	var persona model.Persona
	if err := r.db.WithContext(ctx).Model(&persona).
		Relation("LLM").
		Relation("Prompt").
		Where("persona.id = ?", id).
		Where("persona.tenant_id = ?", tenantID).First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "persona not found")
	}
	return &persona, nil
}

func (r *personaRepository) Create(ctx context.Context, persona *model.Persona) error {
	return r.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		if _, err := tx.Model(persona.LLM).Insert(); err != nil {
			return utils.Internal.Wrap(err, "can not insert LLM")
		}
		persona.LlmID = persona.LLM.ID
		if _, err := tx.Model(persona).Insert(); err != nil {
			return utils.Internal.Wrap(err, "can not insert persona")
		}
		persona.Prompt.PersonaID = persona.ID
		if _, err := tx.Model(persona.Prompt).Insert(); err != nil {
			return utils.Internal.Wrap(err, "can not insert prompt")
		}
		return nil
	})
}

func (r *personaRepository) Update(ctx context.Context, persona *model.Persona) error {
	return r.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		if _, err := tx.Model(persona.LLM).Where("id = ?", persona.LLM.ID).Update(); err != nil {
			return utils.Internal.Wrap(err, "can not update LLM")
		}
		persona.LlmID = persona.LLM.ID
		if _, err := tx.Model(persona).Where("id = ?", persona.ID).Update(); err != nil {
			return utils.Internal.Wrap(err, "can not update persona")
		}
		persona.Prompt.PersonaID = persona.ID
		if _, err := tx.Model(persona.Prompt).Where("id = ?", persona.Prompt.ID).Update(); err != nil {
			return utils.Internal.Wrap(err, "can not update prompt")
		}
		return nil
	})
}
