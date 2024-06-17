package connector

import (
	microsoft_core "cognix.ch/api/v2/core/connector/microsoft-core"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
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
		sessionID     uuid.NullUUID
	}
	OneDriveParameters struct {
		microsoft_core.MSDriveParam
		Token oauth2.Token `json:"token"`
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

//	func (c *OneDrive) execute(ctx context.Context) {
//		defer func() {
//			close(c.resultCh)
//		}()
//
//		body, err := c.request(ctx, getDrive)
//		if err != nil {
//			zap.S().Errorf(err.Error())
//			time.Sleep(50 * time.Millisecond)
//			return
//		}
//		if body != nil {
//			if err := c.handleItems(ctx, "", body.Value); err != nil {
//				zap.S().Errorf(err.Error())
//			}
//		}
//	}
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

//
//func (c *OneDrive) recognizeFiletype(item *DriveChildBody) (string, proto.FileType) {
//
//	mimeTypeParts := strings.Split(item.File.MimeType, ";")
//
//	if fileType, ok := supportedMimeTypes[mimeTypeParts[0]]; ok {
//		return mimeTypeParts[0], fileType
//	}
//	// recognize fileType by filename extension
//	fileNameParts := strings.Split(item.Name, ".")
//	if len(fileNameParts) > 1 {
//		if mimeType, ok := supportedExtensions[strings.ToUpper(fileNameParts[len(fileNameParts)-1])]; ok {
//			return mimeType, supportedMimeTypes[mimeType]
//		}
//	}
//	// recognize filetype by content
//	response, err := c.client.R().
//		SetDoNotParseResponse(true).
//		Get(item.MicrosoftGraphDownloadUrl)
//	if err == nil && !response.IsError() {
//		if mime, err := mimetype.DetectReader(response.RawBody()); err == nil {
//			if fileType, ok := supportedMimeTypes[mime.String()]; ok {
//				return mime.String(), fileType
//			}
//		}
//	}
//	response.RawBody().Close()
//	return "", proto.FileType_UNKNOWN
//}
//func (c *OneDrive) getFolder(ctx context.Context, folder string, id string) error {
//	body, err := c.request(ctx, fmt.Sprintf(getFolderChild, id))
//	if err != nil {
//		return err
//	}
//	return c.handleItems(ctx, folder, body.Value)
//}
//
//func (c *OneDrive) handleItems(ctx context.Context, folder string, items []*DriveChildBody) error {
//	for _, item := range items {
//		// read files if user do not configure folder name
//		// or current folder as a part of configured folder.
//		if !c.isFolderAnalysing(folder) {
//			continue
//		}
//		//if item.File != nil && (strings.Contains(folder, c.param.Folder) || c.param.Folder == "") {
//		if item.File != nil && c.isFilesAnalysing(folder) {
//			if err := c.getFile(item); err != nil {
//				zap.S().Errorf("Failed to get file with id %s : %s ", item.Id, err.Error())
//				continue
//			}
//		}
//		if item.Folder != nil {
//			// do not scan nested folder if user  wants to read dod from single folder
//			if /*item.Name != c.param.Folder*/ strings.Contains(folder, c.param.Folder) && !c.param.Recursive {
//				continue
//			}
//			nextFolder := folder
//			if nextFolder != "" {
//				nextFolder += "/"
//			}
//			if err := c.getFolder(ctx, nextFolder+item.Name, item.Id); err != nil {
//				zap.S().Errorf("Failed to get folder with id %s : %s ", item.Id, err.Error())
//				continue
//			}
//		}
//
//	}
//	return nil
//}
//
//func (c *OneDrive) request(ctx context.Context, url string) (*GetDriveResponse, error) {
//	response, err := c.client.R().
//		SetContext(ctx).
//		SetHeader(authorizationHeader, fmt.Sprintf("%s %s",
//			c.param.Token.TokenType,
//			c.param.Token.AccessToken)).
//		Get(url)
//	if err = utils.WrapRestyError(response, err); err != nil {
//		zap.S().Error(err.Error())
//		return nil, err
//	}
//	var body GetDriveResponse
//	if err = json.Unmarshal(response.Body(), &body); err != nil {
//		zap.S().Errorw("unmarshal failed", "error", err)
//		return nil, err
//	}
//	return &body, nil
//}

// NewOneDrive creates new instance of OneDrive connector
func NewOneDrive(connector *model.Connector) (Connector, error) {
	conn := OneDrive{}
	conn.Base.Config(connector)
	conn.param = &OneDriveParameters{}
	if err := connector.ConnectorSpecificConfig.ToStruct(conn.param); err != nil {
		return nil, err
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
