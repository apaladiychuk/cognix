package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"time"
)

const (
	msTeamsChannelsURL = "https://graph.microsoft.com/v1.0/teams/%s/channels"
	msTeamsMessagesURL = "https://graph.microsoft.com/v1.0/teams/%s/channels/%s/messages"
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
	CreatedDateTime      time.Time     `json:"createdDateTime"`
	LastModifiedDateTime time.Time     `json:"lastModifiedDateTime"`
	From                 *TeamFrom     `json:"from"`
	Body                 *TeamBody     `json:"body"`
	Attachments          []*Attachment `json:"attachments"`
}
type MessageResponse struct {
	Value []*MessageBody `json:"value"`
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
		client        *resty.Client
		fileSizeLimit int
		sessionID     uuid.NullUUID
		chResult      chan *Response
	}
	MSTeamParameters struct {
		Channel string       `json:"channel"`
		Topics  []string     `json:"topics"`
		Chat    string       `json:"chat"`
		Token   oauth2.Token `json:"token"`
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

	teamID, err := c.getGroup(ctx)
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

}

func (c *MSTeams) getMessagesByChannel(ctx context.Context, teamID, channelID string) (string, error) {
	var channels []*ChannelResponse
	response, err := c.client.R().SetContext(ctx).Get(fmt.Sprintf(msTeamsChannelsURL, teamID))

}

func (c *MSTeams) getGroup(ctx context.Context) (string, error) {
	var team TeamResponse
	response, err := c.client.R().
		SetContext(ctx).
		Get(msTeamsInfoURL)
	if err != nil || response.IsError() {
		return "", err
	}
	if err = json.Unmarshal(response.Body(), &team); err != nil {
		return "", err
	}
	if len(team.Value) == 0 {
		return "", fmt.Errorf("team not found")
	}
	return team.Value[0].Id, nil
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
