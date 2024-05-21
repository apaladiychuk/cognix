package connector

import (
	"cognix.ch/api/v2/core/model"
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/go-resty/resty/v2"

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
)

type (
	OneDrive struct {
		Base
		param  *OneDriveParameters
		ctx    context.Context
		client *resty.Client
	}
	OneDriveParameters struct {
		FolderName string `json:"folder_name"`
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
	File                      *File     `json:"file"`
	Size                      int       `json:"size"`
	Folder                    *Folder   `json:"folder"`
}

type File struct {
	Hashes struct {
		QuickXorHash string `json:"quickXorHash"`
	} `json:"hashes"`
	MimeType string `json:"mimeType"`
}

type Folder struct {
	ChildCount int `json:"childCount"`
}

func (c *OneDrive) Execute(ctx context.Context, param map[string]string) chan *Response {
	if len(c.model.DocsMap) == 0 {
		c.model.DocsMap = make(map[string]*model.Document)
	}
	go c.execute(ctx)
	return c.resultCh
}

func (c *OneDrive) execute(ctx context.Context) {
	defer close(c.resultCh)
	body, err := c.request(ctx, getDrive)
	if err != nil {
		zap.S().Errorf(err.Error())
	}
	c.handleItems(ctx, body.Value)
}

func (c *OneDrive) getFile(ctx context.Context, item *DriveChildBody) error {
	doc, ok := c.Base.model.DocsMap[item.Id]
	if !ok {
		doc = &model.Document{
			DocumentID:  item.Id,
			ConnectorID: c.Base.model.ID,
			Link:        "",
			Signature:   "",
		}
		c.Base.model.DocsMap[item.Id] = doc
	}
	if doc.Signature == item.File.Hashes.QuickXorHash {
		return nil
	}
	doc.Signature = item.File.Hashes.QuickXorHash
	response, err := c.client.R().
		SetContext(ctx).
		SetHeader(authorizationHeader, fmt.Sprintf("%s %s",
			c.model.Credential.CredentialJson.Token.TokenType,
			c.model.Credential.CredentialJson.Token.AccessToken)).
		Get(item.MicrosoftGraphDownloadUrl)
	if err != nil || response.IsError() {
		return fmt.Errorf("[%v] %v", err, response.Error())
	}
	c.resultCh <- &Response{
		URL:         "",
		SourceID:    item.Id,
		Name:        item.Name,
		DocumentID:  doc.ID.IntPart(),
		Content:     response.Body(),
		MimeType:    item.File.MimeType,
		SaveContent: true,
	}
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

			if err := c.getFile(ctx, item); err != nil {
				zap.S().Errorf("Failed to get file content with id %s : %s ", item.Id, err.Error())
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
			c.model.Credential.CredentialJson.Token.TokenType,
			c.model.Credential.CredentialJson.Token.AccessToken)).
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
