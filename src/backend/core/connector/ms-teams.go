package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
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
	msTeamsMessagesURL = "https://graph.microsoft.com/v1.0/teams/%s/channels/%s/messages/microsoft.graph.delta()"
	msTeamRepliesURL   = "https://graph.microsoft.com/v1.0/teams/%s/channels/%s/messages/%s/replies"
	msTeamsInfoURL     = "https://graph.microsoft.com/v1.0/teams"

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
soups

TeamMember.Read.All,
TeamMember.ReadWrite.All

Chat.ReadBasic, Chat.Read, Chat.ReadWrite
AuditLogsQuery-OneDrive.Read.All, Files.Read.All

ServiceActivity-OneDrive.Read.All, TeamsApp.Read.All, TeamsApp.ReadWrite.All, User.Read, email, TeamMember.ReadWrite.All'",

	"id": "6fd101cf-ddca-4bef-9fdc-f7fd024c7063",
	"id": "94100e5f-a30f-433d-965e-bde4e817f62a",

94100e5f-a30f-433d-965e-bde4e817f62a
19:65a0a68789ea4abe97c8eec4d6f43786@thread.tacv2

https://graph.microsoft.com/v1.0/groups/94100e5f-a30f-433d-965e-bde4e817f62a/team/channels/19:65a0a68789ea4abe97c8eec4d6f43786@thread.tacv2
https://graph.microsoft.com/v1.0/teams/94100e5f-a30f-433d-965e-bde4e817f62a/channels
https://graph.microsoft.com/v1.0/teams/94100e5f-a30f-433d-965e-bde4e817f62a/channels/19:65a0a68789ea4abe97c8eec4d6f43786@thread.tacv2/messages
https://graph.microsoft.com/v1.0/teams/94100e5f-a30f-433d-965e-bde4e817f62a/channels/19:65a0a68789ea4abe97c8eec4d6f43786@thread.tacv2/messages/
https://graph.microsoft.com/v1.0/teams/94100e5f-a30f-433d-965e-bde4e817f62a/channels/19:65a0a68789ea4abe97c8eec4d6f43786@thread.tacv2/messages/1718016334912/replies

https://graph.microsoft.com/v1.0/teams/94100e5f-a30f-433d-965e-bde4e817f62a/channels/19:65a0a68789ea4abe97c8eec4d6f43786@thread.tacv2/messages/1718121958378/replies

	 -- delete topic


		"id": "19:65a0a68789ea4abe97c8eec4d6f43786@thread.tacv2",
		      "createdDateTime": "2024-06-10T10:45:10.413Z",
		      "displayName": "developmanet",

chat 19:09c30123-8d63-4fca-909a-3af0d3f03a4a_5d51d22a-6b76-4177-928b-28e15caf71cd@unq.gbl.spaces

Team.ReadBasic.All,
TeamSettings.Read.All,
TeamSettings.ReadWrite.All'.
ChannelMessage.Read.All

	APIConnectors.Read.All,

APIConnectors.ReadWrite.All,
AuditLogsQuery-OneDrive.Read.All,
Chat.Read,
Chat.ReadBasic,
Files.Read.All, openid, profile,
ServiceActivity-OneDrive.Read.All,
TeamMember.Read.All, TeamMember.ReadWrite.All, TeamsApp.Read.All, TeamsApp.ReadWrite.All, User.Read, email, Group.Read.All
*/
type (
	MSTeams struct {
		Base
		param         *MSTeamParameters
		state         *MSTeamState
		client        *resty.Client
		fileSizeLimit int
		sessionID     uuid.NullUUID
		chResult      chan *Response
	}
	MSTeamParameters struct {
		Channel string            `json:"channel"`
		Topics  model.StringSlice `json:"topics"`
		Chat    string            `json:"chat"`
		Token   oauth2.Token      `json:"token"`
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
		From    string
		Message string
	}
)

func (c *MSTeams) Validate() error {
	return nil
}

func (c *MSTeams) Execute(ctx context.Context, param map[string]string) chan *Response {
	c.resultCh = make(chan *Response)
	defer close(c.resultCh)
	if err := c.execute(ctx, param); err != nil {
		zap.S().Errorf(err.Error())
	}
	return c.chResult
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
	messages, err := c.getMessagesByChannel(ctx, teamID, channelID)

}

func (c *MSTeams) PrepareTask(ctx context.Context, task Task) error {
	params := make(map[string]string)

	teamID, err := c.getTeamID(ctx)
	if err != nil {
		zap.S().Errorf(err.Error())
	}
	params[msTeamsParamTeamID] = teamID

	return task.RunConnector(ctx, &proto.ConnectorRequest{
		Id:     c.model.ID.IntPart(),
		Params: params,
	})
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

func (c *MSTeams) getReplies(ctx context.Context, teamID, channelID string, msg *MessageBody) ([]*MSTeamsResult, error) {
	var repliesResp MessageResponse
	err := c.requestAndParse(ctx, fmt.Sprintf(msTeamRepliesURL, teamID, channelID, msg.Id), &repliesResp)
	if err != nil {
		return nil, err
	}
	state, ok := c.state.Topics[msg.Id]
	if !ok {
		state = &MSTeamMessageState{}
	}
	lastTime := state.LastCreatedDateTime
	var results []*MSTeamsResult
	for _, repl := range repliesResp.Value {
		if repl.CreatedDateTime.Before(state.LastCreatedDateTime) {
			// ignore messages that were analyzed before
			continue
		}
		if lastTime.Before(repl.CreatedDateTime) {
			// store timestamp of last message
			lastTime = repl.CreatedDateTime
		}

		message := repl.Body.Content
		if repl.Body.ContentType == "html" {
			message, err = html2text.FromString(message, html2text.Options{
				PrettyTables: true,
			})
		}
		results = append(results, &MSTeamsResult{
			From:    repl.From.User.DisplayName,
			Message: message,
		})
	}
	return results, nil
}
func (c *MSTeams) getMessagesByChannel(ctx context.Context, teamID, channelID string) ([]*MessageBody, error) {
	var messagesResp MessageResponse

	if err := c.requestAndParse(ctx, fmt.Sprintf(msTeamsMessagesURL, teamID, channelID), &messagesResp); err != nil {
		return nil, err
	}
	// todo store url for incremental request
	//_ = channelResp.OdataDeltaLink

	// todo add validation on Subject == null - topic was deleted.
	messagesForScan := make([]*MessageBody, 0)
	for _, msg := range messagesResp.Value {
		if c.param.Topics.InArray(msg.Subject) {
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

func (c *MSTeams) requestAndParse(ctx context.Context, url string, result interface{}) error {
	response, err := c.client.R().SetContext(ctx).Get(url)
	if err = utils.WrapleRestyError(response, err); err != nil {
		return err
	}
	return json.Unmarshal(response.Body(), result)
}

// NewMSTeams creates new instance of MsTeams connector
func NewMSTeams(connector *model.Connector) (Connector, error) {
	conn := MSTeams{}
	conn.Base.Config(connector)
	conn.param = &MSTeamParameters{}
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
