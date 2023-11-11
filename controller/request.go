package controller

import (
	"college/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const tokenFailed = "no"

// GetUnionidToken 验证用户登录状态
func GetUnionidToken(context *gin.Context) string {
	//获得用户的唯一id unionid
	aToken, ok := context.Get(CtxtUserIDKey)
	if !ok {
		return tokenFailed
	}
	unionid, ok := aToken.(string)
	if !ok {
		zap.L().Error("UserInputBookBlogAddress aToken.(string) is error")
		return tokenFailed
	}
	return unionid
}

// GetIsBoss 用户鉴权
func GetIsBoss(unionid string) bool {
	return logic.QueryIsBossService(unionid)
}
