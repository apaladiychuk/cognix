package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"context"
)

type Base struct {
	model    *model.Connector
	resultCh chan *proto.ChunkingData
}

type Connector interface {
	Execute(ctx context.Context, param map[string]string) chan *proto.ChunkingData
}

type Builder struct {
	connectorRepo repository.ConnectorRepository
}

type nopConnector struct {
	Base
}

func (n *nopConnector) Execute(ctx context.Context, param map[string]string) chan *proto.ChunkingData {
	ch := make(chan *proto.ChunkingData)
	return ch
}

func New(connectorModel *model.Connector) (Connector, error) {
	switch connectorModel.Source {
	case model.SourceTypeWEB:
		return NewWeb(connectorModel)
	case model.SourceTypeOneDrive:
		return NewOneDrive(connectorModel)
	default:
		return &nopConnector{}, nil
	}
}

func (b *Base) Config(connector *model.Connector) {
	b.model = connector
	b.resultCh = make(chan *proto.ChunkingData, 10)
	return
}
