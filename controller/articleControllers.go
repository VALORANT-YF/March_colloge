package controller

import (
	"college/logic"
	"college/models/bookBlogArticle"
	"college/pkg/article"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ArticleController struct{}

// UserArticle 根据用户的简书主页链接,爬取一周之内的简书文章
func (a ArticleController) UserArticle(context *gin.Context) {
	//首先拿到用户信息和简书主页链接
	err, userArticleInformation := logic.GetArticleUserInformation()
	if err != nil {
		zap.L().Error("logic.GetArticleUserInformation() is failed", zap.Error(err))
		ResponseError(context, CodeServerBusy)
		return
	}
	//所有人的简书信息
	var allPersonArticle []bookBlogArticle.TbBookArticle
	//fmt.Println(userArticleInformation)
	//循环遍历 爬取每一个简书链接不为空的人的简书文章
	for i := 0; i < len(userArticleInformation); i++ {
		if len(userArticleInformation[i].BookAddress) != 0 {
			// 创建 HTTP 请求
			response, err := http.Get(userArticleInformation[i].BookAddress)
			if err != nil {
				zap.L().Error("http.Get(\"https://www.v2ex.com/\") is failed", zap.Error(err))
				continue
			}
			defer response.Body.Close()
			// 使用 NewDocumentFromReader 创建文档
			docHome, err := goquery.NewDocumentFromReader(response.Body)
			if err != nil {
				zap.L().Error("goquery.NewDocumentFromReader(response.Body) is failed", zap.Error(err))
				return
			}
			// 这里可以对主页文档进行解析和提取所需信息
			//最新一周简书文章的链接
			articleUrl := article.ParseBooksHomeHtml(docHome)
			//根据拿到的最新一周的简书文章的链接,爬取简书文章内容
			for _, url := range articleUrl {
				//发送请求
				response, err = http.Get(url)
				if err != nil {
					zap.L().Error("http.Get(url) is failed", zap.Error(err))
					return
				}
				//使用NewDocumentFromReader 创建文档
				docArticle, err := goquery.NewDocumentFromReader(response.Body)
				if err != nil {
					zap.L().Error("docArticle , err := goquery.NewDocumentFromReader(response.Body) is failed", zap.Error(err))
					return
				}
				//提取文章信息
				articleTitle, articleContent := article.ParseBookArticleHtml(docArticle) //简书文章标题和文章内容
				// 创建一个正则表达式，匹配任何空白字符
				re := regexp.MustCompile(`\s+`)
				tempArticleContent := re.ReplaceAllString(articleContent, "")
				deptNameArr := strings.Split(userArticleInformation[i].DeptName, " ")
				for _, deptName := range deptNameArr {
					//判断这个部门是否需要写简书
					if logic.IsWriteBookService(deptName) == 1 {
						//如果不是空白文章 , 将文章和对应信息插入数据库中
						var bookInformations bookBlogArticle.TbBookArticle
						if len(tempArticleContent) > 0 {
							//fmt.Print(articleTitle)
							bookInformations.BookTitle = articleTitle                  //简书标题
							bookInformations.BookArticle = tempArticleContent          //简书文章
							bookInformations.Mobile = userArticleInformation[i].Mobile //简书主人的电话
							bookInformations.Name = userArticleInformation[i].Name     //简书主人的姓名
							bookInformations.DeptName = deptName                       //简书主人所属的部门名称
							bookInformations.ArticleUrl = url                          //简书原文链接
						}
						allPersonArticle = append(allPersonArticle, bookInformations)
					}
				}
			}
		}
	}
	//调用Service层插入数据
	err = logic.InsertAllArticleService(allPersonArticle)
	if err != nil {
		zap.L().Error("logic.InsertAllArticleService(allPersonArticle) is failed", zap.Error(err))
		return
	}
	zap.L().Info("简书插入成功")
}

// UserBlog 定时任务 爬取每周最新的一篇博客
func (a ArticleController) UserBlog(context *gin.Context) {
	//首先拿到用户信息和博客主页链接
	err, userBlogInformation := logic.GetArticleUserInformation()
	if err != nil {
		zap.L().Error("logic.GetArticleUserInformation() is failed", zap.Error(err))
		ResponseError(context, CodeServerBusy)
		return
	}
	//所有人的博客信息
	var allPersonBlog []bookBlogArticle.TbBlog
	//循环遍历 , 爬取所有人的 博客文章
	for i := 0; i < len(userBlogInformation); i++ {
		if len(userBlogInformation[i].BlogAddress) != 0 {

			//###########
			/*if userBlogInformation[i].Name == "张云飞" {
				for j := 0; j <= 10000; j++ {
					这段代码有神奇的力量
				}
			}*/
			//#########

			//创建http请求
			response, err := http.Get(userBlogInformation[i].BlogAddress)
			if err != nil {
				zap.L().Error("http.Get(bolg_address) is failed", zap.Error(err))
				continue
			}
			defer response.Body.Close()
			// 使用 NewDocumentFromReader 创建文档
			docHome, err := goquery.NewDocumentFromReader(response.Body)
			if err != nil {
				zap.L().Error("goquery.NewDocumentFromReader(response.Body) is failed", zap.Error(err))
				continue
			}
			blogUrls := article.ParseBlogHomeHtml(docHome) //拿到最新一周博客内容的链接
			var personBlog bookBlogArticle.TbBlog          // 一篇博客的信息
			deptNameArr := strings.Split(userBlogInformation[i].DeptName, " ")
			for _, value := range deptNameArr {
				//判断这个部门是否需要写博客
				if logic.IsWriteBlogService(value) == 1 {
					//根据拿到的博客链接,循环遍历链接拿到博客内容
					for _, blogUrl := range blogUrls {
						blogTitle, blog := article.ParseBlogInformation(blogUrl) //拿到博客标题和博客文章
						personBlog.Name = userBlogInformation[i].Name            //博客主人姓名
						personBlog.Mobile = userBlogInformation[i].Mobile        //博客主人电话
						personBlog.DeptName = value                              //博客主人所在的部门
						personBlog.BookTitle = blogTitle                         //博客标题
						personBlog.BookArticle = blog                            //博客内容(包含前端页面)
						personBlog.BlogUrl = blogUrl                             //博客原文链接
						allPersonBlog = append(allPersonBlog, personBlog)
					}
				}
			}

		}
	}
	//所有本周博客的信息
	//fmt.Println(allPersonBlog)
	//调用Service层插入数据
	err = logic.InsertAllPersonBlogService(allPersonBlog)
	if err != nil {
		zap.L().Error("logic.InsertAllPersonBlogService(allPersonBlog) is failed", zap.Error(err))
		return
	}
	zap.L().Info("博客插入成功")
}

// GetAllBook 查询所有人的简书
func (a ArticleController) GetAllBook(context *gin.Context) {
	err, allPersonBook := logic.SelectAllPersonBookService()

	if err != nil {
		zap.L().Error("logic.SelectAllPersonBookService() is failed", zap.Error(err))
		return
	}
	ResponseSuccessWithData(context, allPersonBook)
}

// GetAllBlog 查询所有人的博客
func (a ArticleController) GetAllBlog(context *gin.Context) {
	err, allPersonBlog := logic.SelectAllPersonBlogService()
	if err != nil {
		zap.L().Error("logic.SelectAllPersonBlogService() is failed", zap.Error(err))
		return
	}
	ResponseSuccessWithData(context, allPersonBlog)
}
