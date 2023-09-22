package router

import (
	"bluebell/logger"
	"bluebell/settings"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	// 设置 gin 框架日志输出模式
	gin.SetMode(settings.Conf.GinConfig.Mode)
	// 创建一个路由引擎
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "122")
	})

	return r
}
