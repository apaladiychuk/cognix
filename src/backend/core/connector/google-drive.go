package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-pg/pg/v10"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	drive "google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"strconv"
	"strings"
	"time"

	"net/http"

	"google.golang.org/api/option"
)

const (
	googleDriveMIMEApplicationFile  = "application/vnd.google-apps"
	googleDriveMIMEFolder           = "application/vnd.google-apps.folder"
	googleDriveMIMEExternalShortcut = "application/vnd.google-apps.drive-sdk"
	googleDriveMIMEShortcut         = "application/vnd.google-apps.shortcut"

	googleDriveMIMEAudion       = "application/vnd.google-apps.audio"
	googleDriveMIMEDocument     = "application/vnd.google-apps.document"
	googleDriveMIMEPresentation = "application/vnd.google-apps.presentation"
	googleDriveMIMESpreadsheet  = "application/vnd.google-apps.spreadsheet"
	googleDriveMIMEVideo        = "application/vnd.google-apps.video"

	googleDriveRootFolderQuery = "(sharedWithMe or 'root' in parents)"
)

var googleDriveExportFileType = map[string]string{
	googleDriveMIMEDocument:     model.MIMETypeDOCX,
	googleDriveMIMEPresentation: model.MIMETypePPTX,
	googleDriveMIMESpreadsheet:  model.MIMETypeXLSX,
}

type (
	//
	GoogleDrive struct {
		Base
		param               *GoogleDriveParameters
		client              *drive.Service
		fileSizeLimit       int
		sessionID           uuid.NullUUID
		unsupportedMimeType map[string]bool
	}
	//
	GoogleDriveParameters struct {
		Folder    string        `json:"folder"`
		Recursive bool          `json:"recursive"`
		Token     *oauth2.Token `json:"token"`
	}
)

func (p GoogleDriveParameters) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Token, validation.By(func(value interface{}) error {
			if p.Token == nil {
				return fmt.Errorf("missing token")
			}
			if p.Token.AccessToken == "" || p.Token.RefreshToken == "" ||
				p.Token.TokenType == "" {
				return fmt.Errorf("wrong token")
			}
			return nil
		})),
	)
}

func (c *GoogleDrive) Execute(ctx context.Context, param map[string]string) chan *Response {

	var fileSizeLimit int
	if size, ok := param[model.ParamFileLimit]; ok {
		fileSizeLimit, _ = strconv.Atoi(size)
	}
	if fileSizeLimit == 0 {
		fileSizeLimit = 1
	}
	c.fileSizeLimit = fileSizeLimit * model.GB
	paramSessionID, _ := param[model.ParamSessionID]
	if uuidSessionID, err := uuid.Parse(paramSessionID); err != nil {
		c.sessionID = uuid.NullUUID{uuid.New(), true}
	} else {
		c.sessionID = uuid.NullUUID{uuidSessionID, true}
	}

	go func() {
		defer close(c.resultCh)
		folders := []string{""}
		if c.param.Folder != "" {
			rootFolder, err := c.getFolder(ctx)
			if err != nil {
				zap.S().Errorf("can not find folder %s: %s ", c.param.Folder, err.Error())
			}
			folders[0] = rootFolder
		}
		for len(folders) > 0 {
			folders = c.scanFolders(ctx, folders)
		}
		return
	}()
	return c.resultCh
}

func (c *GoogleDrive) scanFolders(ctx context.Context, folders []string) []string {
	nextFolders := make([]string, 0)
	for _, folder := range folders {
		childFolders, err := c.getFolderItems(ctx, folder)
		if err != nil {
			zap.S().Errorf("can not scan folder %s: %s ", folder, err.Error())
			continue
		}
		nextFolders = append(nextFolders, childFolders...)
	}
	return nextFolders
}

func (c *GoogleDrive) PrepareTask(ctx context.Context, sessionID uuid.UUID, task Task) error {
	return task.RunConnector(ctx, &proto.ConnectorRequest{
		Id: c.model.ID.IntPart(),
		Params: map[string]string{
			model.ParamSessionID: sessionID.String(),
		},
	})
}

func (c *GoogleDrive) Validate() error {
	if c.param == nil {
		return fmt.Errorf("file parameter is required")
	}
	return c.param.Validate()
}

func (c *GoogleDrive) getFolderItems(ctx context.Context, folderID string) ([]string, error) {
	var q string
	if folderID == "" {
		q = googleDriveRootFolderQuery
	} else {
		q = fmt.Sprintf(" '%s' in parents ", folderID)
	}
	nextFolders := make([]string, 0)
	var fields googleapi.Field = "nextPageToken, files(name,id, exportLinks, size, mimeType,webContentLink ,fileExtension,md5Checksum, version ) "
	if err := c.client.Files.List().Context(ctx).Q(q).Fields(fields).Pages(ctx, func(l *drive.FileList) error {
		for _, item := range l.Files {
			if item.MimeType == googleDriveMIMEFolder {
				if c.param.Recursive {
					nextFolders = append(nextFolders, item.Id)
				}
				continue
			}
			if err := c.scanFile(ctx, item); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return nextFolders, nil
}

func (c *GoogleDrive) scanFile(ctx context.Context, item *drive.File) error {

	if item.Size > int64(c.fileSizeLimit) {
		return nil
	}
	mimeType, fileType := c.recognizeFiletype(item)
	if fileType == proto.FileType_UNKNOWN {
		return nil
	}
	url := item.WebContentLink
	if len(item.ExportLinks) > 0 {
		url = item.ExportLinks[mimeType]
	}
	doc, ok := c.model.DocsMap[item.Id]
	if !ok {
		doc = &model.Document{
			SourceID:        item.Id,
			ConnectorID:     c.model.ID,
			URL:             url,
			Signature:       "",
			ChunkingSession: c.sessionID,
			CreationDate:    time.Now().UTC(),
			LastUpdate:      pg.NullTime{time.Now().UTC()},
			OriginalURL:     url,
		}
		c.model.DocsMap[item.Id] = doc
	}
	doc.IsExists = true
	checksum := item.Md5Checksum
	if checksum == "" {
		checksum = fmt.Sprintf("%d", item.Version)
	}
	if doc.Signature == checksum {
		return nil
	}
	doc.Signature = checksum

	filename := utils.StripFileName(uuid.New().String() + item.Name)
	response := &Response{
		URL:        url,
		Name:       filename,
		SourceID:   item.Id,
		DocumentID: doc.ID.IntPart(),
		MimeType:   mimeType,
		FileType:   fileType,
		Signature:  doc.Signature,
		Content: &Content{
			Bucket: model.BucketName(c.model.User.EmbeddingModel.TenantID),
		},
	}
	var resp *http.Response
	var err error
	if len(item.ExportLinks) > 0 {
		resp, err = c.client.Files.Export(item.Id, mimeType).Context(ctx).Download()
	} else {
		resp, err = c.client.Files.Get(item.Id).Context(ctx).Download()
	}
	if err != nil {
		zap.S().Errorf("can not download file %s  : %s", item.OriginalFilename, err.Error())
		return nil
	}
	response.Content.Reader = resp.Body
	c.resultCh <- response
	return nil
}
func (c *GoogleDrive) getFolder(ctx context.Context) (string, error) {
	folderParts := strings.Split(c.param.Folder, "/")
	parentID := ""
	for i, part := range folderParts {
		q := fmt.Sprintf("name = '%s' and mimeType = '%s'", part, googleDriveMIMEFolder)
		if parentID == "" {
			//  find in root
			q += " and " + googleDriveRootFolderQuery
		} else {
			// find in folder
			q += fmt.Sprintf(" and '%s' in parents ", parentID)
		}

		folder, err := c.client.Files.List().Context(ctx).Q(q).Do()
		if err != nil {
			return "", err
		}

		if len(folder.Files) == 0 {
			return "", fmt.Errorf("folder %s not found", strings.Join(folderParts[:i], "/"))
		}
		parentID = folder.Files[0].Id
	}
	if parentID == "" {
		return "", fmt.Errorf("folder %s not found", c.param.Folder)
	}
	return parentID, nil

}

func (c *GoogleDrive) recognizeFiletype(item *drive.File) (string, proto.FileType) {
	if item.MimeType == googleDriveMIMEShortcut {
		return "", proto.FileType_UNKNOWN
	}
	if item.FileExtension != "" {
		if mimeType, ok := model.SupportedExtensions[strings.ToUpper(item.FileExtension)]; ok {
			return mimeType, model.SupportedMimeTypes[mimeType]
		}
	}
	if _, ok := c.unsupportedMimeType[item.MimeType]; ok {
		return "", proto.FileType_UNKNOWN
	}
	// recognize file type for google application file
	if mimeType, ok := googleDriveExportFileType[item.MimeType]; ok {
		return mimeType, model.SupportedMimeTypes[mimeType]
	}
	// recognize by mime type
	if ft, ok := model.SupportedMimeTypes[item.MimeType]; ok {
		return item.MimeType, ft
	}
	c.unsupportedMimeType[item.MimeType] = true
	zap.S().Errorf("Unsupported file type: %s  %s", item.OriginalFilename, item.MimeType)
	return "", proto.FileType_UNKNOWN
}

func NewGoogleDrive(connector *model.Connector,
	connectorRepo repository.ConnectorRepository,
	oauthURL string) (Connector, error) {
	conn := GoogleDrive{
		Base: Base{
			connectorRepo: connectorRepo,
			oauthClient: resty.New().
				SetTimeout(time.Minute).
				SetBaseURL(oauthURL),
		},
		param:               &GoogleDriveParameters{},
		unsupportedMimeType: make(map[string]bool),
	}

	conn.Base.Config(connector)
	if err := connector.ConnectorSpecificConfig.ToStruct(conn.param); err != nil {
		return nil, err
	}
	if err := conn.Validate(); err != nil {
		return nil, err
	}
	newToken, err := conn.refreshToken(conn.param.Token)
	if err != nil {
		return nil, err
	}
	if newToken != nil {
		conn.param.Token = newToken
	}
	client, err := drive.NewService(context.Background(),
		option.WithHTTPClient(&http.Client{Transport: utils.NewTransport(conn.param.Token)}))
	if err != nil {
		return nil, err
	}

	conn.client = client
	return &conn, nil
}
