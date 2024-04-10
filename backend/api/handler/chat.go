package handler

import (
	"cognix.ch/api/v2/core/bll"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/responder"
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/server"
	"cognix.ch/api/v2/core/utils"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

type ChatHandler struct {
	chatBL bll.ChatBL
}

func NewChatHandler(chatBL bll.ChatBL) *ChatHandler {
	return &ChatHandler{chatBL: chatBL}

}

func (h *ChatHandler) Mount(route *gin.Engine, authMiddleware gin.HandlerFunc) {
	handler := route.Group("/chats")
	handler.Use(authMiddleware)
	handler.GET("/get-user-chat-sessions", server.HandlerErrorFuncAuth(h.GetSessions))
	handler.GET("/get-chat-session/:id", server.HandlerErrorFuncAuth(h.GetByID))
	handler.POST("/create-chat-session", server.HandlerErrorFuncAuth(h.CreateSession))
	handler.POST("/send-message", server.HandlerErrorFuncAuth(h.SendMessage))
}

// GetSessions return list of chat sessions for current user
// @Summary return list of chat sessions for current user
// @Description return list of chat sessions for current user
// @Tags Chat
// @ID chat_get_sessions
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} model.ChatSession
// @Router /chats/get-user-chat-sessions [get]
func (h *ChatHandler) GetSessions(c *gin.Context, identity *security.Identity) error {
	sessions, err := h.chatBL.GetSessions(c.Request.Context(), identity.User)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, sessions)
}

// GetByID return chat session with messages by given id
// @Summary return chat session with messages by given id
// @Description return chat session with messages by given id
// @Tags Chat
// @ID chat_get_by_id
// @Param id path int true "session id"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.ChatSession
// @Router /chats/get-chat-session/{id} [get]
func (h *ChatHandler) GetByID(c *gin.Context, identity *security.Identity) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.InvalidInput.New("id should be presented")
	}
	session, err := h.chatBL.GetSessionByID(c.Request.Context(), identity.User, id)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, session)
}

// CreateSession creates new chat session
// @Summary creates new chat session
// @Description creates new chat session
// @Tags Chat
// @ID chat_create_session
// @Param payload body parameters.CreateChatSession true "create session parameters"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.ChatSession
// @Router /chats/create-chat-session [post]
func (h *ChatHandler) CreateSession(c *gin.Context, identity *security.Identity) error {
	var param parameters.CreateChatSession
	if err := c.BindJSON(&param); err != nil {
		return utils.InvalidInput.Wrap(err, "wrong payload")
	}
	if err := param.Validate(); err != nil {
		return utils.InvalidInput.Wrap(err, err.Error())
	}
	session, err := h.chatBL.CreateSession(c.Request.Context(), identity.User, &param)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusCreated, session)
}

// SendMessage send message and wait stream response
// @Summary send message and wait stream response
// @Description send message and wait stream response
// @Tags Chat
// @ID chat_send_message
// @Param payload body parameters.CreateChatMessageRequest true "send message parameters"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.ChatMessage
// @Router /chats/create-chat-session [post]
func (h *ChatHandler) SendMessage(c *gin.Context, identity *security.Identity) error {
	var param parameters.CreateChatMessageRequest
	if err := c.BindJSON(&param); err != nil {
		return utils.InvalidInput.Wrap(err, "wrong payload")
	}
	assistant, err := h.chatBL.SendMessage(c, identity.User, &param)
	if err != nil {
		return err
	}
	c.Stream(func(w io.Writer) bool {
		response, ok := assistant.Receive()
		c.SSEvent(responder.ResponseMessage, response.Message)
		return ok
	})
	return nil
}
