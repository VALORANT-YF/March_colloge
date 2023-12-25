package userGroups

import (
	"college/controller"
	middles "college/middlewares"

	"github.com/gin-gonic/gin"
)

var userControllers controller.UserController

// UserRouterGroups 用户操作相关的路由
func UserRouterGroups(r *gin.Engine) {
	users := r.Group("/user")
	{
		users.POST("/login", userControllers.UserLogin) //普通用户登录
	}

	//用户登录之后的路由需要使用路由中间件jwtToken,来验证用户的身份
	userAfterLogin := r.Group("/userLogin", middles.JWTAuthMiddleware())
	{
		//查询简书和博客的主页地址来判断用户是否是第一次登录
		userAfterLogin.POST("/bookAddress", userControllers.UserInitialBookBlog)
		//新用户录入自己的简书和博客主页链接
		userAfterLogin.POST("/inputAddress", userControllers.UserInputBookBlogAddress)
		//用户查看其他人的简书或者博客
		userAfterLogin.GET("/lookArticle", userControllers.GetOtherBookBlog)
		//查看登录用户的信息
		userAfterLogin.GET("/selfInformation", userControllers.GetSelfInformation)
		//用户修改自己的简书博客地址,或者密码
		userAfterLogin.POST("/updateSelfInformation", userControllers.UpdateSelfInformation)
	}
}
