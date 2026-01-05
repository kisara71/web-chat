package middleware

import (
	"net/http"
	errcode "web-chat/pkg/err"

	"github.com/gin-gonic/gin"
)

type response struct {
	Code    errcode.Code `json:"code"`
	Message string       `json:"message"`
}

func abort(c *gin.Context, httpStatus int, code errcode.Code, message string) {
	c.Set("err_code", int(code))
	c.AbortWithStatusJSON(httpStatus, response{
		Code:    code,
		Message: message,
	})
}

func abortInternal(c *gin.Context) {
	abort(c, http.StatusInternalServerError, errcode.CodeInternal, "internal error")
}
