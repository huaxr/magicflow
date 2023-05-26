package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseApiController struct {
}

// Response the unified json structure
type response struct {
	Code    int         `json:"code"`
	Stat    int         `json:"stat"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

type Msg struct {
	Code int
	Msg  string
}

// Respond encapsulates ctx.JSON
func Respond(ctx *gin.Context, status int, msg Msg, data interface{}) {
	respStat := 0
	if msg.Code == 0 {
		respStat = 1
	}
	if data == nil {
		data = gin.H{}
	}
	resp := response{
		Stat:    respStat,
		Code:    msg.Code,
		Message: msg.Msg,
		Data:    data,
	}
	ctx.JSON(status, resp)
}

const (
	SuccessCode      = 0
	DefaultErrorCode = 10000
)

var (
	Success      = Msg{Code: SuccessCode, Msg: "success"}
	DefaultError = Msg{Code: DefaultErrorCode, Msg: "fail"}
)

func (ctl *BaseApiController) Success(ctx *gin.Context, data interface{}) {
	Respond(ctx, http.StatusOK, Success, data)
}

func (ctl *BaseApiController) Error(ctx *gin.Context, code int, message string, data interface{}) {
	msg := Msg{Code: code, Msg: message}
	Respond(ctx, http.StatusOK, msg, data)
}
