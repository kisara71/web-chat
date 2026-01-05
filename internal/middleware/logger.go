package middleware

import (
	"web-chat/pkg/logger"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	lgr := logger.L()
	return gin.LoggerWithWriter(lgr.Writer())
}
