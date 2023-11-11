package articleGroups

import (
	"college/controller"

	"github.com/gin-gonic/gin"
)

var articleController controller.ArticleController

func ArticleGroupRouters(r *gin.Engine) {
	article := r.Group("/article")
	{
		article.POST("/getArticle", articleController.UserArticle) //根据每一个用户简书的主页链接去拿到最新一周的简书文章
		article.POST("/getBlog", articleController.UserBlog)       //根据每一个用户的博客主页链接去拿到最新一周的博客
		article.GET("/selectBook", articleController.GetAllBook)   //查询本周的简书博客
	}
}
