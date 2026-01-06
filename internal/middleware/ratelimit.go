package middleware

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"web-chat/internal/svc"
	errcode "web-chat/pkg/err"
	"web-chat/pkg/logger"

	"github.com/gin-gonic/gin"
)

const rateLimitPrefix = "ratelimit:"

const slidingWindowScript = `
local key = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local member = ARGV[4]
redis.call('ZREMRANGEBYSCORE', key, 0, now - window)
redis.call('ZADD', key, now, member)
local count = redis.call('ZCARD', key)
redis.call('PEXPIRE', key, window)
if count > limit then
	return 0
end
return 1
`

func RateLimit(svcCtx *svc.Context, limit int, window time.Duration) gin.HandlerFunc {
	lgr := logger.L()
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("%s%s", rateLimitPrefix, ip)
		ctx := c.Request.Context()
		now := time.Now().UnixMilli()
		member := fmt.Sprintf("%d-%d", now, rand.Int63())
		allowed, err := svcCtx.Infra.Redis.Eval(ctx, slidingWindowScript, []string{key}, now, window.Milliseconds(), limit, member).Int()
		if err != nil {
			lgr.Errorf("rate limit redis error: %v", err)
			abortInternal(c)
			return
		}
		if allowed == 0 {
			abort(c, http.StatusTooManyRequests, errcode.CodeRateLimited, "rate limit exceeded")
			return
		}
		c.Next()
	}
}
