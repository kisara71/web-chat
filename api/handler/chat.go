package handler

import (
	"encoding/json"
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
	var req httpmodel.Response
	if err := c.ShouldBindJSON(&req); err != nil {
		Respond(c, http.StatusBadRequest, errcode.CodeBadRequest, "invalid request", nil)
		return
	}
	stream, err := h.logic.ResponseStream(c.Request.Context(), &req)
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
		emitStreamEvent(c, ev)
		flusher.Flush()
		if done {
			break
		}
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
