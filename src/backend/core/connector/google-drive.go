package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	drive "google.golang.org/api/drive/v2"
	"strings"
	"time"

	"net/http"

	"google.golang.org/api/option"
	"log"
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

type (
	//
	GoogleDrive struct {
		Base
		param         *GoogleDriveParameters
		client        *drive.Service
		fileSizeLimit int
		sessionID     uuid.NullUUID
	}
	//
	GoogleDriveParameters struct {
		Folder    string        `json:"folder"`
		Recursive bool          `json:"recursive"`
		Token     *oauth2.Token `json:"token"`
	}
)

func (p *GoogleDriveParameters) Validate() error {
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
	//TODO implement me
	panic("implement me")
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

func (c *GoogleDrive) scanFolder(ctx context.Context, folderID string) error {
	var q string
	if folderID == "" {
		q = googleDriveRootFolderQuery
	} else {
		q = fmt.Sprintf(" and '%s' in parents ", folderID)
	}
	if err := c.client.Files.List().Context(ctx).Q(q).Pages(ctx, func(l *drive.FileList) error {
		for _, item := range l.Items {

		}

	}); err != nil {
		return err
	}
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

		folder, err := c.client.Files.List().Context(ctx).Fields("files(id,name)").Q(q).Do()
		if err != nil {
			return "", err
		}
		if len(folder.Items) == 0 {
			return "", fmt.Errorf("folder %s not found", strings.Join(folderParts[:i], "/"))
		}
		parentID = folder.Items[0].Id
	}
	if parentID == "" {
		return "", fmt.Errorf("folder %s not found", c.param.Folder)
	}
	return parentID, nil

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
		param: &GoogleDriveParameters{},
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
		log.Fatalf("Unable to retrieve driveactivity Client %v", err)
	}

	conn.client = client
	return &conn, nil
}

func GetDriver(token *oauth2.Token) {
	ctx := context.Background()

	// If modifying these scopes, delete your previously saved token.json.
	srv, err := drive.NewService(ctx,
		option.WithHTTPClient(&http.Client{Transport: utils.NewTransport(token)}))
	if err != nil {
		log.Fatalf("Unable to retrieve driveactivity Client %v", err)
	}
	//docSrv, err := docs.NewService(ctx, option.WithHTTPClient(&http.Client{Transport: utils.NewTransport(token)}))

	//Fields("nextPageToken, files(id, name)").
	r, err := srv.Drives.List().Do()
	if err != nil {
		log.Fatalf("Unable to retrieve list of activities. %v", err)
	}
	fmt.Println("Recent Activity:")
	for _, dr := range r.Items {
		fmt.Printf(" id %s name %s\n", dr.Id, dr.Name)
	}
	fr, err := srv.Files.List().Do()
	if err != nil {
		log.Fatalf("Unable to retrieve list of activities. %v", err)
	}

	for _, f := range fr.Items {
		fullpath := []string{}
		for _, p := range f.Parents {
			fullpath = append(fullpath, p.Id)
		}

		fmt.Printf("folder %s (%s)\n parent %s  \n", f.Title, f.Id, strings.Join(fullpath, ","))
		// ${dataSource.folder_id}' in parents
		nfr, err := srv.Files.List().Q(fmt.Sprintf("'%s' in parents", f.Id)).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve list of activities. %v", err)
		}
		for _, d := range nfr.Items {
			fmt.Printf("\t---\t\t%s- %s \n", d.Id, d.Title)
			//resp, err := srv.Files.Get(d.Id).Download()
			//if err != nil {
			//	log.Fatalf("Unable to retrieve list of activities. %v", err)
			//}

		}
		fmt.Printf("\t %s- %s \n", f.Id, f.Title)
	}
}
