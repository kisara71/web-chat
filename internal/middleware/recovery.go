package middleware

import (
	"web-chat/pkg/logger"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.RecoveryWithWriter(logger.Writer())
}
