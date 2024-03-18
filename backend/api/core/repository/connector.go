package repository

import "github.com/go-pg/pg/v10"

type (
	ConnectorRepository interface {
	}
	connectorRepository struct {
		db *pg.DB
	}
)

func NewConnectorRepository(db *pg.DB) ConnectorRepository {
	return &connectorRepository{db: db}
}
