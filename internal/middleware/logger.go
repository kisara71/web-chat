package middleware

import (
	"web-chat/pkg/logger"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return gin.LoggerWithWriter(logger.Writer())
}
