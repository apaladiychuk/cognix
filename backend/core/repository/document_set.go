package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type (
	DocumentSetRepository interface {
		FindByUser(ctx context.Context, userID uuid.UUID, param *parameters.ArchivedParam) ([]*model.DocumentSet, error)
		FindByID(ctx context.Context, userID uuid.UUID, id int64) (*model.DocumentSet, error)
		FindByIDWithConnectors(ctx context.Context, userID uuid.UUID, id int64) (*model.DocumentSet, error)
		Create(ctx context.Context, set *model.DocumentSet) error
		Update(ctx context.Context, set *model.DocumentSet) error
		AddConnector(ctx context.Context, pairs ...*model.DocumentSetConnectorPair) error
		DeleteConnector(ctx context.Context, documentSetID int64, connectorIDs []int64) error
	}
	documentSetRepository struct {
		db *pg.DB
	}
)

func (r *documentSetRepository) FindByUser(ctx context.Context, userID uuid.UUID, param *parameters.ArchivedParam) ([]*model.DocumentSet, error) {
	documentSets := make([]*model.DocumentSet, 0)
	stm := r.db.WithContext(ctx).Model(&documentSets).
		Where("user_id = ?", userID)
	if !param.Archived {
		stm = stm.Where("deleted_date = null")
	}
	if err := stm.Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find document sets")
	}
	return documentSets, nil
}

func (r *documentSetRepository) FindByID(ctx context.Context, userID uuid.UUID, id int64) (*model.DocumentSet, error) {
	var documentSet model.DocumentSet
	if err := r.db.WithContext(ctx).Model(&documentSet).
		Where("user_id = ?", userID).
		Where("id = ?", id).
		First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find document set")
	}
	return &documentSet, nil
}

func (r *documentSetRepository) FindByIDWithConnectors(ctx context.Context, userID uuid.UUID, id int64) (*model.DocumentSet, error) {
	var documentSet model.DocumentSet
	if err := r.db.WithContext(ctx).Model(&documentSet).
		Where("document_sets.user_id = ?", userID).
		Where("document_sets.id = ?", id).
		Relation("Pairs").
		First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find document set")
	}
	return &documentSet, nil
}

func (r *documentSetRepository) Create(ctx context.Context, set *model.DocumentSet) error {
	if _, err := r.db.WithContext(ctx).Model(set).Insert(); err != nil {
		return utils.Internal.Wrap(err, "can not create document set")
	}
	return nil
}

func (r *documentSetRepository) Update(ctx context.Context, set *model.DocumentSet) error {
	if _, err := r.db.WithContext(ctx).Model(set).
		Where("id = ?", set.ID).Update(); err != nil {
		return utils.Internal.Wrap(err, "can not update document set")
	}
	return nil
}
func (r *documentSetRepository) AddConnector(ctx context.Context, pairs ...*model.DocumentSetConnectorPair) error {
	if _, err := r.db.WithContext(ctx).Model(&pairs).Insert(); err != nil {
		return utils.Internal.Wrap(err, "can not add connector to document set")
	}
	return nil
}

func (r *documentSetRepository) DeleteConnector(ctx context.Context, documentSetID int64, connectorIDs []int64) error {
	if _, err := r.db.WithContext(ctx).Model(&model.DocumentSetConnectorPair{}).
		Where("document_set_id = ?", documentSetID).
		Where("connector_id in (?)", pg.Array(connectorIDs)).
		Delete(); err != nil {
		return utils.Internal.Wrap(err, "can not delete connector from document set")
	}
	return nil
}

func NewDocumentSetRepository(db *pg.DB) DocumentSetRepository {
	return &documentSetRepository{db: db}
}
