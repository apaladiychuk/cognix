package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"context"
	"fmt"
)

type (
	File struct {
		Base
		param *FileParameters
		ctx   context.Context
	}
	FileParameters struct {
		FileName string `json:"file_name"`
		MIMEType string `json:"mime_type"`
	}
)

func (c *File) PrepareTask(ctx context.Context, task Task) error {

	if len(c.model.Docs) == 0 {
		// send message to connector service if new.
		return task.RunConnector(ctx, &proto.ConnectorRequest{
			Id:     c.model.ID.IntPart(),
			Params: make(map[string]string),
		})
	}
	if !c.model.Docs[0].Analyzed {
		// if file is not  chunked and not stored in vector database send message to chunker
		link := fmt.Sprintf("minio:tenant-%s:%s", c.model.User.EmbeddingModel.TenantID.String(), c.param.FileName)
		return task.RunChunker(ctx, &proto.ChunkingData{
			Url:            link,
			DocumentId:     c.model.Docs[0].ID.IntPart(),
			FileType:       0,
			CollectionName: c.model.CollectionName(),
			ModelName:      c.model.User.EmbeddingModel.ModelID,
			ModelDimension: int32(c.model.User.EmbeddingModel.ModelDim),
		})
	}
	// file already chunked and  stored in vector database. Update connector status.
	return task.UpToDate(ctx)
}

func (c *File) Execute(ctx context.Context, param map[string]string) chan *Response {
	c.ctx = ctx
	go func() {
		defer close(c.resultCh)
		if c.param == nil || c.param.FileName == "" {
			return
		}
		// check id document  already exists
		doc, ok := c.Base.model.DocsMap[c.param.FileName]
		url := fmt.Sprintf("minio:tenant-%s:%s", c.model.User.EmbeddingModel.TenantID, c.param.FileName)
		if !ok {
			doc = &model.Document{
				SourceID:    url,
				ConnectorID: c.model.ID,
				URL:         url,
				Signature:   "",
			}
			c.model.DocsMap[url] = doc
		}
		doc.IsExists = true
		c.resultCh <- &Response{
			URL:      url,
			SourceID: url,
			MimeType: c.param.MIMEType,
		}
	}()
	return c.resultCh
}

// NewFile creates new instance of file connector.
func NewFile(connector *model.Connector) (Connector, error) {
	fileConn := File{}
	fileConn.Base.Config(connector)
	fileConn.param = &FileParameters{}
	if err := connector.ConnectorSpecificConfig.ToStruct(fileConn.param); err != nil {
		return nil, err
	}

	return &fileConn, nil
}
