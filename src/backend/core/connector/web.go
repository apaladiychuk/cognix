package connector

import (
	"cognix.ch/api/v2/core/model"
	"context"
)

type (
	Web struct {
		Base
		param *WebParameters
		ctx   context.Context
	}
	WebParameters struct {
		URL string `url:"url"`
	}
)

func (c *Web) Execute(ctx context.Context, param map[string]string) chan *Response {
	c.ctx = ctx
	go func() {
		doc, ok := c.model.DocsMap[c.param.URL]
		if !ok {
			doc = &model.Document{
				SourceID:    c.param.URL,
				ConnectorID: c.Base.model.ID,
				URL:         c.param.URL,
				Signature:   "",
			}
			c.Base.model.DocsMap[c.param.URL] = doc
		}
		c.resultCh <- &Response{
			URL:        c.param.URL,
			SourceID:   c.param.URL,
			DocumentID: doc.ID.IntPart(),
			MimeType:   mineURL,
		}
		close(c.resultCh)
	}()
	return c.resultCh
}

func NewWeb(connector *model.Connector) (Connector, error) {
	web := Web{}
	web.Base.Config(connector)
	web.param = &WebParameters{}
	if err := connector.ConnectorSpecificConfig.ToStruct(web.param); err != nil {
		return nil, err
	}

	return &web, nil
}
