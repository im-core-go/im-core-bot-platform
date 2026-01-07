package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	httpmodel "web-chat/api/http_model/chat"
	logicchat "web-chat/internal/logic/chat"
	"web-chat/internal/logic/chat/impls/openai"
	"web-chat/internal/middleware"
	"web-chat/internal/svc"
	errcode "web-chat/pkg/err"
	"web-chat/pkg/logger"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	svcCtx *svc.Context
	logic  logicchat.Logic
}

func NewChatHandler(svcCtx *svc.Context) (*ChatHandler, error) {
	logic, err := openai.NewChatLogic(svcCtx)
	if err != nil {
		return nil, err
	}
	return &ChatHandler{svcCtx: svcCtx, logic: logic}, nil
}

func (h *ChatHandler) RegisterRoutes(engine *gin.Engine) {
	chatGroup := engine.Group("/api/chat")
	chatGroup.Use(middleware.Auth(h.svcCtx))
	chatGroup.GET("/models", h.Models)
	chatGroup.POST("/conversation", h.CreateConversation)
	chatGroup.GET("/conversations", h.ListConversations)
	chatGroup.GET("/conversation", h.GetConversation)
	chatGroup.PATCH("/conversation/title", h.UpdateConversationTitle)
	chatGroup.DELETE("/conversation", h.DeleteConversation)
	chatGroup.GET("/messages", h.ListMessages)
	chatGroup.DELETE("/messages", h.ClearMessages)
	chatGroup.POST("/stream", h.Stream)
}

func (h *ChatHandler) Models(c *gin.Context) {
	resp, err := h.logic.PullModules(c.Request.Context())
	if err != nil {
		logger.L().Errorf("pull models error: %v", err)
		Respond(c, http.StatusBadRequest, errcode.CodeBadRequest, err.Error(), nil)
		return
	}
	OK(c, resp)
}

func (h *ChatHandler) Stream(c *gin.Context) {
	var req httpmodel.Completion
	if err := c.ShouldBindJSON(&req); err != nil {
		Respond(c, http.StatusBadRequest, errcode.CodeBadRequest, "invalid request", nil)
		return
	}
	userID, ok := c.Get("user_id")
	if !ok {
		Respond(c, http.StatusUnauthorized, errcode.CodeUnauthorized, "unauthorized", nil)
		return
	}
	stream, conversationID, err := h.logic.ResponseStream(c.Request.Context(), &req, userIDToString(userID))
	if err != nil {
		logger.L().Errorf("chat stream error: %v", err)
		Respond(c, http.StatusBadRequest, errcode.CodeBadRequest, err.Error(), nil)
		return
	}
	defer stream.Close()

	c.Set("err_code", int(errcode.CodeOK))
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Status(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		Respond(c, http.StatusInternalServerError, errcode.CodeInternal, "stream unsupported", nil)
		return
	}

	for {
		ev, done, err := stream.Next()
		if err != nil {
			logger.L().Errorf("chat stream next error: %v", err)
			c.Set("err_code", int(errcode.CodeInternal))
			emitStreamEvent(c, httpmodel.StreamEvent{Type: httpmodel.EventError, Delta: "internal error"})
			break
		}
		if done && ev.ConversationID == "" {
			ev.ConversationID = conversationID
		}
		emitStreamEvent(c, ev)
		flusher.Flush()
		if done {
			break
		}
	}
}

func (h *ChatHandler) CreateConversation(c *gin.Context) {
	var req httpmodel.CreateConversationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Respond(c, http.StatusBadRequest, errcode.CodeBadRequest, "invalid request", nil)
		return
	}
	userID, ok := c.Get("user_id")
	if !ok {
		Respond(c, http.StatusUnauthorized, errcode.CodeUnauthorized, "unauthorized", nil)
		return
	}
	resp, err := h.logic.CreateConversation(c.Request.Context(), &req, userIDToString(userID))
	if err != nil {
		logger.L().Errorf("create conversation error: %v", err)
		Respond(c, http.StatusBadRequest, errcode.CodeBadRequest, err.Error(), nil)
		return
	}
	OK(c, resp)
}

func (h *ChatHandler) ListConversations(c *gin.Context) {
	var req httpmodel.ListConversationsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Respond(c, http.StatusBadRequest, errcode.CodeBadRequest, "invalid request", nil)
		return
	}
	userID, ok := c.Get("user_id")
	if !ok {
		Respond(c, http.StatusUnauthorized, errcode.CodeUnauthorized, "unauthorized", nil)
		return
	}
	resp, err := h.logic.ListConversations(c.Request.Context(), &req, userIDToString(userID))
	if err != nil {
		logger.L().Errorf("list conversations error: %v", err)
		Respond(c, http.StatusBadRequest, errcode.CodeBadRequest, err.Error(), nil)
		return
	}
	OK(c, resp)
}

func (h *ChatHandler) ListMessages(c *gin.Context) {
	var req httpmodel.ListMessagesReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Respond(c, http.StatusBadRequest, errcode.CodeBadRequest, "invalid request", nil)
		return
	}
	userID, ok := c.Get("user_id")
	if !ok {
		Respond(c, http.StatusUnauthorized, errcode.CodeUnauthorized, "unauthorized", nil)
		return
	}
	resp, err := h.logic.ListMessages(c.Request.Context(), &req, userIDToString(userID))
	if err != nil {
		logger.L().Errorf("list messages error: %v", err)
		Respond(c, http.StatusBadRequest, errcode.CodeBadRequest, err.Error(), nil)
		return
	}
	OK(c, resp)
}

func (h *ChatHandler) GetConversation(c *gin.Context) {
	var req httpmodel.GetConversationReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Respond(c, http.StatusBadRequest, errcode.CodeBadRequest, "invalid request", nil)
		return
	}
	userID, ok := c.Get("user_id")
	if !ok {
		Respond(c, http.StatusUnauthorized, errcode.CodeUnauthorized, "unauthorized", nil)
		return
	}
	resp, err := h.logic.GetConversation(c.Request.Context(), &req, userIDToString(userID))
	if err != nil {
		logger.L().Errorf("get conversation error: %v", err)
		Respond(c, http.StatusBadRequest, errcode.CodeBadRequest, err.Error(), nil)
		return
	}
	OK(c, resp)
}

func (h *ChatHandler) UpdateConversationTitle(c *gin.Context) {
	var req httpmodel.UpdateConversationTitleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Respond(c, http.StatusBadRequest, errcode.CodeBadRequest, "invalid request", nil)
		return
	}
	userID, ok := c.Get("user_id")
	if !ok {
		Respond(c, http.StatusUnauthorized, errcode.CodeUnauthorized, "unauthorized", nil)
		return
	}
	if err := h.logic.UpdateConversationTitle(c.Request.Context(), &req, userIDToString(userID)); err != nil {
		logger.L().Errorf("update conversation title error: %v", err)
		Respond(c, http.StatusBadRequest, errcode.CodeBadRequest, err.Error(), nil)
		return
	}
	OK(c, nil)
}

func (h *ChatHandler) DeleteConversation(c *gin.Context) {
	var req httpmodel.DeleteConversationReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Respond(c, http.StatusBadRequest, errcode.CodeBadRequest, "invalid request", nil)
		return
	}
	userID, ok := c.Get("user_id")
	if !ok {
		Respond(c, http.StatusUnauthorized, errcode.CodeUnauthorized, "unauthorized", nil)
		return
	}
	if err := h.logic.DeleteConversation(c.Request.Context(), &req, userIDToString(userID)); err != nil {
		logger.L().Errorf("delete conversation error: %v", err)
		Respond(c, http.StatusBadRequest, errcode.CodeBadRequest, err.Error(), nil)
		return
	}
	OK(c, nil)
}

func (h *ChatHandler) ClearMessages(c *gin.Context) {
	var req httpmodel.ClearMessagesReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Respond(c, http.StatusBadRequest, errcode.CodeBadRequest, "invalid request", nil)
		return
	}
	userID, ok := c.Get("user_id")
	if !ok {
		Respond(c, http.StatusUnauthorized, errcode.CodeUnauthorized, "unauthorized", nil)
		return
	}
	if err := h.logic.ClearMessages(c.Request.Context(), &req, userIDToString(userID)); err != nil {
		logger.L().Errorf("clear messages error: %v", err)
		Respond(c, http.StatusBadRequest, errcode.CodeBadRequest, err.Error(), nil)
		return
	}
	OK(c, nil)
}

func userIDToString(val any) string {
	switch v := val.(type) {
	case string:
		return v
	case int64:
		return fmt.Sprintf("%d", v)
	case int:
		return fmt.Sprintf("%d", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func emitStreamEvent(c *gin.Context, ev httpmodel.StreamEvent) {
	payload, err := json.Marshal(ev)
	if err != nil {
		logger.L().Errorf("chat stream marshal error: %v", err)
		return
	}
	_, _ = c.Writer.Write([]byte("data: "))
	_, _ = c.Writer.Write(payload)
	_, _ = c.Writer.Write([]byte("\n\n"))
}
