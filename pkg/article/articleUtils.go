package article

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

var start int64 = 1676822400
var end int64 = 1677427200

//爬取简书博客文章相关的第三方包的方法

// ParseBooksHomeHtml 解析简书主页爬取到的Html元素
func ParseBooksHomeHtml(doc *goquery.Document) []string {
	var articleUrls []string //文章链接
	homeInformationHTML, err := doc.Html()
	if err != nil {
		zap.L().Error("doc.Html() is error", zap.Error(err))
		return articleUrls
	}
	homeInformation, err := goquery.NewDocumentFromReader(strings.NewReader(homeInformationHTML))
	if err != nil {
		zap.L().Error("goquery.NewDocumentFromReader(strings.NewReader(homeInformationHTML)) is error", zap.Error(err))
		return articleUrls
	}
	// 获取当前时间
	currentTime := time.Now()
	//获取文章发布时间
	//使用Each方法迭代所有匹配的span标签
	homeInformation.Find("span.time").Each(func(i int, selection *goquery.Selection) {
		// 获取data-shared-at属性的值并添加到切片中
		timeAttr, exists := selection.Attr("data-shared-at")
		if exists {
			//解析时间字符串为时间对象
			sharedTime, err := time.Parse("2006-01-02T15:04:05-07:00", timeAttr)
			if err == nil {
				//计算时间差
				duration := currentTime.Sub(sharedTime)
				//如果时间差在一周内,提取相应的文章链接
				if duration.Hours() < 7*24 {
					//找到相应<a>元素的href属性
					url := selection.Parent().Find("a[target=_blank]").AttrOr("href", "")
					url = fmt.Sprintf("https://www.jianshu.com%s", url)
					articleUrls = append(articleUrls, url)
				}
			}
		}
	})
	//返回简书文章链接切片
	return articleUrls
}

// ParseBookArticleHtml 解析爬取到的简书文章
func ParseBookArticleHtml(doc *goquery.Document) (string, string) {
	var article string
	var articleTitle string
	articleInformationHtml, err := doc.Html()
	if err != nil {
		zap.L().Error("articleInformation , err := doc.Html() is failed", zap.Error(err))
		return article, articleTitle
	}
	articleInformation, err := goquery.NewDocumentFromReader(strings.NewReader(articleInformationHtml))
	if err != nil {
		zap.L().Error("goquery.NewDocumentFromReader(strings.NewReader(homeInformationHTML)) is error", zap.Error(err))
		return article, articleTitle
	}
	// articleInformation.Find("h1._1RuRku").Text()获取文章标题
	//fmt.Println(articleInformation.Find("h1._1RuRku").Text())
	articleTitle = articleInformation.Find("h1._1RuRku").Text()
	// articleInformation.Find("article._2rhmJa").Text() 获取文章内容
	//fmt.Println(articleInformation.Find("article._2rhmJa").Text())
	article = articleInformation.Find("article._2rhmJa").Text()
	//获取文章内容
	return articleTitle, article
}

// ParseBlogHomeHtml 解析爬取到的博客主页信息
func ParseBlogHomeHtml(doc *goquery.Document) []string {
	var blogUrls []string //博客链接
	homeInformationHTML, err := doc.Html()
	if err != nil {
		zap.L().Error("doc.Html() is error", zap.Error(err))
		return blogUrls
	}
	//使用goquery提取文章发布时间
	homeBlogInformation, err := goquery.NewDocumentFromReader(strings.NewReader(homeInformationHTML))
	if err != nil {
		zap.L().Error("goquery.NewDocumentFromReader(strings.NewReader(homeInformationHTML)) is error", zap.Error(err))
		return nil
	}
	//blogTimeText := homeBlogInformation.Find("div.view-time-box").Text() //Text()获得指定元素下所有的文本内容,Html()返回第一个匹配元素的html内容
	// 定义匹配日期格式的正则表达式
	currentTime := time.Now() // 获取当前时间
	datePattern := `\d{4}.\d{2}.\d{2}|\d+\s*小时前|.天\d*`
	//使用Each获得所有的html内容 即 所有博客的发布时间
	homeBlogInformation.Find("div.view-time-box").Each(func(i int, selection *goquery.Selection) {
		// 提取HTML内容
		htmlBlogHomeContent, _ := selection.Html()
		// 使用正则表达式提取日期部分
		re := regexp.MustCompile(datePattern)
		matcheBlog := re.FindString(htmlBlogHomeContent)
		if matcheBlog != "" {
			//昨天,前天,包括几个小时之前,均是最新一周的文章
			if strings.Contains(matcheBlog, "天") {
				blogUrls = append(blogUrls, selection.Parent().Parent().Parent().Parent().Parent().Find("a[target=_blank]").AttrOr("href", ""))
			} else if strings.Contains(matcheBlog, "小时前") {
				// 处理相对时间，例如 "14 小时前"
				hoursAgo, err := time.ParseDuration(strings.ReplaceAll(matcheBlog, " 小时前", "h"))
				if err == nil {
					// 计算实际发布日期
					publishDate := currentTime.Add(-hoursAgo)
					// 判断是否在最新一周内发布
					oneWeekAgo := currentTime.AddDate(0, 0, -7)
					//publishDate 文章发布的时间
					if publishDate.After(oneWeekAgo) {
						//最新一周的博客文章的链接
						blogUrls = append(blogUrls, selection.Parent().Parent().Parent().Parent().Parent().Find("a[target=_blank]").AttrOr("href", ""))
						//fmt.Println("文章在最新一周内发布:", publishDate)
					}
				}
			} else {
				// 处理日期格式 "YYYY.MM.DD"
				publishDate, err := time.Parse("2006.01.02", matcheBlog)
				if err == nil {
					// 判断是否在最新一周内发布
					oneWeekAgo := currentTime.AddDate(0, 0, -7)
					//publishDate 文章发布的时间
					if publishDate.After(oneWeekAgo) {
						//最新一周的博客文章的链接
						blogUrls = append(blogUrls, selection.Parent().Parent().Parent().Parent().Parent().Find("a[target=_blank]").AttrOr("href", ""))
						//fmt.Println("文章在最新一周内发布:", publishDate)
					}
				}
			}
		}
	})
	//blogUrl := fmt.Sprintf(`<a href="(%s/article/details/[^"]+)`, homeUrl) //得到文章标签
	////使用正则表达式提取链接
	//reUrl := regexp.MustCompile(blogUrl)
	//matchesUrl := reUrl.FindAllStringSubmatch(homeInformationHTML, -1)
	//for _, match := range matchesUrl {
	//	blogUrls = append(blogUrls, match[1])
	//}
	//fmt.Println(blogUrls)
	return blogUrls
}

// ParseBlogInformation 爬取博客文章信息
func ParseBlogInformation(blogUrl string) (string, string) {
	var blogTitle string //博客标题
	var blog string      //博客内容
	//发送请求
	response, err := http.Get(blogUrl)
	if err != nil {
		zap.L().Error("ParseBlogInformation http.Get(url) is failed", zap.Error(err))
		return blogTitle, blog
	}
	//使用NewDocumentFromReader 创建文档
	blogContent, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		zap.L().Error("docArticle , err := goquery.NewDocumentFromReader(response.Body) is failed", zap.Error(err))
		return blogTitle, blog
	}
	blogTitle = blogContent.Find("h1.title-article").Text() //获得博客标题
	blog, _ = blogContent.Find("div#content_views").Html()  //获得博客内容(前端页面)
	//fmt.Println(blogTitle)
	//fmt.Println(blog)
	return blogTitle, blog
}
