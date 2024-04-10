package bll

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/repository"
	"context"
	"fmt"
	"time"
)

type (
	PersonaBL interface {
		GetAll(ctx context.Context, user *model.User) ([]*model.Persona, error)
		GetByID(ctx context.Context, user *model.User, id int64) (*model.Persona, error)
		Create(ctx context.Context, user *model.User, param *parameters.PersonaParam) (*model.Persona, error)
		Update(ctx context.Context, id int64, user *model.User, param *parameters.PersonaParam) (*model.Persona, error)
	}
	personaBL struct {
		personaRepo repository.PersonaRepository
	}
)

func (b *personaBL) Create(ctx context.Context, user *model.User, param *parameters.PersonaParam) (*model.Persona, error) {

	persona := model.Persona{
		Name:            param.Name,
		LlmID:           0,
		DefaultPersona:  true,
		Description:     param.Description,
		TenantID:        user.TenantID,
		IsVisible:       true,
		DisplayPriority: 0,
		StarterMessages: nil,
		LLM: &model.LLM{
			Name:     fmt.Sprintf("%s %s", user.FirstName, param.ModelID),
			ModelID:  param.ModelID,
			Url:      param.URL,
			ApiKey:   param.APIKey,
			Endpoint: param.Endpoint,
		},
		Prompt: &model.Prompt{
			UserID:           user.ID,
			Name:             param.Name,
			Description:      param.Description,
			SystemPrompt:     param.SystemPrompt,
			TaskPrompt:       param.TaskPrompt,
			IncludeCitations: false,
			DatetimeAware:    false,
			DefaultPrompt:    false,
			CreatedDate:      time.Now().UTC(),
		},
	}
	if err := b.personaRepo.Create(ctx, &persona); err != nil {
		return nil, err
	}
	return &persona, nil
}

func (b *personaBL) Update(ctx context.Context, id int64, user *model.User, param *parameters.PersonaParam) (*model.Persona, error) {
	persona, err := b.personaRepo.GetByID(ctx, id, user.TenantID)
	if err != nil {
		return nil, err
	}
	persona.Name = param.Name
	persona.Description = param.Description
	persona.LLM.Endpoint = param.Endpoint
	persona.LLM.ModelID = param.ModelID
	persona.LLM.ApiKey = param.APIKey
	persona.Prompt.Name = param.Name
	persona.Prompt.Description = param.Description
	persona.Prompt.SystemPrompt = param.SystemPrompt
	persona.Prompt.TaskPrompt = param.TaskPrompt

	if err = b.personaRepo.Update(ctx, persona); err != nil {
		return nil, err
	}
	return persona, nil
}

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
