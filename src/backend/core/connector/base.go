package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"context"
	"time"
)

const ()

type Task interface {
	RunConnector(ctx context.Context, data *proto.ConnectorRequest) error
	RunSemantic(ctx context.Context, data *proto.SemanticData) error
	UpToDate(ctx context.Context) error
}

type Connector interface {
	Execute(ctx context.Context, param map[string]string) chan *Response
	PrepareTask(ctx context.Context, task Task) error
	Validate() error
}

type Base struct {
	model    *model.Connector
	resultCh chan *Response
}
type Response struct {
	URL              string
	SiteMap          string
	SearchForSitemap bool
	Name             string
	SourceID         string
	DocumentID       int64
	//Content          []byte
	MimeType  string
	FileType  proto.FileType
	Signature string
	Content   *Content
	UpToData  bool
}

// Content  defines action for stop content in minio database
type Content struct {
	Bucket        string
	URL           string // URL for download
	Body          []byte // Body raw content  for store
	AppendContent bool   // if true content will be added to existing file on minio
}

type Builder struct {
	connectorRepo repository.ConnectorRepository
}

type nopConnector struct {
	Base
}

func (n *nopConnector) Validate() error {
	return nil
}

func (n *nopConnector) PrepareTask(ctx context.Context, task Task) error {
	return nil
}

func (n *nopConnector) Execute(ctx context.Context, param map[string]string) chan *Response {
	ch := make(chan *Response)
	go func() {
		time.Sleep(5 * time.Second)
		close(ch)
	}()

	return ch
}

func New(connectorModel *model.Connector) (Connector, error) {
	switch connectorModel.Type {
	case model.SourceTypeWEB:
		return NewWeb(connectorModel)
	case model.SourceTypeOneDrive:
		return NewOneDrive(connectorModel)
	case model.SourceTypeFile:
		return NewFile(connectorModel)
	case model.SourceTypeYoutube:
		return NewYoutube(connectorModel)
	default:
		return &nopConnector{}, nil
	}
}

func (b *Base) Config(connector *model.Connector) {
	b.model = connector
	b.resultCh = make(chan *Response, 10)
	return
}
