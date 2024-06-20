package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
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
	model         *model.Connector
	connectorRepo repository.ConnectorRepository
	resultCh      chan *Response
	oauthClient   *resty.Client
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

func New(connectorModel *model.Connector,
	connectorRepo repository.ConnectorRepository,
	oauthURL string) (Connector, error) {
	switch connectorModel.Type {
	case model.SourceTypeWEB:
		return NewWeb(connectorModel)
	case model.SourceTypeOneDrive:
		return NewOneDrive(connectorModel, connectorRepo, oauthURL)
	case model.SourceTypeFile:
		return NewFile(connectorModel)
	case model.SourceTypeYoutube:
		return NewYoutube(connectorModel)
	case model.SourceTypeMsTeams:
		return NewMSTeams(connectorModel, connectorRepo, oauthURL)
	default:
		return &nopConnector{}, nil
	}
}

func (b *Base) Config(connector *model.Connector) {
	b.model = connector
	b.resultCh = make(chan *Response, 10)
	return
}

// refreshToken  refresh OAuth token and store credential in database
func (b *Base) refreshToken(token *oauth2.Token) (*oauth2.Token, error) {
	if token.Expiry.UTC().After(time.Now().UTC()) {
		return nil, nil
	}
	provider, ok := model.ConnectorAuthProvider[b.model.Type]
	if !ok {
		return nil, nil
	}

	response, err := b.oauthClient.R().
		SetBody(token).Post(fmt.Sprintf("/api/oauth/%s/refresh_token", provider))
	if err = utils.WrapRestyError(response, err); err != nil {
		return nil, err
	}
	var payload struct {
		Data *oauth2.Token `json:"data"`
	}

	if err = json.Unmarshal(response.Body(), &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshl token: %v : %v", err, response.Error())
	}
	b.model.ConnectorSpecificConfig["token"] = payload.Data
	if err = b.connectorRepo.Update(context.Background(), b.model); err != nil {
		return nil, err
	}
	return payload.Data, nil
}
