package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"context"
)

type Connector interface {
	Config(ctx context.Context, connector *model.Connector) error
	Execute(ctx context.Context, param model.JSONMap) error
}

type Builder struct {
	connectorRepo repository.ConnectorRepository
}
