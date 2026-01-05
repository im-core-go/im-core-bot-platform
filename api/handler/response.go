package handler

import (
	"net/http"
	errcode "web-chat/pkg/err"

	"github.com/gin-gonic/gin"
)

type response struct {
	Code    errcode.Code `json:"code"`
	Message string       `json:"message"`
	Data    any          `json:"data,omitempty"`
}

func Respond(c *gin.Context, httpStatus int, code errcode.Code, message string, data any) {
	c.Set("err_code", int(code))
	c.JSON(httpStatus, response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func OK(c *gin.Context, data any) {
	Respond(c, http.StatusOK, errcode.CodeOK, "OK", data)
}
