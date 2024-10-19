package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func response(ctx *gin.Context, httpCode int, response *Response) {
	ctx.JSON(httpCode, response)
}

// FailWithMsg returns a failed response with a message
func ErrorWithMsg(ctx *gin.Context, msg string) {
	response(ctx, http.StatusOK, &Response{
		Code: CodeServerInternalError,
		Msg:  msg,
		Data: nil,
	})
}

// FailWithCode returns a failed response with a code
func ErrorWithCode(ctx *gin.Context, code Code) {
	response(ctx, http.StatusOK, &Response{
		Code: int(code),
		Msg:  code.ToMsg(),
		Data: nil,
	})
}

func ErrorWithCodeAndMsg(ctx *gin.Context, code Code, msg string) {
	response(ctx, http.StatusOK, &Response{
		Code: int(code),
		Msg:  msg,
		Data: nil,
	})
}

func ErrorWithData(ctx *gin.Context, code Code, data interface{}) {
	response(ctx, http.StatusOK, &Response{
		Code: int(code),
		Msg:  code.ToMsg(),
		Data: data,
	})
}

// Ok returns a successful response
func OkWithMsg(ctx *gin.Context, msg string) {
	response(ctx, http.StatusOK, &Response{
		Code: CodeSuccess,
		Msg:  msg,
		Data: nil,
	})
}

// OkWithData returns a successful response with data
func OkWithData(ctx *gin.Context, data interface{}) {
	response(ctx, http.StatusOK, &Response{
		Code: CodeSuccess,
		Msg:  Code(CodeSuccess).ToMsg(),
		Data: data,
	})
}
