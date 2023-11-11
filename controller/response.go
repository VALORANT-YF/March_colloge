package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
	{
		"code" : 1001 //程序中的错误码
		"msg" : xx 提示的信息
		"data" : {} //返回的数据
	}
*/

//封装响应数据

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

func ResponseError(context *gin.Context, code ResCode) {
	rd := &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	}
	context.JSON(http.StatusOK, rd)
}

func ResponseErrorWithMsg(context *gin.Context, code ResCode, msg interface{}) {
	rd := &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	context.JSON(http.StatusOK, rd)
}

func ResponseSuccessWithData(context *gin.Context, data interface{}) {
	rd := &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	}
	context.JSON(http.StatusOK, rd)
}

func ResponseSuccess(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"code": CodeSuccess,
		"msg":  CodeSuccess.Msg(),
	})
}
