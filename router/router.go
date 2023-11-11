package router

import (
	"college/logger"
	middles "college/middlewares"
	"college/router/routersGroup/adminGroups"
	"college/router/routersGroup/articleGroups"
	"college/router/routersGroup/deptGroups"
	"college/router/routersGroup/dingOfficialGroups"
	"college/router/routersGroup/userGroups"
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	// 设置 gin 框架日志输出模式
	//gin.SetMode(settings.Conf.GinConfig.Mode)
	// 创建一个路由引擎
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	//使用全局中间件 配置跨域
	r.Use(middles.CORSMiddleware())
	zap.L().Debug("跨域配置成功")
	//抽取路由组件
	adminGroups.AdminRouters(r)               //管理员相关的路由
	dingOfficialGroups.DingOfficialRouters(r) //钉钉官方接口的路由组件
	userGroups.UserRouterGroups(r)            //用户相关的路由组件
	articleGroups.ArticleGroupRouters(r)      //简书博客文章相关的路由组件
	deptGroups.DeptRouterGroups(r)            //部门相关的路由组件

	r.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "122")
	})

	return r
}
