package api

import (
	"time"
	"web-chat/api/handler"
	"web-chat/internal/middleware"
	"web-chat/internal/svc"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(svcCtx *svc.Context) (*gin.Engine, error) {
	engine := gin.New()
	engine.Use(
		middleware.Logger(),
		middleware.Recovery(),
		middleware.CORS(),
		middleware.Metrics(),
		middleware.RateLimit(svcCtx, 120, time.Second*10),
	)
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	handler.NewUserHandler(svcCtx).RegisterRoutes(engine)
	chatHandler, err := handler.NewChatHandler(svcCtx)
	if err != nil {
		return nil, err
	}
	chatHandler.RegisterRoutes(engine)

	return engine, nil
}
