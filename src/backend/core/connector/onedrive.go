package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/gabriel-vasile/mimetype"
	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
	"strconv"

	//"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"time"
)

// todo max file size 1G
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
	}
	OneDriveParameters struct {
		FolderName string       `json:"folder_name"`
		Token      oauth2.Token `json:"token"`
	}
)

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
	if size, ok := param[ParamFileLimit]; ok {
		fileSizeLimit, _ = strconv.Atoi(size)
	}
	if fileSizeLimit == 0 {
		fileSizeLimit = 1
	}
	c.fileSizeLimit = fileSizeLimit * GB

	if len(c.model.DocsMap) == 0 {
		c.model.DocsMap = make(map[string]*model.Document)
	}
	go c.execute(ctx)
	return c.resultCh
}

func (c *OneDrive) execute(ctx context.Context) {
	defer func() {
		close(c.resultCh)
	}()
	body, err := c.request(ctx, getDrive)
	if err != nil {
		zap.S().Errorf(err.Error())
		time.Sleep(50 * time.Millisecond)
		return
	}
	if body != nil {
		if err := c.handleItems(ctx, body.Value); err != nil {
			zap.S().Errorf(err.Error())
		}
	}
}

func (c *OneDrive) getFile(item *DriveChildBody) error {
	doc, ok := c.Base.model.DocsMap[item.Id]
	if !ok {
		doc = &model.Document{
			SourceID:    item.Id,
			ConnectorID: c.Base.model.ID,
			Link:        item.MicrosoftGraphDownloadUrl,
			Signature:   "",
		}
		c.Base.model.DocsMap[item.Id] = doc
	}
	doc.IsExists = true
	// todo need to clarify how to check that file was chunked, embedded and stored in milvus
	//if doc.Signature == item.File.Hashes.QuickXorHash {
	//	return nil
	//}
	doc.Signature = item.File.Hashes.QuickXorHash
	payload := &Response{
		URL:         item.MicrosoftGraphDownloadUrl,
		SourceID:    item.Id,
		Name:        item.Name,
		DocumentID:  doc.ID.IntPart(),
		MimeType:    item.File.MimeType,
		SaveContent: false,
	}
	// do not process files that size greater than limit
	if item.Size > c.fileSizeLimit {
		return nil
	}
	// try to recognize type of file by content
	if _, ok = supportedMimeTypes[item.File.MimeType]; !ok {
		response, err := c.client.R().
			SetDoNotParseResponse(true).
			Get(item.MicrosoftGraphDownloadUrl)
		if err == nil && !response.IsError() {
			if mime, err := mimetype.DetectReader(response.RawBody()); err == nil {
				payload.MimeType = mime.String()
			}
		}
		response.RawBody().Close()
	}
	if payload.GetType() == proto.FileType_UNKNOWN {
		zap.S().Infof("unsupported file %s type %s -- %s", item.Name, item.File.MimeType, payload.MimeType)
		return nil
	}
	c.resultCh <- payload
	return nil
}

func (c *OneDrive) getFolder(ctx context.Context, id string) error {
	body, err := c.request(ctx, fmt.Sprintf(getFolderChild, id))
	if err != nil {
		return err
	}
	return c.handleItems(ctx, body.Value)
}

func (c *OneDrive) handleItems(ctx context.Context, items []*DriveChildBody) error {
	for _, item := range items {
		if item.Folder != nil {
			if err := c.getFolder(ctx, item.Id); err != nil {
				zap.S().Errorf("Failed to get folder with id %s : %s ", item.Id, err.Error())
				continue
			}
		}
		if item.File != nil {

			if err := c.getFile(item); err != nil {
				zap.S().Errorf("Failed to get file with id %s : %s ", item.Id, err.Error())
				continue
			}
		}
	}
	return nil
}

func (c *OneDrive) request(ctx context.Context, url string) (*GetDriveResponse, error) {
	response, err := c.client.R().
		SetContext(ctx).
		SetHeader(authorizationHeader, fmt.Sprintf("%s %s",
			c.param.Token.TokenType,
			c.param.Token.AccessToken)).
		Get(url)
	if err != nil || response.IsError() {
		zap.S().Errorw("Error executing OneDrive", "error", err, "response", response)
		return nil, fmt.Errorf("%v:%v", err, response.Error())
	}
	var body GetDriveResponse
	if err = json.Unmarshal(response.Body(), &body); err != nil {
		zap.S().Errorw("unmarshal failed", "error", err)
		return nil, err
	}
	return &body, nil
}

// NewOneDrive creates new instance of OneDrive connector
func NewOneDrive(connector *model.Connector) (Connector, error) {
	conn := OneDrive{}
	conn.Base.Config(connector)
	conn.param = &OneDriveParameters{}
	if err := connector.ConnectorSpecificConfig.ToStruct(conn.param); err != nil {
		return nil, err
	}

	conn.client = resty.New().
		SetTimeout(time.Minute)

	return &conn, nil
}
