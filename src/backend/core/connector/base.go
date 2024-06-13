package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"context"
	"strings"
	"time"
)

const (
	mineURL        = "url"
	mimeYoutube    = "youtube"
	ParamFileLimit = "file_limit"
	GB             = 1024 * 1024 * 1024
)

var supportedMimeTypes = map[string]proto.FileType{
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":       proto.FileType_XLS,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": proto.FileType_DOC,
	"application/pdf": proto.FileType_PDF,
	"text/plain":      proto.FileType_TXT,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": proto.FileType_PPT,
}

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
	MimeType    string
	Signature   string
	Bucket      string
	SaveContent bool
	UpToData    bool
}

type Builder struct {
	connectorRepo repository.ConnectorRepository
}

func (r *Response) GetType() proto.FileType {
	switch r.MimeType {
	case mineURL:
		return proto.FileType_URL
	}
	mimeType := strings.Split(r.MimeType, ";")

	if fileType, ok := supportedMimeTypes[mimeType[0]]; ok {
		return fileType
	}
	return proto.FileType_UNKNOWN
}

type nopConnector struct {
	Base
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
