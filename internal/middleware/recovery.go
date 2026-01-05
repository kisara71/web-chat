package middleware

import (
	"web-chat/pkg/logger"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	lgr := logger.L()
	return gin.RecoveryWithWriter(lgr.Writer())
}
