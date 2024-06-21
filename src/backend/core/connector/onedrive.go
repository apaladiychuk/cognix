package connector

import (
	microsoft_core "cognix.ch/api/v2/core/connector/microsoft-core"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"context"
	"fmt"
	_ "github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"strconv"
	"strings"

	//"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"time"
)

const (
	authorizationHeader = "Authorization"
	apiBase             = "https://graph.microsoft.com/v2.0"
	getFilesURL         = "/me/drive/special/%s/children"
	getDrive            = "https://graph.microsoft.com/v1.0/me/drive/root/children"
	getFolderChild      = "https://graph.microsoft.com/v1.0/me/drive/items/%s/children"
	createSharedLink    = "https://graph.microsoft.com/v1.0/me/drive/items/%s/createLink"
)

type (
	OneDrive struct {
		Base
		param         *OneDriveParameters
		ctx           context.Context
		client        *resty.Client
		fileSizeLimit int
		sessionID     uuid.NullUUID
	}
	OneDriveParameters struct {
		microsoft_core.MSDriveParam
		Token *oauth2.Token `json:"token"`
	}
)

func (c *OneDrive) Validate() error {
	return nil
}

type GetDriveResponse struct {
	Value []*DriveChildBody `json:"value"`
}

type DriveChildBody struct {
	MicrosoftGraphDownloadUrl string    `json:"@microsoft.graph.downloadUrl"`
	MicrosoftGraphDecorator   string    `json:"@microsoft.graph.Decorator"`
	Id                        string    `json:"id"`
	LastModifiedDateTime      time.Time `json:"lastModifiedDateTime"`
	Name                      string    `json:"name"`
	WebUrl                    string    `json:"webUrl"`
	File                      *MsFile   `json:"file"`
	Size                      int       `json:"size"`
	Folder                    *Folder   `json:"folder"`
}

type MsFile struct {
	Hashes struct {
		QuickXorHash string `json:"quickXorHash"`
	} `json:"hashes"`
	MimeType string `json:"mimeType"`
}

type Folder struct {
	ChildCount int `json:"childCount"`
}

func (c *OneDrive) PrepareTask(ctx context.Context, task Task) error {
	//	for one drive always send message to connector
	return task.RunConnector(ctx, &proto.ConnectorRequest{
		Id:     c.model.ID.IntPart(),
		Params: make(map[string]string),
	})
}

func (c *OneDrive) Execute(ctx context.Context, param map[string]string) chan *Response {
	var fileSizeLimit int
	c.sessionID = uuid.NullUUID{
		UUID:  uuid.New(),
		Valid: true,
	}
	if size, ok := param[model.ParamFileLimit]; ok {
		fileSizeLimit, _ = strconv.Atoi(size)
	}
	if fileSizeLimit == 0 {
		fileSizeLimit = 1
	}
	c.fileSizeLimit = fileSizeLimit * model.GB

	if len(c.model.DocsMap) == 0 {
		c.model.DocsMap = make(map[string]*model.Document)
	}
	msDrive := microsoft_core.NewMSDrive(
		&c.param.MSDriveParam,
		c.model,
		c.sessionID,
		c.client,
		getDrive,
		getFolderChild,
		c.getFile,
	)
	go func() {
		defer close(c.resultCh)
		if err := msDrive.Execute(ctx, param); err != nil {
			zap.S().Errorf(err.Error())
		}
	}()

	return c.resultCh
}

func (c *OneDrive) getFile(payload *microsoft_core.Response) {
	response := &Response{
		URL:        payload.URL,
		Name:       payload.Name,
		SourceID:   payload.SourceID,
		DocumentID: payload.DocumentID,
		MimeType:   payload.MimeType,
		FileType:   payload.FileType,
		Signature:  payload.Signature,
		Content: &Content{
			Bucket: model.BucketName(c.model.User.EmbeddingModel.TenantID),
			URL:    payload.URL,
		},
	}
	c.resultCh <- response
}

// NewOneDrive creates new instance of OneDrive connector
func NewOneDrive(connector *model.Connector,
	connectorRepo repository.ConnectorRepository,
	oauthURL string) (Connector, error) {
	conn := OneDrive{
		Base: Base{
			connectorRepo: connectorRepo,
			oauthClient: resty.New().
				SetTimeout(time.Minute).
				SetBaseURL(oauthURL),
		},
		param: &OneDriveParameters{},
	}

	conn.Base.Config(connector)
	if err := connector.ConnectorSpecificConfig.ToStruct(conn.param); err != nil {
		return nil, err
	}
	newToken, err := conn.refreshToken(conn.param.Token)
	if err != nil {
		return nil, err
	}
	if newToken != nil {
		conn.param.Token = newToken
	}

	conn.client = resty.New().
		SetTimeout(time.Minute).
		SetHeader(authorizationHeader, fmt.Sprintf("%s %s",
			conn.param.Token.TokenType,
			conn.param.Token.AccessToken))
	return &conn, nil
}
func (c *OneDrive) isFolderAnalysing(current string) bool {
	mask := c.param.Folder
	if len(current) < len(c.param.Folder) {
		mask = c.param.Folder[:len(current)]
	}
	// if user does not  set folder name. scan whole oneDrive or only root if recursive is false
	if c.param.Folder == "" {
		return len(current) == 0 || c.param.Recursive
	}
	// verify is current folder is   part of folder that user configure for scan
	if c.param.Recursive {
		return strings.HasPrefix(current+"/", mask+"/") || current == c.param.Folder
	}
	return strings.HasPrefix(current+"/", mask+"/") && len(current) <= len(c.param.Folder)
}

func (c *OneDrive) isFilesAnalysing(current string) bool {
	mask := c.param.Folder
	if len(current) < len(c.param.Folder) {
		mask = c.param.Folder[:len(mask)]
	}
	// if user does not  set folder name. scan whole oneDrive or only root if recursive is false
	if c.param.Folder == "" {
		return len(current) == 0 || c.param.Recursive

	}

	if c.param.Recursive {
		// recursive
		return strings.HasPrefix(current+"/", mask+"/") || current == c.param.Folder
	}
	// only one folder
	return current == c.param.Folder
}
