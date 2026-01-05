package middleware

import (
	"net/http"
	"strconv"
	errcode "web-chat/pkg/err"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var apiErrCodeCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "api_errcode_total",
		Help: "Total number of API responses by errcode.",
	},
	[]string{"method", "path", "errcode"},
)

func init() {
	prometheus.MustRegister(apiErrCodeCounter)
}

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		code := resolveErrCode(c)
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		apiErrCodeCounter.WithLabelValues(
			c.Request.Method,
			path,
			strconv.Itoa(int(code)),
		).Inc()
	}
}

func resolveErrCode(c *gin.Context) errcode.Code {
	if val, ok := c.Get("err_code"); ok {
		if code, ok := val.(int); ok {
			return errcode.Code(code)
		}
	}
	switch c.Writer.Status() {
	case http.StatusUnauthorized:
		return errcode.CodeUnauthorized
	case http.StatusForbidden:
		return errcode.CodeForbidden
	case http.StatusNotFound:
		return errcode.CodeNotFound
	case http.StatusTooManyRequests:
		return errcode.CodeRateLimited
	case http.StatusInternalServerError:
		return errcode.CodeInternal
	default:
		if c.Writer.Status() >= 500 {
			return errcode.CodeInternal
		}
	}
	return errcode.CodeOK
}
