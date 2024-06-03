package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"context"
	"fmt"
)

type (
	Web struct {
		Base
		param *WebParameters
	}
	WebParameters struct {
		URL              string `url:"url"`
		SiteMap          string `json:"site_map"`
		SearchForSitemap bool   `json:"search_for_sitemap"`
	}
)

func (c *Web) PrepareTask(ctx context.Context, task Task) error {

	// if this connector new we need to run connectorTask for prepare document table
	if len(c.model.Docs) == 0 {
		doc, ok := c.model.DocsMap[c.param.URL]
		if !ok {
			doc = &model.Document{
				SourceID:    c.param.URL,
				ConnectorID: c.Base.model.ID,
				URL:         c.param.URL,
				Signature:   "",
			}
			c.model.Docs = append(c.model.Docs, doc)
		}
	}
	var rootDoc *model.Document
	for _, doc := range c.model.Docs {
		if !doc.ParentID.Valid {
			rootDoc = doc
			break
		}
	}
	if rootDoc == nil {
		return fmt.Errorf("root document not found")
	}
	return task.RunSemantic(ctx, &proto.SemanticData{
		Url:              c.param.URL,
		SiteMap:          c.param.SiteMap,
		SearchForSitemap: c.param.SearchForSitemap,
		DocumentId:       rootDoc.ID.IntPart(),
		FileType:         proto.FileType_URL,
		CollectionName:   c.model.CollectionName(),
		ModelName:        c.model.User.EmbeddingModel.ModelID,
		ModelDimension:   int32(c.model.User.EmbeddingModel.ModelDim),
	})
}

func (c *Web) Execute(ctx context.Context, param map[string]string) chan *Response {
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
			URL:              c.param.URL,
			SourceID:         c.param.URL,
			SiteMap:          c.param.SiteMap,
			SearchForSitemap: c.param.SearchForSitemap,
			DocumentID:       doc.ID.IntPart(),
			MimeType:         mineURL,
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
