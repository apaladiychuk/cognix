package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"context"
	"time"
)

const mineURL = "url"

var supportedMimeTypes = map[string]proto.FileType{
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":       proto.FileType_XLS,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": proto.FileType_DOC,
	"application/pdf": proto.FileType_PDF,
	"text/rtf":        proto.FileType_RTF,
}

type Base struct {
	model    *model.Connector
	resultCh chan *Response
}
type Response struct {
	URL         string
	Name        string
	SourceID    string
	DocumentID  int64
	Content     []byte
	MimeType    string
	SaveContent bool
}
type Connector interface {
	Execute(ctx context.Context, param map[string]string) chan *Response
}

type Builder struct {
	connectorRepo repository.ConnectorRepository
}

func (r *Response) GetType() proto.FileType {
	switch r.MimeType {
	case mineURL:
		return proto.FileType_URL
	}

	if fileType, ok := supportedMimeTypes[r.MimeType]; ok {
		return fileType
	}
	return proto.FileType_URL
}

type nopConnector struct {
	Base
}

func (n *nopConnector) Execute(ctx context.Context, param map[string]string) chan *Response {
	ch := make(chan *Response)
	go func() {
		time.Sleep(1 * time.Second)
		close(ch)
	}()

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
	b.resultCh = make(chan *Response, 10)
	return
}
