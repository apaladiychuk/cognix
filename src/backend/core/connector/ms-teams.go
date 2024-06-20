package connector

import (
	microsoft_core "cognix.ch/api/v2/core/connector/microsoft-core"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"jaytaylor.com/html2text"
	"time"
)

const (
	msTeamsChannelsURL = "https://graph.microsoft.com/v1.0/teams/%s/channels"
	//msTeamsMessagesURL = "https://graph.microsoft.com/v1.0/teams/%s/channels/%s/messages/microsoft.graph.delta()"
	msTeamsMessagesURL = "https://graph.microsoft.com/v1.0/teams/%s/channels/%s/messages"
	msTeamRepliesURL   = "https://graph.microsoft.com/v1.0/teams/%s/channels/%s/messages/%s/replies"
	msTeamsInfoURL     = "https://graph.microsoft.com/v1.0/teams"

	msTeamsFilesFolder   = "https://graph.microsoft.com/v1.0/teams/%s/channels/%s/filesFolder"
	msTeamsFolderContent = "https://graph.microsoft.com/v1.0/groups/%s/drive/items/%s/children"

	msTeamsChats           = "https://graph.microsoft.com/v1.0/chats"
	msTeamsChatMessagesURL = "https://graph.microsoft.com/v1.0/chats/%s/messages"

	msTeamsParamTeamID    = "team_id"
	msTeamsParamChannelID = "channel_id"
)

type (
	Team struct {
		Id          string `json:"id"`
		DisplayName string `json:"displayName"`
		Description string `json:"description"`
	}

	TeamResponse struct {
		Value []*Team `json:"value"`
	}
	ChannelResponse struct {
		Value []*ChannelBody `json:"value"`
	}
	ChannelBody struct {
		Id              string    `json:"id"`
		CreatedDateTime time.Time `json:"createdDateTime"`
		DisplayName     string    `json:"displayName"`
		Description     string    `json:"description"`
	}
	TeamUser struct {
		OdataType        string `json:"@odata.type"`
		Id               string `json:"id"`
		DisplayName      string `json:"displayName"`
		UserIdentityType string `json:"userIdentityType"`
		TenantId         string `json:"tenantId"`
	}
	TeamFrom struct {
		User *TeamUser `json:"user"`
	}

	TeamBody struct {
		ContentType string `json:"contentType"`
		Content     string `json:"content"`
	}
)

type TeamFilesFolder struct {
	Id string `json:"id"`
}

type MessageBody struct {
	Id                   string        `json:"id"`
	Etag                 string        `json:"etag"`
	MessageType          string        `json:"messageType"`
	ReplyToId            string        `json:"replyToId"`
	Subject              string        `json:"subject"`
	CreatedDateTime      time.Time     `json:"createdDateTime"`
	LastModifiedDateTime time.Time     `json:"lastModifiedDateTime"`
	DeletedDateTime      pg.NullTime   `json:"deletedDateTime"`
	From                 *TeamFrom     `json:"from"`
	Body                 *TeamBody     `json:"body"`
	Attachments          []*Attachment `json:"attachments"`
}
type MessageResponse struct {
	OdataContext   string         `json:"@odata.context"`
	OdataNextLink  string         `json:"@odata.nextLink"`
	OdataDeltaLink string         `json:"@odata.deltaLink"`
	Value          []*MessageBody `json:"value"`
}
type Attachment struct {
	Id           string      `json:"id"`
	ContentType  string      `json:"contentType"`
	ContentUrl   string      `json:"contentUrl"`
	Content      interface{} `json:"content"`
	Name         string      `json:"name"`
	ThumbnailUrl interface{} `json:"thumbnailUrl"`
	TeamsAppId   interface{} `json:"teamsAppId"`
}

/*
https://graph.microsoft.com/v1.0/groups/94100e5f-a30f-433d-965e-bde4e817f62a/team/channels/19:65a0a68789ea4abe97c8eec4d6f43786@thread.tacv2
https://graph.microsoft.com/v1.0/teams/94100e5f-a30f-433d-965e-bde4e817f62a/channels
https://graph.microsoft.com/v1.0/teams/94100e5f-a30f-433d-965e-bde4e817f62a/channels/19:65a0a68789ea4abe97c8eec4d6f43786@thread.tacv2/messages
https://graph.microsoft.com/v1.0/teams/94100e5f-a30f-433d-965e-bde4e817f62a/channels/19:65a0a68789ea4abe97c8eec4d6f43786@thread.tacv2/messages/
https://graph.microsoft.com/v1.0/teams/94100e5f-a30f-433d-965e-bde4e817f62a/channels/19:65a0a68789ea4abe97c8eec4d6f43786@thread.tacv2/messages/1718016334912/replies

https://graph.microsoft.com/v1.0/teams/94100e5f-a30f-433d-965e-bde4e817f62a/channels/19:65a0a68789ea4abe97c8eec4d6f43786@thread.tacv2/messages/1718121958378/replies

https://graph.microsoft.com/v1.0/drives/01SZITRJYIBUNPFKHAYJCLISV4DXJPIJV4/items/

https://graph.microsoft.com/v1.0/drives/b!oxsuyS45_EKmyHYegUv4SmEjVp8sBIFPvH1TNMZJZqPviFyz50UFTqjI-nC6wDfJ

https://graph.microsoft.com/v1.0/groups/94100e5f-a30f-433d-965e-bde4e817f62a/drive/root/children

// get drive items for channel
https://graph.microsoft.com/v1.0/teams/94100e5f-a30f-433d-965e-bde4e817f62a/channels/19:65a0a68789ea4abe97c8eec4d6f43786@thread.tacv2/filesFolder
// get files from channel
https://graph.microsoft.com/v1.0/groups/94100e5f-a30f-433d-965e-bde4e817f62a/drive/items/01SZITRJYIBUNPFKHAYJCLISV4DXJPIJV4/children

	/groups/94100e5f-a30f-433d-965e-bde4e817f62a/drive/items/01SZITRJYIBUNPFKHAYJCLISV4DXJPIJV4
*/
type (
	MSTeams struct {
		Base
		param         *MSTeamParameters
		state         *MSTeamState
		client        *resty.Client
		fileSizeLimit int
		sessionID     uuid.NullUUID
	}
	MSTeamParameters struct {
		Channel string                       `json:"channel"`
		Topics  model.StringSlice            `json:"topics"`
		Chat    string                       `json:"chat"`
		Token   *oauth2.Token                `json:"token"`
		Files   *microsoft_core.MSDriveParam `json:"files"`
	}
	// MSTeamState store ms team state after each execute
	MSTeamState struct {
		// Link for request changes after last execution
		DeltaLink string                         `json:"delta_link"`
		Topics    map[string]*MSTeamMessageState `json:"topics"`
	}
	// MSTeamMessageState store
	MSTeamMessageState struct {
		LastCreatedDateTime time.Time `json:"last_created_date_time"`
	}
	MSTeamsResult struct {
		PrevLoadTime string
		Messages     []*MSTeamsResultMessages
	}
	MSTeamsResultMessages struct {
		User    string `json:"user"`
		Message string `json:"message"`
	}
)

func (c *MSTeams) Validate() error {
	return nil
}

func (c *MSTeams) PrepareTask(ctx context.Context, task Task) error {
	params := make(map[string]string)

	teamID, err := c.getTeamID(ctx)
	if err != nil {
		zap.S().Errorf(err.Error())
	}
	params[msTeamsParamTeamID] = teamID

	channelID, err := c.getChannel(ctx, teamID)
	if err != nil {
		zap.S().Errorf(err.Error())
	}
	params[msTeamsParamChannelID] = channelID
	return task.RunConnector(ctx, &proto.ConnectorRequest{
		Id:     c.model.ID.IntPart(),
		Params: params,
	})
}

func (c *MSTeams) Execute(ctx context.Context, param map[string]string) chan *Response {
	c.resultCh = make(chan *Response, 10)
	for _, doc := range c.model.Docs {
		if doc.Signature == "" {
			// do not delete document with chat history.
			doc.IsExists = true
		}
	}
	go func() {
		defer close(c.resultCh)
		if err := c.execute(ctx, param); err != nil {
			zap.S().Errorf(err.Error())
		}
		return
	}()
	return c.resultCh
}
func (c *MSTeams) execute(ctx context.Context, param map[string]string) error {

	teamID, ok := param[msTeamsParamTeamID]
	if !ok {
		return fmt.Errorf("team_id is not configured")
	}
	channelID, ok := param[msTeamsParamChannelID]
	if !ok {
		return fmt.Errorf("channel_id is not configured")
	}
	topics, err := c.getTopicsByChannel(ctx, teamID, channelID)
	if err != nil {
		return err
	}
	c.sessionID = uuid.NullUUID{
		UUID:  uuid.New(),
		Valid: true,
	}
	//  load topics
	for _, topic := range topics {
		// create unique id for store new messages in new document
		sourceID := fmt.Sprintf("%s-%s", topic.Id, uuid.New().String())
		replies, err := c.getReplies(ctx, teamID, channelID, topic)
		if err != nil {
			return err
		}
		body, err := json.Marshal(replies.Messages)
		if err != nil {
			return err
		}
		doc := &model.Document{
			SourceID:        sourceID,
			ConnectorID:     c.model.ID,
			URL:             "",
			ChunkingSession: c.sessionID,
			Analyzed:        false,
			CreationDate:    time.Now().UTC(),
			LastUpdate:      pg.NullTime{time.Now().UTC()},
			IsExists:        true,
		}
		c.model.DocsMap[sourceID] = doc

		fileName := topic.Id + "_" + topic.Subject
		if replies.PrevLoadTime != "" {
			fileName += "-" + replies.PrevLoadTime
		}
		fileName += ".json"
		c.resultCh <- &Response{
			URL:        doc.URL,
			Name:       fileName,
			SourceID:   doc.SourceID,
			DocumentID: doc.ID.IntPart(),
			MimeType:   "plain/text",
			FileType:   proto.FileType_TXT,
			Signature:  "",
			Content: &Content{
				Bucket:        model.BucketName(c.model.User.EmbeddingModel.TenantID),
				URL:           "",
				AppendContent: true,
				Body:          body,
			},
			UpToData: false,
		}
	}
	if c.param.Files != nil {
		if err = c.loadFiles(ctx, param, teamID, channelID); err != nil {
			return err
		}
	}
	// save current state
	if err = c.model.State.FromStruct(c.state); err == nil {
		return c.connectorRepo.Update(ctx, c.model)
	}
	return nil
}

func (c *MSTeams) loadFiles(ctx context.Context, param map[string]string, teamID, channelID string) error {
	var folderInfo TeamFilesFolder
	if err := c.requestAndParse(ctx, fmt.Sprintf(msTeamsFilesFolder, teamID, channelID), &folderInfo); err != nil {
		return err
	}
	baseUrl := fmt.Sprintf(msTeamsFolderContent, teamID, folderInfo.Id)
	folderURL := fmt.Sprintf(msTeamsFolderContent, teamID, `%s`)
	msDrive := microsoft_core.NewMSDrive(c.param.Files,
		c.model,
		c.sessionID, c.client,
		baseUrl, folderURL,
		c.getFile,
	)
	return msDrive.Execute(ctx, param)

}

func (c *MSTeams) getChannel(ctx context.Context, teamID string) (string, error) {
	var channelResp ChannelResponse
	if err := c.requestAndParse(ctx, fmt.Sprintf(msTeamsChannelsURL, teamID), &channelResp); err != nil {
		return "", err
	}
	for _, channel := range channelResp.Value {
		if channel.DisplayName == c.param.Channel {
			return channel.Id, nil
		}
	}
	return "", fmt.Errorf("channel not found")
}

func (c *MSTeams) getReplies(ctx context.Context, teamID, channelID string, msg *MessageBody) (*MSTeamsResult, error) {
	var repliesResp MessageResponse
	err := c.requestAndParse(ctx, fmt.Sprintf(msTeamRepliesURL, teamID, channelID, msg.Id), &repliesResp)
	if err != nil {
		return nil, err
	}
	var result MSTeamsResult
	state, ok := c.state.Topics[msg.Id]
	if !ok {
		state = &MSTeamMessageState{}
		c.state.Topics[msg.Id] = state
	} else {
		result.PrevLoadTime = state.LastCreatedDateTime.Format("2006-01-02-15-04-05")
	}
	lastTime := state.LastCreatedDateTime

	for _, repl := range repliesResp.Value {
		if state.LastCreatedDateTime.After(repl.CreatedDateTime) ||
			state.LastCreatedDateTime.Equal(repl.CreatedDateTime) {
			// ignore messages that were analyzed before
			continue
		}
		if repl.CreatedDateTime.After(lastTime) {
			// store timestamp of last message
			lastTime = repl.CreatedDateTime
		}

		message := repl.Body.Content
		if repl.Body.ContentType == "html" {
			message, err = html2text.FromString(message, html2text.Options{
				PrettyTables: true,
			})
		}
		result.Messages = append(result.Messages, &MSTeamsResultMessages{
			User:    repl.From.User.DisplayName,
			Message: message,
		})
	}
	state.LastCreatedDateTime = lastTime
	return &result, nil
}

func (c *MSTeams) getTopicsByChannel(ctx context.Context, teamID, channelID string) ([]*MessageBody, error) {
	var messagesResp MessageResponse
	// Get url from state. Load changes from previous scan.

	//url := c.state.DeltaLink
	//if url == "" {
	//	// Load all history if stored lin is empty
	//	url = fmt.Sprintf(msTeamsMessagesURL, teamID, channelID)incremental request
	//}
	url := fmt.Sprintf(msTeamsMessagesURL, teamID, channelID)
	if err := c.requestAndParse(ctx, url, &messagesResp); err != nil {
		return nil, err
	}
	//if messagesResp.OdataNextLink != "" {
	//	c.state.DeltaLink = messagesResp.OdataNextLink
	//}
	//if messagesResp.OdataDeltaLink != "" {
	//	c.state.DeltaLink = messagesResp.OdataDeltaLink
	//}

	messagesForScan := make([]*MessageBody, 0)
	for _, msg := range messagesResp.Value {
		if msg.Subject == "" {
			// todo add validation on Subject == null - topic was deleted.
			//for _, doc := range c.model.DocsMap {
			//	if strings.HasPrefix(doc.SourceID, msg.Id) {
			//		doc.IsExists = false
			//	}
			//}
			//delete(c.state.Topics, msg.Id)
			continue
		}
		if len(c.param.Topics) == 0 || c.param.Topics.InArray(msg.Subject) {
			messagesForScan = append(messagesForScan, msg)
		}
	}
	return messagesForScan, nil
}

// getTeamID get team id for current user
func (c *MSTeams) getTeamID(ctx context.Context) (string, error) {
	var team TeamResponse

	if err := c.requestAndParse(ctx, msTeamsInfoURL, &team); err != nil {
		return "", err
	}
	if len(team.Value) == 0 {
		return "", fmt.Errorf("team not found")
	}
	return team.Value[0].Id, nil
}

// requestAndParse request graph endpoint and parse result.
func (c *MSTeams) requestAndParse(ctx context.Context, url string, result interface{}) error {
	response, err := c.client.R().SetContext(ctx).Get(url)
	if err = utils.WrapRestyError(response, err); err != nil {
		return err
	}
	return json.Unmarshal(response.Body(), result)
}

// getFile callback for receive files
func (c *MSTeams) getFile(payload *microsoft_core.Response) {
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

// NewMSTeams creates new instance of MsTeams connector
func NewMSTeams(connector *model.Connector,
	connectorRepo repository.ConnectorRepository,
	oauthURL string) (Connector, error) {
	conn := MSTeams{
		Base: Base{
			connectorRepo: connectorRepo,
			oauthClient: resty.New().
				SetTimeout(time.Minute).
				SetBaseURL(oauthURL),
		},
		param: &MSTeamParameters{},
		state: &MSTeamState{},
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
	if err = connector.State.ToStruct(conn.state); err != nil {
		zap.S().Infof("can not parse state %v", err)
	}
	if conn.state.Topics == nil {
		conn.state.Topics = make(map[string]*MSTeamMessageState)
	}

	conn.client = resty.New().
		SetTimeout(time.Minute).
		SetHeader(authorizationHeader, fmt.Sprintf("%s %s",
			conn.param.Token.TokenType,
			conn.param.Token.AccessToken))
	return &conn, nil
}
