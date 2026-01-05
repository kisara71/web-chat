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

	userGroup := engine.Group("/api/user")
	userGroup.Use(middleware.Auth(h.svcCtx))
	userGroup.PUT("/update", h.Update)
	userGroup.POST("/logout", h.Logout)
}

func (h *UserHandler) Register(c *gin.Context) {
	var req http_model.UserRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, 400, errcode.CodeBadRequest, "invalid request", nil)
		return
	}
	if err := h.logic.Register(&req); err != nil {
		logger.L().Printf("register error: %v", err)
		respond(c, 400, errcode.CodeBadRequest, err.Error(), nil)
		return
	}
	ok(c, nil)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req http_model.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, 400, errcode.CodeBadRequest, "invalid request", nil)
		return
	}
	token, err := h.logic.Login(&req)
	if err != nil {
		logger.L().Printf("login error: %v", err)
		respond(c, 401, errcode.CodeUnauthorized, err.Error(), nil)
		return
	}
	ok(c, gin.H{"token": token})
}

func (h *UserHandler) Update(c *gin.Context) {
	var req http_model.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, 400, errcode.CodeBadRequest, "invalid request", nil)
		return
	}
	userID, ok := c.Get("user_id")
	if !ok {
		respond(c, 401, errcode.CodeUnauthorized, "unauthorized", nil)
		return
	}
	id, ok := userID.(int64)
	if !ok {
		respond(c, 401, errcode.CodeUnauthorized, "invalid user", nil)
		return
	}
	if err := h.logic.Update(id, &req); err != nil {
		logger.L().Printf("update error: %v", err)
		respond(c, 400, errcode.CodeBadRequest, err.Error(), nil)
		return
	}
	ok(c, nil)
}

func (h *UserHandler) Logout(c *gin.Context) {
	var req http_model.LogoutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, 400, errcode.CodeBadRequest, "invalid request", nil)
		return
	}
	if err := h.logic.Logout(&req); err != nil {
		logger.L().Printf("logout error: %v", err)
		respond(c, 400, errcode.CodeBadRequest, err.Error(), nil)
		return
	}
	ok(c, nil)
}
