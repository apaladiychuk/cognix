package connector

import (
	"cognix.ch/api/v2/core/model"
	"context"
)

type Web struct {
	url            string
	collectionName string
	model          *model.Connector
}

func (c *Web) Execute(ctx context.Context, param model.JSONMap) error {
	//TODO implement me
	panic("implement me")
}

func NewWeb(connector *model.Connector) Connector {
	return &Web{}
}
