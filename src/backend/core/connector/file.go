package connector

import (
	"cognix.ch/api/v2/core/model"
	"context"
)

type (
	File struct {
		Base
		param *FileParameters
		ctx   context.Context
	}
	FileParameters struct {
		URL string `url:"url"`
	}
)

func (c *File) Execute(ctx context.Context, param map[string]string) chan *Response {
	c.ctx = ctx
	go func() {
		defer close(c.resultCh)
		if c.param == nil || c.param.URL == "" {
			return
		}
		doc, ok := c.Base.model.DocsMap[c.param.URL]
		if !ok {
			doc = &model.Document{
				DocumentID:  c.param.URL,
				ConnectorID: c.Base.model.ID,
				Link:        c.param.URL,
				Signature:   "",
			}
			c.model.DocsMap[c.param.URL] = doc
		}
		doc.IsExists = true
		c.resultCh <- &Response{
			URL:         c.param.URL,
			SourceID:    c.param.URL,
			SaveContent: true,
			MimeType:    mineURL,
		}
	}()
	return c.resultCh
}

func NewFile(connector *model.Connector) (Connector, error) {
	fileConn := File{}
	fileConn.Base.Config(connector)
	fileConn.param = &FileParameters{}
	if err := connector.ConnectorSpecificConfig.ToStruct(fileConn.param); err != nil {
		return nil, err
	}

	return &fileConn, nil
}
