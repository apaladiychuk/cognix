package bll

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"time"
)

type (
	PersonaBL interface {
		GetAll(ctx context.Context, user *model.User, archived bool) ([]*model.Persona, error)
		GetByID(ctx context.Context, user *model.User, id int64) (*model.Persona, error)
		Create(ctx context.Context, user *model.User, param *parameters.PersonaParam) (*model.Persona, error)
		Update(ctx context.Context, id int64, user *model.User, param *parameters.PersonaParam) (*model.Persona, error)
		Archive(ctx context.Context, user *model.User, id int64, restore bool) (*model.Persona, error)
	}
	personaBL struct {
		personaRepo repository.PersonaRepository
	}
)

func (b *personaBL) Archive(ctx context.Context, user *model.User, id int64, restore bool) (*model.Persona, error) {
	if !user.HasRoles(model.RoleSuperAdmin, model.RoleAdmin) {
		return nil, utils.ErrorPermission.New("do not have permission")
	}
	var relations []string

	if !restore {
		relations = append(relations, "ChatSessions")
	}
	persona, err := b.personaRepo.GetByID(ctx, id, user.TenantID, relations...)
	if err != nil {
		return nil, err
	}
	if len(persona.ChatSessions) > 0 {
		return nil, utils.ErrorBadRequest.New("persona is used in chat sessions")
	}
	if restore {
		persona.DeletedDate = pg.NullTime{}
	} else {
		persona.DeletedDate = pg.NullTime{time.Now().UTC()}
	}
	persona.UpdatedDate = pg.NullTime{time.Now().UTC()}
	if err = b.personaRepo.Update(ctx, persona); err != nil {
		return nil, err
	}
	return persona, nil
}

func (b *personaBL) Create(ctx context.Context, user *model.User, param *parameters.PersonaParam) (*model.Persona, error) {

	starterMessages, err := json.Marshal(param.StarterMessages)
	if err != nil {
		return nil, utils.ErrorBadRequest.Wrap(err, "fail to marshal starter messages")
	}
	persona := model.Persona{
		Name:            param.Name,
		DefaultPersona:  true,
		Description:     param.Description,
		TenantID:        user.TenantID,
		IsVisible:       true,
		StarterMessages: starterMessages,
		CreatedDate:     time.Now().UTC(),
		LLM: &model.LLM{
			Name:        fmt.Sprintf("%s %s", user.FirstName, param.ModelID),
			ModelID:     param.ModelID,
			TenantID:    user.TenantID,
			CreatedDate: time.Now().UTC(),
			Url:         param.URL,
			ApiKey:      param.APIKey,
			Endpoint:    param.Endpoint,
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
	starterMessages, err := json.Marshal(param.StarterMessages)
	if err != nil {
		return nil, utils.ErrorBadRequest.Wrap(err, "fail to marshal starter messages")
	}
	persona.Name = param.Name
	persona.Description = param.Description
	persona.UpdatedDate = pg.NullTime{time.Now().UTC()}
	persona.StarterMessages = starterMessages
	persona.LLM.Endpoint = param.Endpoint
	persona.LLM.ModelID = param.ModelID
	persona.LLM.ApiKey = param.APIKey
	persona.LLM.UpdatedDate = pg.NullTime{time.Now().UTC()}
	persona.Prompt.Name = param.Name
	persona.Prompt.Description = param.Description
	persona.Prompt.SystemPrompt = param.SystemPrompt
	persona.Prompt.TaskPrompt = param.TaskPrompt
	persona.Prompt.UpdatedDate = pg.NullTime{time.Now().UTC()}

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
func (b *personaBL) GetAll(ctx context.Context, user *model.User, archived bool) ([]*model.Persona, error) {
	return b.personaRepo.GetAll(ctx, user.TenantID, archived)
}

func (b *personaBL) GetByID(ctx context.Context, user *model.User, id int64) (*model.Persona, error) {
	return b.personaRepo.GetByID(ctx, id, user.TenantID)
}
