package handler

import (
	"web-chat/api/http_model"
	"web-chat/internal/logic/user"
	"web-chat/internal/middleware"
	"web-chat/internal/svc"
	errcode "web-chat/pkg/err"
	"web-chat/pkg/logger"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	logic  user.Logic
	svcCtx *svc.Context
}

func NewUserHandler(svcCtx *svc.Context) *UserHandler {
	return &UserHandler{
		logic:  user.NewLogic(svcCtx),
		svcCtx: svcCtx,
	}
}

func (h *UserHandler) RegisterRoutes(engine *gin.Engine) {
	engine.POST("/api/user/register", h.Register)
	engine.POST("/api/user/login", h.Login)
	engine.POST("/api/user/login/code", h.LoginByCode)
	engine.POST("/api/user/email/code", h.SendEmailCode)

	userGroup := engine.Group("/api/user")
	userGroup.Use(middleware.Auth(h.svcCtx))
	userGroup.PUT("/update", h.Update)
	userGroup.POST("/logout", h.Logout)
}

func (h *UserHandler) Register(c *gin.Context) {
	var req http_model.UserRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Respond(c, 400, errcode.CodeBadRequest, "invalid request", nil)
		return
	}
	if err := h.logic.Register(&req); err != nil {
		logger.L().Errorf("register error: %v", err)
		Respond(c, 400, errcode.CodeBadRequest, err.Error(), nil)
		return
	}
	OK(c, nil)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req http_model.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Respond(c, 400, errcode.CodeBadRequest, "invalid request", nil)
		return
	}
	token, err := h.logic.Login(&req)
	if err != nil {
		logger.L().Errorf("login error: %v", err)
		Respond(c, 401, errcode.CodeUnauthorized, err.Error(), nil)
		return
	}
	OK(c, gin.H{"token": token})
}

func (h *UserHandler) LoginByCode(c *gin.Context) {
	var req http_model.LoginCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Respond(c, 400, errcode.CodeBadRequest, "invalid request", nil)
		return
	}
	token, err := h.logic.LoginByCode(&req)
	if err != nil {
		logger.L().Errorf("login by code error: %v", err)
		Respond(c, 401, errcode.CodeUnauthorized, err.Error(), nil)
		return
	}
	OK(c, gin.H{"token": token})
}

func (h *UserHandler) SendEmailCode(c *gin.Context) {
	var req http_model.SendEmailCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Respond(c, 400, errcode.CodeBadRequest, "invalid request", nil)
		return
	}
	if err := h.logic.SendEmailCode(&req); err != nil {
		logger.L().Errorf("send email code error: %v", err)
		Respond(c, 400, errcode.CodeBadRequest, err.Error(), nil)
		return
	}
	OK(c, nil)
}

func (h *UserHandler) Update(c *gin.Context) {
	var req http_model.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Respond(c, 400, errcode.CodeBadRequest, "invalid request", nil)
		return
	}
	userID, ok := c.Get("user_id")
	if !ok {
		Respond(c, 401, errcode.CodeUnauthorized, "unauthorized", nil)
		return
	}
	id, ok := userID.(int64)
	if !ok {
		Respond(c, 401, errcode.CodeUnauthorized, "invalid user", nil)
		return
	}
	if err := h.logic.Update(id, &req); err != nil {
		logger.L().Errorf("update error: %v", err)
		Respond(c, 400, errcode.CodeBadRequest, err.Error(), nil)
		return
	}
	OK(c, nil)
}

func (h *UserHandler) Logout(c *gin.Context) {
	var req http_model.LogoutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Respond(c, 400, errcode.CodeBadRequest, "invalid request", nil)
		return
	}
	if err := h.logic.Logout(&req); err != nil {
		logger.L().Errorf("logout error: %v", err)
		Respond(c, 400, errcode.CodeBadRequest, err.Error(), nil)
		return
	}
	OK(c, nil)
}
