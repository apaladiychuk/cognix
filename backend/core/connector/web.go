package connector

import (
	"cognix.ch/api/v2/core/model"
	"context"
)

type (
	Web struct {
		Base
		param *WebParameters
	}
	WebParameters struct {
		URL string `url:"url"`
	}
)

func (c *Web) Config(connector *model.Connector) (Connector, error) {
	c.Base.Config(connector)
	if err := connector.ConnectorSpecificConfig.ToStruct(c.param); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Web) Execute(ctx context.Context, param model.JSONMap) error {
	//TODO implement me
	panic("implement me")
}

func NewWeb(connector *model.Connector) (Connector, error) {
	var web Web
	return web.Config(connector)
}
