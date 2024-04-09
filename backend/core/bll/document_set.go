package bll

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"time"
)

type (
	DocumentSetBL interface {
		GetByUser(ctx context.Context, user *model.User) ([]*model.DocumentSet, error)
		GetByID(ctx context.Context, user *model.User, id int64) (*model.DocumentSet, error)
		Create(ctx context.Context, user *model.User, name, description string) (*model.DocumentSet, error)
		Update(ctx context.Context, user *model.User, id int64, name, description string) (*model.DocumentSet, error)
		AddConnector(ctx context.Context, user *model.User, documentSetID int64, connectorIDs ...int64) ([]*model.DocumentSetConnectorPair, error)
		DeleteConnector(ctx context.Context, user *model.User, ids ...int64) error
	}
	documentSetBL struct {
		documentSetRepo repository.DocumentSetRepository
	}
)

func (b *documentSetBL) GetByUser(ctx context.Context, user *model.User) ([]*model.DocumentSet, error) {
	return b.documentSetRepo.FindByUser(ctx, user.ID)
}

func (b *documentSetBL) GetByID(ctx context.Context, user *model.User, id int64) (*model.DocumentSet, error) {
	return b.documentSetRepo.FindByID(ctx, user.ID, id)
}

func (b *documentSetBL) Create(ctx context.Context, user *model.User, name, description string) (*model.DocumentSet, error) {
	documentSet := model.DocumentSet{
		UserID:      uuid.New(),
		Name:        name,
		Description: description,
		IsUpToDate:  false,
		CreatedDate: time.Now().UTC(),
	}
	if err := b.documentSetRepo.Create(ctx, &documentSet); err != nil {
		return nil, err
	}
	return &documentSet, nil
}

func (b *documentSetBL) Update(ctx context.Context, user *model.User, id int64, name, description string) (*model.DocumentSet, error) {
	documentSet, err := b.documentSetRepo.FindByID(ctx, user.ID, id)
	if err != nil {
		return nil, err
	}
	documentSet.Name = name
	documentSet.Description = description
	documentSet.UpdatedDate = pg.NullTime{Time: time.Now().UTC()}
	if err = b.documentSetRepo.Update(ctx, documentSet); err != nil {
		return nil, err
	}
	return documentSet, nil
}

// AddConnector
func (b *documentSetBL) AddConnector(ctx context.Context, user *model.User, documentSetID int64, connectorIDs ...int64) ([]*model.DocumentSetConnectorPair, error) {
	documentSet, err := b.documentSetRepo.FindByIDWithConnectors(ctx, user.ID, documentSetID)
	if err != nil {
		return nil, err
	}
	existingContainers := make(map[int64]*model.DocumentSetConnectorPair)
	for _, pair := range documentSet.Pairs {
		existingContainers[pair.ConnectorID] = pair
	}
	var newPairs []*model.DocumentSetConnectorPair
	for _, connectorID := range connectorIDs {
		if _, ok := existingContainers[connectorID]; ok {
			continue
		}
		newPairs = append(newPairs, &model.DocumentSetConnectorPair{
			DocumentSetID: documentSetID,
			ConnectorID:   connectorID,
			IsCurrent:     false,
		})
	}
	if err = b.documentSetRepo.AddConnector(ctx, newPairs...); err != nil {
		return nil, err
	}
	return newPairs, nil
}

func (b *documentSetBL) DeleteConnector(ctx context.Context, user *model.User, ids ...int64) error {
	//TODO implement me
	panic("implement me")
}

func NewDocumentSetBL(documentSetRepo repository.DocumentSetRepository) DocumentSetBL {
	return &documentSetBL{documentSetRepo: documentSetRepo}
}
