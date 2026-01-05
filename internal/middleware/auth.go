package middleware

import (
	"net/http"
	"strings"
	"web-chat/internal/svc"
	"web-chat/pkg/auth"
	errcode "web-chat/pkg/err"
	"web-chat/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const blacklistPrefix = "jwt:blacklist:"

func Auth(svcCtx *svc.Context) gin.HandlerFunc {
	lgr := logger.L()
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			abort(c, http.StatusUnauthorized, errcode.CodeUnauthorized, "missing authorization")
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			abort(c, http.StatusUnauthorized, errcode.CodeUnauthorized, "invalid authorization")
			return
		}
		token := strings.TrimSpace(parts[1])
		if token == "" {
			abort(c, http.StatusUnauthorized, errcode.CodeUnauthorized, "invalid authorization")
			return
		}
		claim := &auth.UserClaim{}
		_, err := svcCtx.Auth.TrackAuthToken(token, claim)
		if err != nil {
			lgr.Printf("auth token error: %v", err)
			abort(c, http.StatusUnauthorized, errcode.CodeUnauthorized, "invalid token")
			return
		}
		ctx := c.Request.Context()
		key := blacklistPrefix + token
		val, err := svcCtx.Infra.Redis.Get(ctx, key).Result()
		if err == nil && val != "" {
			abort(c, http.StatusUnauthorized, errcode.CodeUnauthorized, "token revoked")
			return
		}
		if err != nil && err != redis.Nil {
			lgr.Printf("redis error: %v", err)
			abortInternal(c)
			return
		}
		c.Set("user_id", claim.UserID)
		c.Next()
	}
}
