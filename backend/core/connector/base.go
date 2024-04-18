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

func (b *Builder) New(ctx context.Context, connector *model.Connector) Connector {
	switch connector.Source {
	case model.SourceTypeWEB:
		return NewWeb(connector)
	}
}
