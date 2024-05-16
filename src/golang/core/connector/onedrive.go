package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"context"
)

type (
	OneDrive struct {
		Base
		param *OneDriveParameters
		ctx   context.Context
	}
	OneDriveParameters struct {
		SubscriptionID string `json:"subscription_id"`
	}
)

func (c *OneDrive) Execute(ctx context.Context, param map[string]string) chan *proto.ChunkingData {
	chResult := make(chan *proto.ChunkingData, 1)

	return chResult
}

func NewOneDrive(connector *model.Connector) (Connector, error) {
	conn := OneDrive{}
	conn.Base.Config(connector)
	conn.param = &OneDriveParameters{}
	if err := connector.ConnectorSpecificConfig.ToStruct(conn.param); err != nil {
		return nil, err
	}

	return &conn, nil
}
